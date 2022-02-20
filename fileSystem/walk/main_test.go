package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestRun(t *testing.T) {
	testCases := []struct {
		name     string
		root     string
		cfg      config
		expected string
	}{
		{"NoFilter", "testdata", config{ext: "", size: 0, list: true}, "testdata/dir.log\ntestdata/dir2/script.sh\n"},
		{"FilterExtensionMatch", "testdata", config{ext: ".log", size: 0, list: true}, "testdata/dir.log\n"},
		{"FilterExtensionSizeMatch", "testdata", config{ext: ".log", size: 10, list: true}, "testdata/dir.log\n"},
		{"FilterExtensionSizeNoMatch", "testdata", config{ext: ".log", size: 20, list: true}, ""},
		{"FilterExtensionNoMatch", "testdata", config{ext: ".gz", size: 0, list: true}, ""},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer

			if err := run(tc.root, &buf, tc.cfg); err != nil {
				t.Fatal(err)
			}

			res := buf.String()

			if tc.expected != res {
				t.Errorf("Expected %q, go %q instead\n", tc.expected, res)
			}
		})
	}
}

func TestRunDelExtension(t *testing.T) {
	testCases := []struct {
		name        string
		cfg         config
		extNoDelete string
		nDelete     int
		nNoDelete   int
		expected    string
	}{
		{},
		{"DeleteExtensionNoMatch", config{ext: ".log", del: true}, ".gz", 0, 10, ""},
		{"DeleteExtensionMatch", config{ext: ".log", del: true}, "", 10, 0, ""},
		{"DeleteExtensionMixed", config{ext: ".log", del: true}, ".gz", 5, 5, ""}}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer

			tempDir, cleanup := createTempDir(t, map[string]int{
				tc.cfg.ext:     tc.nDelete,
				tc.extNoDelete: tc.nNoDelete,
			})
			defer cleanup()

			if err := run(tempDir, &buf, tc.cfg); err != nil {
				t.Fatal(err)
			}

			res := buf.String()

			if tc.expected != res {
				t.Errorf("Expected %q, got %q instead\n", tc.expected, res)
			}

			filesLeft, err := os.ReadDir(tempDir)
			if err != nil {
				t.Error(err)
			}

			if len(filesLeft) != tc.nNoDelete {
				t.Errorf("Expected %d files left, go %d instead\n", tc.nNoDelete, len(filesLeft))
			}
		})
	}
}
func createTempDir(t *testing.T, files map[string]int) (dirname string, cleanup func()) {
	t.Helper()

	tempDir, err := os.MkdirTemp("", "walktest")
	if err != nil {
		t.Fatal(err)
	}

	for k, n := range files {
		for i := 1; i <= n; i++ {
			fname := fmt.Sprintf("file%d%s", i, k)
			fpath := filepath.Join(tempDir, fname)
			if err := os.WriteFile(fpath, []byte("dummy"), 0644); err != nil {
				t.Fatal(err)
			}
		}
	}
	return tempDir, func() { os.RemoveAll(tempDir) }
}
