package cui2vec

import (
	"testing"
	"os"
	"log"
	"fmt"
)

func TestName(t *testing.T) {
	log.Println("loading file")
	f, err := os.Open("cui2vec_pretrained.csv")
	if err != nil {
		t.Fatal(err)
	}

	v, err := Load(f, true)
	if err != nil {
		t.Fatal(err)
	}
	log.Println("loaded file")

	log.Println("computing distance")
	s, err := v.Similar("C0000052")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(s[:5])
	log.Println("computed distance")
	return
}
