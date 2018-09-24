package main

import (
	"encoding/json"
	"fmt"
	"github.com/alexflint/go-arg"
	"github.com/go-errors/errors"
	"github.com/hscells/cui2vec"
	"os"
)

type args struct {
	CUI       string `help:"input cui,required"`
	Model     string `help:"path to cui2vec model"`
	Type      string `help:"what kind of cui2vec model is loaded (default/precomputed)"`
	SkipFirst bool   `help:"skip first line in cui2vec model?"`
	NumCUIS   int    `arg:"-n" help:"number of cuis to output"`
	Mapping   string `help:"path to cui mapping"`
	Verbose   bool   `arg:"-v" help:"verbose output"`
}

func (args) Version() string {
	return "cui2vec 24.Sep.2018"
}

func (args) Description() string {
	return `tool for using cui2vec
paper: https://arxiv.org/pdf/1804.01486.pdf
repository: https://github.com/hscells/cui2vec

the author of this program is not affiliated with the authors of the paper`
}

func main() {
	var args args
	arg.MustParse(&args)

	if len(args.Model) > 0 {
		if args.Verbose {
			fmt.Println("loading model...")
		}

		f, err := os.OpenFile(args.Model, os.O_RDONLY, os.ModePerm)
		if err != nil {
			panic(err)
		}

		var e cui2vec.Embeddings
		if args.Type == "default" {
			e, err = cui2vec.NewUncompressedEmbeddings(f, args.SkipFirst)
			if err != nil {
				panic(err)
			}
		} else if args.Type == "precomputed" {
			e, err = cui2vec.NewPrecomputedEmbeddings(f)
			if err != nil {
				panic(err)
			}
		} else {
			panic(errors.New("unrecognised model type"))
		}

		if args.Verbose {
			fmt.Println("computing similarity...")
		}
		concepts, err := e.Similar(args.CUI)
		if err != nil {
			panic(err)
		}

		if args.NumCUIS > 0 {
			// Resize the slice.
			if args.Verbose {
				fmt.Println("resizing slice...")
			}

			if len(concepts) > args.NumCUIS {
				concepts = concepts[:args.NumCUIS]
			}
		}

		b, err := json.Marshal(concepts)
		if err != nil {
			panic(err)
		}

		_, err = os.Stdout.Write(b)
		if err != nil {
			panic(err)
		}
		return
	}

	if len(args.Mapping) > 0 {
		if args.Verbose {
			fmt.Println("loading mapping...")
		}
		m, err := cui2vec.LoadCUIMapping(args.Mapping)
		if err != nil {
			panic(err)
		}

		_, err = os.Stdout.Write([]byte(m[args.CUI]))
		if err != nil {
			panic(err)
		}
		return
	}

	fmt.Println("please provide an argument, use --help for help")
	return
}
