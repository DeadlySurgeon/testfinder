package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"
)

func TestMain(t *testing.T) {
	type testCase struct {
		location        string
		expectedErr     string
		expectedResults []testInfo
	}

	performTest := func(tc testCase) func(*testing.T) {
		return func(t *testing.T) {
			testOutput := runTest(t, tc.location)

			if tc.expectedErr == "" && testOutput.err != nil {
				t.Fatal(`Expected no error but got "` + testOutput.err.Error() + `"`)
			} else if tc.expectedErr != "" && tc.expectedErr != errStr(testOutput.err) {
				t.Fatal(`Expected error "` + tc.expectedErr + `" but got "` + errStr(testOutput.err) + `"`)
			}

			fmt.Printf("%#v\n", testOutput.results)
			if len(tc.expectedResults) != len(testOutput.results) {
				t.Fatalf("Expected %d results but got %d", len(tc.expectedResults), len(testOutput.results))
			}

			for i, expected := range tc.expectedResults {
				result := testOutput.results[i]
				if result.Line != expected.Line {
					t.Fatalf("Expected result[%d] line to be %d but got %d", i, expected.Line, result.Line)
				}
				if result.SubTest != expected.SubTest {
					t.Fatalf("Expected result[%d] sub test to be %s but got %s", i, expected.SubTest, result.SubTest)
				}
				if result.TestFunc != expected.TestFunc {
					t.Fatalf("Expected result[%d] test func to be %s but got %s", i, expected.TestFunc, result.TestFunc)
				}
			}
		}
	}

	t.Run("direct_map_test", performTest(testCase{
		location: "./examples/direct_map_test.go",
		expectedResults: []testInfo{
			{TestFunc: "TestRngSplitMapTest", SubTest: "", Line: 10},
			{TestFunc: "TestRngSplitMapTest", SubTest: "simple", Line: 16},
			{TestFunc: "TestRngSplitMapTest", SubTest: "wrong sep", Line: 21},
			{TestFunc: "TestRngSplitMapTest", SubTest: "no sep", Line: 26},
			{TestFunc: "TestRngSplitMapTest", SubTest: "trailing sep", Line: 31},
		},
	}))

	t.Run("direct_range_slice_test", performTest(testCase{
		location: "./examples/direct_range_slice_test.go",
		expectedResults: []testInfo{
			{TestFunc: "TestRngSplitSliceTest", SubTest: "", Line: 10},
			{TestFunc: "TestRngSplitSliceTest", SubTest: "simple", Line: 18},
			{TestFunc: "TestRngSplitSliceTest", SubTest: "wrong sep", Line: 24},
			{TestFunc: "TestRngSplitSliceTest", SubTest: "no sep", Line: 30},
			{TestFunc: "TestRngSplitSliceTest", SubTest: "trailing sep", Line: 36},
		},
	}))

	t.Run("individual_runs_single_test", performTest(testCase{
		location: "./examples/individual_runs_single_test.go",
		expectedResults: []testInfo{
			{TestFunc: "TestSplitIndividualSingleTest", SubTest: "", Line: 10},
			{TestFunc: "TestSplitIndividualSingleTest", SubTest: "simple", Line: 26},
			{TestFunc: "TestSplitIndividualSingleTest", SubTest: "wrong sep", Line: 32},
			{TestFunc: "TestSplitIndividualSingleTest", SubTest: "no sep", Line: 38},
			{TestFunc: "TestSplitIndividualSingleTest", SubTest: "trailing sep", Line: 44},
		},
	}))

	t.Run("individual_runs_test", performTest(testCase{
		location: "./examples/individual_runs_test.go",
		expectedResults: []testInfo{
			{TestFunc: "TestSplitIndivual", SubTest: "", Line: 8},
			{TestFunc: "TestSplitIndivual", SubTest: "simple", Line: 11},
			{TestFunc: "TestSplitIndivual", SubTest: "wrong sep", Line: 19},
			{TestFunc: "TestSplitIndivual", SubTest: "no sep", Line: 27},
			{TestFunc: "TestSplitIndivual", SubTest: "trailing sep", Line: 35},
		},
	}))

	t.Run("individual_test", performTest(testCase{
		location: "./examples/individual_test.go",
		expectedResults: []testInfo{
			{TestFunc: "TestSplitSimple", SubTest: "", Line: 8},
			{TestFunc: "TestSplitWrongSep", SubTest: "", Line: 16},
			{TestFunc: "TestSplitNoSep", SubTest: "", Line: 24},
			{TestFunc: "TestSplitTrailingSep", SubTest: "", Line: 32},
		},
	}))

	t.Run("map_test", performTest(testCase{
		location: "./examples/map_test.go",
		expectedResults: []testInfo{
			{TestFunc: "TestSplitMapTest", SubTest: "", Line: 10},
			{TestFunc: "TestSplitMapTest", SubTest: "simple", Line: 16},
			{TestFunc: "TestSplitMapTest", SubTest: "wrong sep", Line: 21},
			{TestFunc: "TestSplitMapTest", SubTest: "no sep", Line: 26},
			{TestFunc: "TestSplitMapTest", SubTest: "trailing sep", Line: 31},
		},
	}))

	t.Run("slice_test", performTest(testCase{
		location: "./examples/slice_test.go",
		expectedResults: []testInfo{
			{TestFunc: "TestSplitSliceTest", SubTest: "", Line: 10},
			{TestFunc: "TestSplitSliceTest", SubTest: "simple", Line: 18},
			{TestFunc: "TestSplitSliceTest", SubTest: "wrong sep", Line: 24},
			{TestFunc: "TestSplitSliceTest", SubTest: "no sep", Line: 30},
			{TestFunc: "TestSplitSliceTest", SubTest: "trailing sep", Line: 36},
		},
	}))
}

func errStr(err error) string {
	if err == nil {
		return "[nil]"
	}
	return err.Error()
}

func runTest(t *testing.T, location string) testOutput {
	oldStdout := os.Stdout
	defer func() { os.Stdout = oldStdout }()

	// Create a pipe
	r, w, _ := os.Pipe()
	os.Stdout = w

	runErr := run(location)
	w.Close()

	// Read output
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatal("Failed to copy response data: " + err.Error())
	}

	to := testOutput{err: runErr}
	if runErr != nil {
		return to
	}

	if err := json.Unmarshal(buf.Bytes(), &to.results); err != nil {
		t.Fatal("Failed to read bytes out from response: " + err.Error())
	}

	return to
}

type testOutput struct {
	results []testInfo
	err     error
}
