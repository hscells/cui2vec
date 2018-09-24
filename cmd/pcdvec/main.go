package main

import (
	"fmt"
	"github.com/alexflint/go-arg"
	"github.com/hscells/cui2vec"
	"io"
	"math"
	"os"
	"sort"
	"sync"
)

const (
	part = 10000000
)

type args struct {
	CUI       string `arg:"required" help:"path to cui2vec model"`
	Output    string `arg:"-o" help:"where to output distances to (default stdout)"`
	SkipFirst bool   `help:"skip first line in cui2vec model?"`
}

func (args) Version() string {
	return "pcdvec 24.Sep.2018"
}

func (args) Description() string {
	return `pre-compute distances for cui2vec`
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func distance(embeddings map[string][]float64) ([][]int, error) {
	defer fmt.Println()
	fmt.Println("computing distances")
	var (
		max    int
		matrix [][]int
	)

	// First, find the largest cui.
	for cui := range embeddings {
		c, err := cui2vec.CUI2Int(cui)
		if err != nil {
			return nil, err
		}
		if c > max {
			max = c
		}
	}

	fmt.Println("maximum size of sparse matrix:", max)

	// Now, create a sparse matrix of size `max`.
	matrix = make([][]int, max+1)

	// Next, compute the distances between each cui and every other cui.
	for c1, e1 := range embeddings {
		fmt.Print(".")
		// Minus one because we don't include this cui.
		concepts := make([]cui2vec.Concept, len(embeddings)-1)
		i := 0

		var wg sync.WaitGroup
		for c2, e2 := range embeddings {
			// Don't compute distances for the same cui.
			if c1 == c2 {
				continue
			}

			if len(c2) == 0 {
				continue
			}
			wg.Add(1)
			go func(idx int, v1, v2 []float64, concept string) {
				defer wg.Done()
				sim, err := cui2vec.Cosine(v1, v2)
				if err != nil {
					panic(err)
				}
				concepts[idx] = cui2vec.Concept{
					CUI:   concept,
					Value: sim,
				}
			}(i, e1, e2, c2)
			i++
		}

		wg.Wait()

		// Normalise the concepts.
		concepts = cui2vec.Softmax(concepts)

		// Sort them by score.
		sort.Slice(concepts, func(i, j int) bool {
			return concepts[i].Value > concepts[j].Value
		})

		// Take the top n concepts (in this case n=20; hardcoded).
		n := 20
		if len(concepts) < n {
			n = len(concepts)
		}
		concepts = concepts[:n]

		// Convert the string value of the cui into an int.
		c, err := cui2vec.CUI2Int(c1)
		if err != nil {
			return nil, err
		}

		// Create the row for the current concept in the matrix.
		matrix[c] = make([]int, len(concepts)*2)

		// Insert the concepts one by one into the row.
		for i, j := 0, 0; i < len(concepts); i++ {
			cc, err := cui2vec.CUI2Int(concepts[i].CUI)
			if err != nil {
				return nil, err
			}

			v := int(round(concepts[i].Value * part))
			v = v - v/part*part

			matrix[c][j] = cc
			matrix[c][j+1] = v
			j += 2
		}
	}
	return matrix, nil
}

func main() {
	var (
		args   args
		err    error
		input  io.ReadCloser
		output io.WriteCloser
	)
	arg.MustParse(&args)

	// Open the output file, defaulting to stdout.
	if len(args.Output) == 0 {
		output = os.Stdout
	} else {
		output, err = os.OpenFile(args.Output, os.O_WRONLY|os.O_CREATE, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	defer output.Close()

	// Open the input file to read the embeddings.
	input, err = os.OpenFile(args.CUI, os.O_RDONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer input.Close()

	// Load the embeddings into memory.
	ue, err := cui2vec.NewUncompressedEmbeddings(input, args.SkipFirst)
	if err != nil {
		panic(err)
	}

	// Create a new pre-computed embeddings with distance calculations.
	m, err := distance(ue.Embeddings)
	if err != nil {
		panic(err)
	}
	pe := cui2vec.PrecomputedEmbeddings{
		Matrix: m,
		Cols:   20,
	}

	// Output the pre-computed distances to file.
	err = pe.WriteModel(output)
	if err != nil {
		panic(err)
	}

	return
}
