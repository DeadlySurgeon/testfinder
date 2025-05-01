package examples

import (
	"reflect"
	"testing"

	"github.com/deadlysurgeon/testfinder/extra"
)

func TestRngSplitSliceTest(t *testing.T) {
	for _, test := range []struct {
		name  string
		input string
		sep   string
		want  []string
	}{
		{
			name:  "simple",
			input: "a/b/c",
			sep:   "/",
			want:  []string{"a", "b", "c"},
		},
		{
			name:  "wrong sep",
			input: "a/b/c",
			sep:   ",",
			want:  []string{"a/b/c"},
		},
		{
			name:  "no sep",
			input: "abc",
			sep:   "/",
			want:  []string{"abc"},
		},
		{
			name:  "trailing sep",
			input: "a/b/c/",
			sep:   "/",
			want:  []string{"a", "b", "c", ""},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			got := Split(extra.ID(test.input), test.sep)
			if !reflect.DeepEqual(test.want, got) {
				t.Fatalf("expected: %v, got: %v", test.want, got)
			}
		})
	}
}
