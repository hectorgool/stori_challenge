package csv

import (
	"path/filepath"
	"testing"
)

// testPair defines a structure for holding test case information,
// including the file path and whether an error is expected.
type testPair struct {
	filePath      string // Path to the CSV file
	expectedError bool   // Indicates if an error is expected for this test case
}

// List of test cases with corresponding expected outcomes
var tests = []testPair{
	{"test1.csv", true},
	{"test2.csv", true},
	{"test3.csv", true},
	{"test4.csv", true},
	{"no_existe.csv", true}, // Expecting an error for a non-existent file
}

// TestProcessCSVFile tests the ProcessCSVFile function with various inputs
func TestProcessCSVFile(t *testing.T) {
	for _, pair := range tests {
		// Construct the full file path
		csvFile := filepath.Join(".", pair.filePath)

		// Call the ProcessCSVFile function
		err := ProcessCSVFile(csvFile)

		// Check if the error result matches the expected outcome
		if (err != nil) != pair.expectedError {
			// Log an error if the actual result does not match the expected result
			t.Errorf("For file %s, expected error: %v, got: %v", pair.filePath, pair.expectedError, err != nil)
		}
	}
}
