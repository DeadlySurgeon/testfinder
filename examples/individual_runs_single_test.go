package examples

import (
	"reflect"
	"testing"

	"github.com/deadlysurgeon/testfinder/extra"
)

func TestSplitIndividualSingleTest(t *testing.T) {
	type test struct {
		input string
		sep   string
		want  []string
	}

	perform := func(tc test) func(t *testing.T) {
		return func(t *testing.T) {
			got := Split(extra.ID(tc.input), tc.sep)
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		}
	}

	t.Run("simple", perform(test{
		input: "a/b/c",
		sep:   "/",
		want:  []string{"a", "b", "c"},
	}))

	t.Run("wrong sep", perform(test{
		input: "a/b/c",
		sep:   ",",
		want:  []string{"a/b/c"},
	}))

	t.Run("no sep", perform(test{
		input: "abc",
		sep:   "/",
		want:  []string{"abc"},
	}))

	t.Run("trailing sep", perform(test{
		input: "a/b/c/",
		sep:   "/",
		want:  []string{"a", "b", "c", ""},
	}))
}
