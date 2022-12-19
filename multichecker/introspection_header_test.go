package multichecker_test

import (
	"testing"

	"github.com/gqlgo/gqlanalysis/multichecker"
)

func TestIntrospectionHeader_Set(t *testing.T) {
	cases := []struct {
		args string
		want map[string]string
	}{
		{"key1:value1", map[string]string{"key1": "value1"}},
		{"key1:value1,key2:value2", map[string]string{"key1": "value1", "key2": "value2"}},
	}

	for _, tt := range cases {
		actual := make(multichecker.ExportedIntrospectionHeader)
		actual.Set(tt.args)

		if len(actual) != len(tt.want) {
			t.Errorf("len(actual) != len(tt.want). actual: %v, want: %v", actual, tt.want)
		}
		for wantKey, wantValue := range tt.want {
			actualValue, ok := actual[wantKey]
			if !ok {
				t.Errorf("Does not contain key: %s", wantKey)
			}
			if wantValue != actualValue[0] {
				t.Errorf("wantValue is not equal actualValue. key: %v, wantKey: %v, actualValue: %v", wantKey, wantValue, actualValue)
			}
		}
	}
}

func TestIntrospectionHeader_Value(t *testing.T) {
	cases := []struct {
		want string
	}{
		{"key1:value1"},
		{"key1:value1,key2:value2"},
	}

	for _, tt := range cases {
		ih := make(multichecker.ExportedIntrospectionHeader)
		ih.Set(tt.want)
		actual := ih.String(
)
		if actual != tt.want {
			t.Errorf("actual != tt.want. actual: %v, want: %v", actual, tt.want)
		}
	}
}
