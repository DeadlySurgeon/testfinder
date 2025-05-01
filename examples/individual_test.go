package examples

import (
	"reflect"
	"testing"
)

func TestSplitSimple(t *testing.T) {
	got := Split("a/b/c", "/")
	want := []string{"a", "b", "c"}
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("expected: %v, got: %v", want, got)
	}
}

func TestSplitWrongSep(t *testing.T) {
	got := Split("a/b/c", ",")
	want := []string{"a/b/c"}
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("expected: %v, got: %v", want, got)
	}
}

func TestSplitNoSep(t *testing.T) {
	got := Split("abc", "/")
	want := []string{"abc"}
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("expected: %v, got: %v", want, got)
	}
}

func TestSplitTrailingSep(t *testing.T) {
	got := Split("a/b/c/", "/")
	want := []string{"a", "b", "c", ""}
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("expected: %v, got: %v", want, got)
	}
}
