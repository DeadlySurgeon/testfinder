package examples

import (
	"reflect"
	"testing"
)

func TestSplitIndivual(t *testing.T) {
	sep := "/"

	t.Run("simple", func(t *testing.T) {
		got := Split("a/b/c", sep)
		want := []string{"a", "b", "c"}
		if !reflect.DeepEqual(want, got) {
			t.Fatalf("expected: %v, got: %v", want, got)
		}
	})

	t.Run("wrong sep", func(t *testing.T) {
		got := Split("a/b/c", ",")
		want := []string{"a/b/c"}
		if !reflect.DeepEqual(want, got) {
			t.Fatalf("expected: %v, got: %v", want, got)
		}
	})

	t.Run("no sep", func(t *testing.T) {
		got := Split("abc", sep)
		want := []string{"abc"}
		if !reflect.DeepEqual(want, got) {
			t.Fatalf("expected: %v, got: %v", want, got)
		}
	})

	t.Run("trailing sep", func(t *testing.T) {
		got := Split("a/b/c/", sep)
		want := []string{"a", "b", "c", ""}
		if !reflect.DeepEqual(want, got) {
			t.Fatalf("expected: %v, got: %v", want, got)
		}
	})
}
