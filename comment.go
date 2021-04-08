package gqlanalysis

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/vektah/gqlparser/v2/ast"
)

// Comment represents a comment line of a graphql file.
type Comment struct {
	Value string
	Pos   ast.Position
}

// String implements fmt.Stringer.
func (c *Comment) String() string {
	return c.Value
}

var _ fmt.Stringer = (*Comment)(nil)

// ReadComments reads comments from src.
// Ofcourse, src is io.Reader better than *ast.Source but it use *ast.Soruce
// for compatibility with gqlparser.
// Unfortunately, gqlparser does not toknize comments.
// See: https://github.com/vektah/gqlparser/issues/145
func ReadComments(src *ast.Source) ([]*Comment, error) {
	r := &readCounter{org: strings.NewReader(src.Input)}
	scanner := bufio.NewScanner(r)
	var (
		comments []*Comment
		start    int
		l        int
	)

	for scanner.Scan() {
		line := scanner.Text()
		l++
		i := strings.Index(line, "#")
		if i == -1 {
			continue
		}

		comments = append(comments, &Comment{
			Value: line[i:],
			Pos: ast.Position{
				Start:  start,
				End:    r.cnt,
				Line:   l,
				Column: i + 1,
				Src:    src,
			},
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("analysis.ReadComments: %w", err)
	}

	return comments, nil
}

type readCounter struct {
	org io.Reader
	cnt int
}

func (r *readCounter) Read(p []byte) (int, error) {
	n, err := r.org.Read(p)
	r.cnt += n
	return n, err
}
