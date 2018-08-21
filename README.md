# cui2vec

package cui2vec implements utilities for dealing with cui2vec embeddings and mapping cuis to text.

Paper (author not affiliated with this code): https://arxiv.org/pdf/1804.01486.pdf

Documentation: https://godoc.org/github.com/hscells/cui2vec

Example: See [cmd/cui2vec/main.go](cmd/cui2vec/main.go)

## Data

Pretrained embeddings (model) can be downloaded from https://figshare.com/s/00d69861786cd0156d81.

---

Example file structure of mapping file:

```
ICUI,title
5,(131)i-maa
39,dipalmitoylphosphatidylcholine
96,1-methyl-3-isobutylxanthine
107,"1-(n-methylglycine)-8-l-isoleucine-angiotensin ii"
139,"16,16-dimethyl-pge2"
151,"17 beta hydroxy 5 beta androstan 3 one"
167,17-ketosteroids
172,18-hydroxycorticosterone
173,"18 hydroxydesoxycorticosterone"
```

One way this can be constructed is detailed in:

```
Jimmy, Zuccon G., Koopman B. (2018) Choices in Knowledge-Base Retrieval for Consumer Health Search.
In: Pasi G., Piwowarski B., Azzopardi L., Hanbury A. (eds) Advances in Information Retrieval. ECIR 2018.
Lecture Notes in Computer Science, vol 10772. Springer, Cham
```

## Command-line

Command-line utility can be installed with:

```bash
go install github.com/hscells/cui2vec/cmd/cui2vec
```

```go
Usage: cui2vec --cui CUI [--model MODEL] [--skipfirst] [--n N] [--mapping MAPPING] [--v]

Options:
  --cui CUI              input cui
  --model MODEL          path to cui2vec model
  --skipfirst            skip first line in cui2vec model?
  --n N                  number of cuis to output
  --mapping MAPPING      path to cui mapping
  --v                    verbose output
  --help, -h             display this help and exit
  --version              display version and exit
```

