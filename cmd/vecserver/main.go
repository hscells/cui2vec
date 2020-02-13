package main

import (
	"fmt"
	"github.com/alexflint/go-arg"
	"github.com/hscells/cui2vec"
	"net"
	"net/rpc"
	"os"
	"time"
)

type args struct {
	CUI       string `arg:"required" help:"path to uncompressed model"`
	Delimiter rune   `help:"What is the delimiter (default:' ')"`
	SkipFirst bool   `help:"skip first line in cui2vec model?"`
}

func (args) Version() string {
	return "vecserver 03.Oct.2019"
}

func (args) Description() string {
	return `vector server for fast access to elements`
}

type similarCache map[string][]cui2vec.Concept

type EmbeddingsRPC struct {
	embeddings *cui2vec.UncompressedEmbeddings
	cache      similarCache
}

func logf(message string, args ...interface{}) {
	fmt.Printf(fmt.Sprintf("[%s] %s\n", time.Now().Format(time.RFC3339), fmt.Sprintf(message, args...)))
}

func (e *EmbeddingsRPC) GetVector(cui string, vec *cui2vec.VecResponse) error {
	if v, ok := e.embeddings.Embeddings[cui]; ok {
		vec.V = v
		logf("request for %s, found: %d", cui, len(vec.V))
		return nil
	}
	logf("request for %s, found nothing", cui)
	return nil
}

func (e *EmbeddingsRPC) GetSimilar(cui string, vec *cui2vec.SimResponse) error {
	if v, ok := e.cache[cui]; ok {
		vec.V = v
		return nil
	}
	logf("request for %s", cui)
	v, err := e.embeddings.Similar(cui)
	if err != nil {
		e.cache[cui] = v
	}
	vec.V = v
	return err
}

func main() {
	var args args
	arg.MustParse(&args)

	logf("initialising server...")
	addy, err := net.ResolveTCPAddr("tcp", "0.0.0.0:8003")
	if err != nil {
		panic(err)
	}

	inbound, err := net.ListenTCP("tcp", addy)
	if err != nil {
		panic(err)
	}

	logf("loading embeddings...")
	f, err := os.OpenFile(args.CUI, os.O_RDONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}

	if args.Delimiter == 0 {
		args.Delimiter = ' '
	}

	e, err := cui2vec.NewUncompressedEmbeddings(f, args.SkipFirst, args.Delimiter)
	if err != nil {
		panic(err)
	}

	fmt.Println(e.Embeddings["C0243032"])

	logf("registering listener...")
	listener := new(EmbeddingsRPC)
	x := EmbeddingsRPC{embeddings: e, cache: make(similarCache)}
	listener = &x
	err = rpc.Register(listener)
	if err != nil {
		panic(err)
	}
	logf("ready to go!")
	rpc.Accept(inbound)

}
