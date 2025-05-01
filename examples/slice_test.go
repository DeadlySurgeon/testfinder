package examples

import (
	"reflect"
	"testing"

	"github.com/deadlysurgeon/testfinder/extra"
)

func TestSplitSliceTest(t *testing.T) {
	tests := []struct {
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
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := Split(extra.ID(test.input), test.sep)
			if !reflect.DeepEqual(test.want, got) {
				t.Fatalf("expected: %v, got: %v", test.want, got)
			}
		})
	}
}
