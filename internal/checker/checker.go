package checker

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Yamashou/gqlgenc/client"
	"github.com/Yamashou/gqlgenc/introspection"
	"github.com/mattn/go-zglob"
	gqlparser "github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/parser"
	"github.com/vektah/gqlparser/v2/validator"
	"go.uber.org/multierr"

	"github.com/gqlgo/gqlanalysis"
)

type Checker struct {
	Schema              string
	Query               string
	Stderr              io.Writer
	Fsys                fs.FS
	HTTPClient          *http.Client
	IntrospectionHeader http.Header
}

func (c *Checker) stderr() io.Writer {
	if c.Stderr != nil {
		return c.Stderr
	}
	return os.Stderr
}

func (c *Checker) readFile(path string) ([]byte, error) {
	if c.Fsys != nil {
		return fs.ReadFile(c.Fsys, path)
	}
	return os.ReadFile(filepath.FromSlash(path))
}

func (c *Checker) walk(pattern string, fn filepath.WalkFunc) error {
	if c.Fsys != nil {
		return fs.WalkDir(c.Fsys, ".", func(path string, d fs.DirEntry, _err error) error {
			if _err != nil {
				return _err
			}
			info, err := d.Info()
			if err != nil {
				return err
			}

			matches, err := zglob.Match(pattern, path)
			if err != nil {
				return fmt.Errorf("matching to glob: %w", err)
			}

			if matches {
				return fn(path, info, _err)
			}

			return nil
		})
	}

	files, err := zglob.Glob(pattern)
	if err != nil {
		return fmt.Errorf("parse glob: %w", err)
	}

	for _, file := range files {
		info, err := os.Stat(file)
		if err := fn(file, info, err); err != nil {
			return err
		}
	}

	return nil
}

func (c *Checker) httpClient() *http.Client {
	if c.HTTPClient != nil {
		return c.HTTPClient
	}
	return http.DefaultClient
}

func (c *Checker) Run(analyzers ...*gqlanalysis.Analyzer) error {

	var comments []*gqlanalysis.Comment

	schema, scomments, err := c.parseSchema()
	if err != nil {
		return err
	}
	comments = append(comments, scomments...)

	queries, qcomments, err := c.parseQuery(schema)
	if err != nil {
		return err
	}
	comments = append(comments, qcomments...)

	var errs error
	acts := c.analyze(schema, queries, comments, analyzers)
	for _, act := range acts {
		errs = multierr.Append(errs, act.err)
	}
	if errs != nil {
		return errs
	}

	for _, act := range acts {
		if len(act.diagnostics) == 0 {
			continue
		}
		fmt.Fprintf(c.stderr(), "# results of analyzer %s", act.a.Name)
		for _, d := range act.diagnostics {
			fmt.Fprintf(c.stderr(), "%s:%d %s\n", d.Pos.Src.Name, d.Pos.Line, d.Message)
		}
	}

	return nil
}

func (c *Checker) parseSchema() (*ast.Schema, []*gqlanalysis.Comment, error) {
	switch {
	case strings.HasPrefix(c.Schema, "https://"), strings.HasPrefix(c.Schema, "http://"):
		return c.parseSchemaFromIntrospection()
	default:
		return c.parseSchemaFromFiles()
	}
}

