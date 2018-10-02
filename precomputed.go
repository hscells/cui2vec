package cui2vec

import (
	"encoding/binary"
	"io"
	"io/ioutil"
	"math"
)

// PrecomputedEmbeddings is a type of cui2vec container where the distances between CUIs have been pre-computed.
// It contains a sparse Matrix where the rows are CUIs and the columns are the distances to other CUIs.
// Each row is formatted in the form [CUI, score, CUI, score, ...].
// Each CUI must be converted back to a string, and each score must be re-normalised from an int back to a float (taken care of by the Similar method).
type PrecomputedEmbeddings struct {
	Matrix [][]int
	Cols   int
}

// LoadModel reads a model from disk into memory. The file format of the pre-computed distances
// file is that of a single, continuous byte sequence starting with four bytes indicating the rows in the matrix.
// The first four bytes indicate a single Uint32 number representing the size of the matrix.
// This is used to create a fixed-size sparse matrix. The `Cols` attribute of the `PrecomputedEmbeddings` type
// is used to read N four-byte Uint32 numbers at a time to populate the columns of the matrix.
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

// WriteModel writes a pre-computed distance matrix to disk.
// The write begins with a four-byte sequence to be parsed as a Uint32 representing the
// size of the matrix. Each value of the matrix is then written one by one in a continuous
// byte sequence, where each element in the matrix is encoded as a four-byte sequence to
// be parsed as a Uint32. Elements of the matrix are written row-by-row, and each row
// is exactly `Cols` wide. If there are less than `Cols` elements in a row, the row is
// padded with zeros.
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

// Similar matches a given input CUI to the `Cols`-closest CUIs in the cui2vec embedding space.
// As each row in the matrix is encoded into (CUI, score) pairs, this method handles that.
// It also converts each int value in the matrix into either a string CUI or a re-normalised softmax score float64.
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

	// Exit early if the CUI is malformed.
	if c >= len(v.Matrix) || c < 0 {
		return concepts, nil
	}

	// Create a slice of concepts from the number of pre-computed scores.
	concepts = make([]Concept, v.Cols/2)
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
			concepts[j] = c
			j++
		} else {
			concept = Int2CUI(val)
		}
	}

	return concepts, nil
}
