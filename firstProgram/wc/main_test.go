package main

import (
	"bytes"
	"testing"
)

func TestCountWords(t *testing.T) {
	b := bytes.NewBufferString(`word1 word2 word3 word4
`)
	exp := 4
	res := count(b, countWords)
	if res != exp {
		t.Errorf("Expected %d, got %d instead", exp, res)
	}
}

func TestCountLines(t *testing.T) {
	b := bytes.NewBufferString(`word1 word2 word3
line2
line3 word1`)

	exp := 3
	res := count(b, countLines)

	if res != exp {
		t.Errorf("Expected %d, go %d instead", exp, res)
	}

}
func TestCountBytes(t *testing.T) {
	b := bytes.NewBufferString(`word1 word2`)

	exp := 11
	res := count(b, countBytes)

	if res != exp {
		t.Errorf("Expected %d, go %d instead", exp, res)
	}

}