func (c *Checker) parseSchemaFromFiles() (*ast.Schema, []*gqlanalysis.Comment, error) {
	var srcs []*ast.Source
	err := c.walk(c.Schema, func(_path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("parse schema in %s: %w", _path, err)
		}

		slashPath := filepath.ToSlash(_path)

		if info.IsDir() {
			return nil
		}

		content, err := c.readFile(slashPath)
		if err != nil {
			return fmt.Errorf("read file %s: %w", _path, err)
		}

		srcs = append(srcs, &ast.Source{
			Name:  _path,
			Input: string(content),
		})

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	// LoadSchemas returns *gqlerror.Error but not error
	schema, loadErr := gqlparser.LoadSchema(srcs...)
	if loadErr != nil {
		return nil, nil, loadErr
	}

	var comments []*gqlanalysis.Comment
	for _, src := range srcs {
		comment, err := gqlanalysis.ReadComments(src)
		if err != nil {
			return nil, nil, err
		}
		comments = append(comments, comment...)
	}

	return schema, comments, nil
}

func (c *Checker) parseSchemaFromIntrospection() (*ast.Schema, []*gqlanalysis.Comment, error) {
	gqlclient := client.NewClient(c.httpClient(), c.Schema, func(req *http.Request) {
		for key := range c.IntrospectionHeader {
			for _, value := range c.IntrospectionHeader[key] {
				req.Header.Add(key, value)
			}
		}
	})
	var res introspection.Query
	if err := gqlclient.Post(context.Background(), "Query", introspection.Introspection, &res, nil); err != nil {
		return nil, nil, fmt.Errorf("introspection query failed: %w", err)
	}
	schema, err := validator.ValidateSchemaDocument(introspection.ParseIntrospectionQuery(c.Schema, res))
	if err != nil {
		return nil, nil, fmt.Errorf("validation error: %w", err)
	}
	return schema, nil, nil
}

func (c *Checker) parseQuery(schema *ast.Schema) ([]*ast.QueryDocument, []*gqlanalysis.Comment, error) {
	var (
		queries  []*ast.QueryDocument
		comments []*gqlanalysis.Comment
	)

	err := c.walk(c.Query, func(_path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("parse query in %s: %w", _path, err)
		}

		slashPath := filepath.ToSlash(_path)

		if info.IsDir() {
			return nil
		}

		content, err := c.readFile(slashPath)
		if err != nil {
			return fmt.Errorf("read file %s: %w", _path, err)
		}

		src := &ast.Source{
			Name:  _path,
			Input: string(content),
		}

		// ParseQuery returns gqlerror.List but not error
		q, parseErr := parser.ParseQuery(src)
		if parseErr != nil {
			return parseErr
		}
		queries = append(queries, q)

		cmts, err := gqlanalysis.ReadComments(src)
		if err != nil {
			return err
		}
		comments = append(comments, cmts...)

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	merge(schema, queries)

	return queries, comments, nil
}

func merge(schema *ast.Schema, queries []*ast.QueryDocument) {
	// Because validator.Walk accepts only single QueryDocument,
	// queries must be merged into single QueryDocument.
	var query ast.QueryDocument
	for _, q := range queries {
		query.Operations = append(query.Operations, q.Operations...)
		query.Fragments = append(query.Fragments, q.Fragments...)
	}
	validator.Walk(schema, &query, new(validator.Events))
}

func (c *Checker) analyze(schema *ast.Schema, queries []*ast.QueryDocument, comments []*gqlanalysis.Comment, analyzers []*gqlanalysis.Analyzer) []*action {
	actions := make(map[*gqlanalysis.Analyzer]*action)
	var mkAction func(a *gqlanalysis.Analyzer) *action
	mkAction = func(a *gqlanalysis.Analyzer) *action {
		act, ok := actions[a]
		if !ok {
			act = &action{
				a:        a,
				schema:   schema,
				queries:  queries,
				comments: comments,
			}
			for _, req := range a.Requires {
				act.deps = append(act.deps, mkAction(req))
			}
			actions[a] = act
		}
		return act
	}

	var roots []*action
	for _, a := range analyzers {
		root := mkAction(a)
		root.isroot = true
		roots = append(roots, root)
	}

	execAll(roots)

	return roots
}

type AnalyzerResult struct {
	Pass        *gqlanalysis.Pass
	Diagnostics []*gqlanalysis.Diagnostic
	Result      interface{}
	Err         error
}

func (c *Checker) RunSingle(a *gqlanalysis.Analyzer) ([]*AnalyzerResult, error) {

	var comments []*gqlanalysis.Comment

	schema, scomments, err := c.parseSchema()
	if err != nil {
		return nil, err
	}
	comments = append(comments, scomments...)

	queries, qcomments, err := c.parseQuery(schema)
	if err != nil {
		return nil, err
	}
	comments = append(comments, qcomments...)

	acts := c.analyze(schema, queries, comments, []*gqlanalysis.Analyzer{a})
	if acts == nil {
		return nil, nil
	}

	results := make([]*AnalyzerResult, len(acts))
	for i, act := range acts {
		results[i] = &AnalyzerResult{
			Pass:        act.pass,
			Diagnostics: act.diagnostics,
			Result:      act.result,
			Err:         act.err,
		}
	}

	return results, nil
}
