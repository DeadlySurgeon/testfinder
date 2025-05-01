package examples

import (
	"reflect"
	"testing"

	"github.com/deadlysurgeon/testfinder/extra"
)

func TestRngSplitMapTest(t *testing.T) {
	for name, test := range map[string]struct {
		input string
		sep   string
		want  []string
	}{
		"simple": {
			input: "a/b/c",
			sep:   "/",
			want:  []string{"a", "b", "c"},
		},
		"wrong sep": {
			input: "a/b/c",
			sep:   ",",
			want:  []string{"a/b/c"},
		},
		"no sep": {
			input: "abc",
			sep:   "/",
			want:  []string{"abc"},
		},
		"trailing sep": {
			input: "a/b/c/",
			sep:   "/",
			want:  []string{"a", "b", "c", ""},
		},
	} {
		t.Run(name, func(t *testing.T) {
			got := Split(extra.ID(test.input), test.sep)
			if !reflect.DeepEqual(test.want, got) {
				t.Fatalf("expected: %v, got: %v", test.want, got)
			}
		})
	}
}
