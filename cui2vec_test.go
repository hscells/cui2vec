package cui2vec_test

import (
	"fmt"
	"github.com/hscells/cui2vec"
	"os"
	"testing"
)

func TestUncompressed(t *testing.T) {
	f, err := os.Open("cui2vec_pretrained_medium.csv")
	if err != nil {
		t.Fatal(err)
	}

	v, err := cui2vec.NewUncompressedEmbeddings(f, true)
	if err != nil {
		t.Fatal(err)
	}

	s, err := v.Similar("C0000052")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(s[:5])
	return
}

func TestPrecomputed(t *testing.T) {
	f, err := os.Open("cui2vec_precomputed_medium.bin")
	if err != nil {
		t.Fatal(err)
	}

	v, err := cui2vec.NewPrecomputedEmbeddings(f)
	if err != nil {
		t.Fatal(err)
	}

	s, err := v.Similar("C0000052")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(s[:5])
	return
}
