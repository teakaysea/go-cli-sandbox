package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func filterOut(path, ext string, minSize int64, info os.FileInfo) bool {
	if info.IsDir() || info.Size() < minSize {
		return true
	}

	return ext != "" && filepath.Ext(path) != ext
}

func listFile(path string, out io.Writer) error {
	_, err := fmt.Fprintln(out, path)
	return err
}

func delFile(path string, delLogger *log.Logger) error {
	if err := os.Remove(path); err != nil {
		return err
	}

	delLogger.Println(path)
	return nil
}
