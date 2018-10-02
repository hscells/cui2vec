package cui2vec

import (
	"testing"
	"fmt"
)

func TestMapping(t *testing.T) {
	m, err := LoadCUIMapping("cuis.csv")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(m["C0000294"])
}
