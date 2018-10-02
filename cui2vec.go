// package cui2vec implements utilities for dealing with cui2vec Embeddings and mapping cuis to text.
package cui2vec

import (
	"fmt"
	"github.com/go-errors/errors"
	"io"
	"regexp"
	"strconv"
)

// Embeddings is a complete cui2vec file loaded into memory.
type Embeddings interface {
	LoadModel(r io.Reader) error
	Similar(cui string) ([]Concept, error)
}

// Concept is a CUI that has a similarity score in relation to a target CUI.
type Concept struct {
	CUI   string
	Value float64
}

var re, _ = regexp.Compile("C[0]*(?P<CUI>[0-9]+)")

// CUI2Int converts a string CUI into an integer.
func CUI2Int(cui string) (int, error) {
	m := re.FindAllStringSubmatch(cui, -1)
	if len(m) != 1 {
		return 0, errors.New(fmt.Sprintf("%s is not a cui", cui))
	}
	v, err := strconv.Atoi(m[0][1])
	if err != nil {
		return 0, err
	}
	return v, nil
}

// Int2CUI converts an integer value to a CUI.
func Int2CUI(val int) string {
	cui := strconv.Itoa(val)
	for len(cui) < 7 {
		cui = "0" + cui
	}
	return "C" + cui
}

func NewUncompressedEmbeddings(r io.Reader, skipFirst bool) (*UncompressedEmbeddings, error) {
	v := &UncompressedEmbeddings{
		SkipFirst:  skipFirst,
		Embeddings: make(map[string][]float64),
	}
	err := v.LoadModel(r)
	return v, err
}

func NewPrecomputedEmbeddings(r io.Reader) (*PrecomputedEmbeddings, error) {
	v := &PrecomputedEmbeddings{
		Cols: 20,
	}
	err := v.LoadModel(r)
	return v, err
}
