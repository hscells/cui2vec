package cui2vec

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"runtime"
	"sort"
	"strconv"
	"sync"
)

type UncompressedEmbeddings struct {
	SkipFirst  bool
	Comma      rune
	Embeddings map[string][]float64
}

// LoadModel a cui2vec pre-trained model into memory.
// The pre-trained file from:
// 	https://arxiv.org/pdf/1804.01486.pdf
// which was downloaded from:
//	https://figshare.com/s/00d69861786cd0156d81
// is a csv file. The skipFirst parameter determines if the first line of the file should be skipped.
func (v *UncompressedEmbeddings) LoadModel(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	if v.SkipFirst {
		scanner.Scan()
	}

	concurrency := runtime.NumCPU()
	var mu sync.Mutex
	queue := make(chan string)
	complete := make(chan bool)
	embeddings := make(map[string][]float64)

	// Read the pre-trained vector file line by line.
	go func() {
		for scanner.Scan() {
			queue <- scanner.Text()
		}
		close(queue)
	}()

	for i := 0; i < concurrency; i++ {
		go func(q chan string, complete chan bool) {
			for b := range q {
				// Use a csv parser to read the line.
				reader := csv.NewReader(bytes.NewBufferString(b))
				reader.Comma = v.Comma
				line, err := reader.Read()
				if err != nil {
					panic(err)
				}
				if len(line) > 0 {
					cui := line[0]
					vec := make([]float64, len(line))
					for i := 1; i < len(line); i++ {
						// The features come in as strings and must be parsed.
						vec[i], err = strconv.ParseFloat(line[i], 64)
						if err != nil {
							fmt.Println(len(line), line)
							panic(err)
						}
					}
					mu.Lock()
					embeddings[cui] = vec
					mu.Unlock()
				}
			}
			complete <- true
		}(queue, complete)
	}

	// Wait until the last goroutine has read from the semaphore.
	for i := 0; i < concurrency; i++ {
		<-complete
	}
	v.Embeddings = embeddings
	return nil
}

// Similar computes cuis that a similar to an input CUI. The distance function used is Cosine similarity. The CUIs are
// then run through Softmax and sorted.
func (v *UncompressedEmbeddings) Similar(cui string) ([]Concept, error) {
	vec := v.Embeddings[cui]
	var cuis []Concept
	i := 0

	concurrency := runtime.NumCPU() * 2
	sem := make(chan bool, concurrency)
	var mu sync.Mutex

	// Compute the Cosine similarity for each value.
	for vectorCui, vectorVector := range v.Embeddings {
		sem <- true
		go func(c string, f []float64) {
			defer func() { <-sem }()
			if c != cui {
				sim, err := Cosine(vec, f)
				if err != nil {
					return
				}

				if len(c) == 0 {
					return
				}

				mu.Lock()
				cuis = append(cuis, Concept{
					CUI:   c,
					Value: sim,
				})
				i++
				mu.Unlock()
			}
		}(vectorCui, vectorVector)
	}

	// Wait until the last goroutine has read from the semaphore.
	for i := 0; i < cap(sem); i++ {
		sem <- true
	}

	// Softmax the values.
	cuis = Softmax(cuis)

	// Sort the values.
	sort.Slice(cuis, func(i, j int) bool {
		return cuis[i].Value > cuis[j].Value
	})

	return cuis, nil
}
