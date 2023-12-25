package main

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var lineCases = []struct {
	testName string
	line     string
	name     string
	duration float64
}{
	{"empty line", "", "", 0},
	{
		"benchmem",
		`BenchmarkRand-12       	38350470	        29.16 ns/op	       0 B/op	       0 allocs/op`,
		"BenchmarkRand-12",
		29.16,
	},
}

func Test_parseLine(t *testing.T) {
	for _, tc := range lineCases {
		t.Run(tc.testName, func(t *testing.T) {
			n, d := parseLine(tc.line)
			require.Equal(t, tc.name, n)
			require.Equal(t, tc.duration, d)
		})
	}
}

func openFile(t *testing.T, path string) *os.File {
	file, err := os.Open(path)
	require.NoError(t, err)
	t.Cleanup(func() { file.Close() })
	return file
}

func Test_parseFile(t *testing.T) {
	file := openFile(t, "testdata/wrand.txt")

	ds, err := parseFile(file)
	require.NoError(t, err, "parse")

	expected := map[string][]float64{
		"BenchmarkRand-12":    {29.16, 29.1, 29.03, 29.2, 26.21},
		"BenchmarkRandBig-12": {128.7, 130, 128.8, 128.4, 128.8},
	}
	require.Equal(t, expected, ds)
}

func Test_parseConfig(t *testing.T) {
	file := openFile(t, "testdata/wrand.yml")
	out, err := parseConfig(file)
	require.NoError(t, err)

	expected := map[string]time.Duration{
		"BenchmarkRand-12":    30 * time.Nanosecond,
		"BenchmarkRandBig-12": 112 * time.Nanosecond,
	}

	require.Equal(t, expected, out)
}
