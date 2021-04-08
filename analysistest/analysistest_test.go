package analysistest_test

import (
	"testing"

	"github.com/gqlgo/gqlanalysis"
	"github.com/gqlgo/gqlanalysis/analysistest"
)

type noopTestingT struct {
	*testing.T
	reported bool
}

func (t *noopTestingT) Errorf(_ string, _ ...interface{}) { t.reported = true }
func (t *noopTestingT) Error(_ ...interface{})            { t.reported = true }
func (t *noopTestingT) Fatalf(_ string, _ ...interface{}) { t.reported = true }
func (t *noopTestingT) Fatal(_ ...interface{})            { t.reported = true }

func TestRun(t *testing.T) {
	t.Parallel()
	type D struct { // want diagnositc
		line int
		msg  string
	}
	cases := map[string]struct {
		dir        string
		wantReport bool
		want       *D
	}{
		"nodiagnotics":    {"a", false, nil},
		"diagnotics":      {"b", false, &D{2, "NG"}},
		"lackwant":        {"a", true, &D{2, "NG"}},
		"unnecessarywant": {"b", true, nil},
	}

	for name, tt := range cases {
		name, tt := name, tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			testdata := analysistest.TestData(t)
			a := &gqlanalysis.Analyzer{
				Name: name,
				Doc:  name,
				Run: func(pass *gqlanalysis.Pass) (interface{}, error) {
					if tt.want == nil {
						return nil, nil
					}
					for _, q := range pass.Queries {
						for _, op := range q.Operations {
							for _, sel := range op.SelectionSet {
								if tt.want.line == sel.GetPosition().Line {
									pass.Reportf(sel.GetPosition(), tt.want.msg)
								}
							}
						}
					}
					return nil, nil
				},
			}

			nt := &noopTestingT{t, false}
			result := analysistest.Run(nt, testdata, a, tt.dir)[0]

			if nt.reported != tt.wantReport {
				t.Errorf("reported want %v but got %v", tt.wantReport, nt.reported)
			}

			switch {
			case tt.want != nil && len(result.Diagnostics) == 0:
				t.Fatal("expected diagnositc was not reported")
			case tt.want == nil && len(result.Diagnostics) != 0:
				t.Fatalf("unexpected diagnositc was occured: %v", result.Diagnostics)
			case tt.want == nil && len(result.Diagnostics) == 0:
				return
			}

			got := result.Diagnostics[0]
			if tt.want.msg != got.Message {
				t.Errorf("diagnositc message want %v but got %v", tt.want.msg, got.Message)
			}

			if tt.want.line != got.Pos.Line {
				t.Errorf("diagnositc message want %v but got %v", tt.want.msg, got.Pos.Line)
			}
		})
	}
}
