package main

import (
	"github.com/alexflint/go-arg"
	"fmt"
	"github.com/hscells/cui2vec"
	"os"
	"encoding/json"
)

type args struct {
	CUI       string `arg:"help:input cui,required"`
	Model     string `arg:"help:path to cui2vec model"`
	SkipFirst bool   `arg:"help:skip first line in cui2vec model?"`
	N         int    `arg:"help:number of cuis to output"`
	Mapping   string `arg:"help:path to cui mapping"`
	V         bool   `arg:"help:verbose output"`
}

func (args) Version() string {
	return "cui2vec 21.Aug.2018"
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
		if args.V {
			fmt.Println("loading model...")
		}

		f, err := os.OpenFile(args.Model, os.O_RDONLY, os.ModePerm)
		if err != nil {
			panic(err)
		}
		e, err := cui2vec.LoadModel(f, args.SkipFirst)
		if err != nil {
			panic(err)
		}

		if args.V {
			fmt.Println("computing similarity...")
		}
		concepts, err := e.Similar(args.CUI)
		if err != nil {
			panic(err)
		}

		if args.N > 0 {
			// Resize the slice.
			if args.V {
				fmt.Println("resizing slice...")
			}

			if len(concepts) > args.N {
				concepts = concepts[:args.N]
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
		if args.V {
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
