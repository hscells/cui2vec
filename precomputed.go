package cui2vec

import (
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

// PrecomputedEmbeddings is a type of cui2vec container where the distances between CUIs have been pre-computed.
// It contains a sparse Matrix where the rows are CUIs and the columns are the distances to other CUIs.
// Each row is formatted in the form [CUI, score, CUI, score, ...].
// Each CUI must be converted back to a string, and each score must be re-normalised from an int back to a float (taken care of by the Similar method).
// The file format of the
type PrecomputedEmbeddings struct {
	Matrix [][]int
	Cols   int
}

func (v *PrecomputedEmbeddings) LoadModel(r io.Reader) error {
	var (
		matrix [][]int
		size   int
		idx    int
	)

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	size = int(binary.LittleEndian.Uint32(b[:4]))
	matrix = make([][]int, size)

	// i = matrix row; j = bytes position; k = matrix column.
	for i, j, k := 0, 4, 0; i < size; i++ {
		if k == 0 { // Start of new section
			idx = int(binary.LittleEndian.Uint32(b[j : j+4]))
			matrix[idx] = make([]int, v.Cols)
		} else {
			matrix[idx][k-1] = int(binary.LittleEndian.Uint32(b[j : j+4]))
		}
		k++
		if k == v.Cols+1 {
			k = 0
		}
		j += 4
	}
	v.Matrix = matrix
	return nil
}

func (v *PrecomputedEmbeddings) WriteModel(w io.Writer) error {
	d := make([]byte, 4)
	binary.LittleEndian.PutUint32(d, uint32(len(v.Matrix)))
	_, err := w.Write(d)
	if err != nil {
		return err
	}
	for i := range v.Matrix {
		if len(v.Matrix[i]) == 0 {
			continue
		}
		b := make([]byte, (v.Cols*4)+4)
		d := make([]byte, 4)
		binary.LittleEndian.PutUint32(d, uint32(i))
		b[0] = d[0]
		b[1] = d[1]
		b[2] = d[2]
		b[3] = d[3]
		for j, k := 0, 4; j < v.Cols; j++ {
			d := make([]byte, 4)
			if j < len(v.Matrix[i]) {
				binary.LittleEndian.PutUint32(d, uint32(v.Matrix[i][j]))
			} else {
				d = []byte{0, 0, 0, 0}
			}
			b[k] = d[0]
			b[k+1] = d[1]
			b[k+2] = d[2]
			b[k+3] = d[3]
			k += 4
		}
		_, err := w.Write(b)
		if err != nil {
			return err
		}
	}
	return nil
}

func (v *PrecomputedEmbeddings) Similar(cui string) ([]Concept, error) {
	c, err := CUI2Int(cui)
	if err != nil {
		return nil, err
	}

	var (
		concept  string
		score    float64
		concepts []Concept
	)

	// Create a slice of concepts from the number of pre-computed scores.
	concepts = make([]Concept, len(v.Matrix[c])/2)
	j := 0
	for i, val := range v.Matrix[c] {
		if i%2 != 0 {
			// Convert the very large int score back into the softmax score.
			// This works by finding the log-10 of the value, taking the ceiling of it, and taking the power of 10 to that value.
			// This finds the value by which the very large number should be divided by in order to convert it back to the original score.
			score = float64(val) / math.Pow(10, math.Ceil(math.Log10(float64(val))))

			c := Concept{
				CUI:   concept,
				Value: score,
			}
			fmt.Println(c)
			concepts[j] = c
			j++
		} else {
			concept = Int2CUI(val)
		}
	}

	return concepts, nil
}
