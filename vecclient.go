package cui2vec

import (
	"net/rpc"
)

type VecClient struct {
	client *rpc.Client
}

type VecResponse struct {
	V []float64
}

type SimResponse struct {
	V []Concept
}

func NewVecClient(addr string) (*VecClient, error) {
	client, err := rpc.Dial("tcp", "localhost:8003")
	if err != nil {
		return nil, err
	}

	return &VecClient{
		client: client,
	}, nil
}

func (c *VecClient) Vec(cui string) ([]float64, error) {
	vec := new(VecResponse)
	err := c.client.Call("EmbeddingsRPC.GetVector", cui, vec)
	return vec.V, err
}

func (c *VecClient) Sim(cui string) ([]Concept, error) {
	vec := new(SimResponse)
	err := c.client.Call("EmbeddingsRPC.GetSimilar", cui, vec)
	return vec.V, err
}
