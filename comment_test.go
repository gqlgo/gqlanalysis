package gqlanalysis_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/vektah/gqlparser/v2/ast"

	"github.com/gqlgo/gqlanalysis"
)

func TestReadComments(t *testing.T) {
	t.Parallel()

	type (
		S    = []string
		I    = []int
		want struct {
			comments S
			lines    I
			cols     I
			err      bool
		}
	)

	cases := map[string]struct {
		content string
		want    want
	}{
		"empty":              {"", want{S{}, I{}, I{}, false}},
		"normal":             {" # test", want{S{"# test"}, I{1}, I{2}, false}},
		"2line":              {"\n# test", want{S{"# test"}, I{2}, I{1}, false}},
		"2comments":          {"# test1\n# test2", want{S{"# test1", "# test2"}, I{1, 2}, I{1, 1}, false}},
		"2comments-sameline": {"# test1# test2", want{S{"# test1# test2"}, I{1}, I{1}, false}},
		"double-sharp":       {"## test1", want{S{"## test1"}, I{1}, I{1}, false}},
	}

	for name, tt := range cases {
		name, tt := name, tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			src := &ast.Source{
				Name:  name,
				Input: tt.content,
			}

			got, err := gqlanalysis.ReadComments(src)

			switch {
			case !tt.want.err && err != nil:
				t.Fatal("unexpected error", err)
			case tt.want.err && err == nil:
				t.Fatal("the expected error does not occur", err)
			}

			comments := make([]string, len(got))
			lines := make([]int, len(got))
			cols := make([]int, len(got))

			for i := range got {
				comments[i] = got[i].Value
				lines[i] = got[i].Pos.Line
				cols[i] = got[i].Pos.Column
			}

			if diff := cmp.Diff(tt.want.comments, comments); diff != "" {
				t.Error("comments", diff)
			}

			if diff := cmp.Diff(tt.want.lines, lines); diff != "" {
				t.Error("lines", diff)
			}

			if diff := cmp.Diff(tt.want.cols, cols); diff != "" {
				t.Error("cols", diff)
			}
		})
	}
}
