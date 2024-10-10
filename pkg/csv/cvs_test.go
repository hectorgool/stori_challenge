package csv

import (
	"path/filepath"
	"testing"
)

type testpair struct {
	filePath      string
	expectedError bool
}

var tests = []testpair{
	{"test1.csv", true},
	{"test2.csv", true},
	{"test3.csv", true},
	{"test4.csv", true},
	{"no_existe.csv", true},
}

func TestProcessCSVFile(t *testing.T) {
	for _, pair := range tests {
		csvFile := filepath.Join(".", pair.filePath)
		err := ProcessCSVFile(csvFile)
		if (err != nil) != pair.expectedError {
			t.Errorf("For file %s, expected error: %v, got: %v", pair.filePath, pair.expectedError, err != nil)
		}
	}
}
