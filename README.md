# cui2vec

package cui2vec implements utilities for dealing with cui2vec embeddings and mapping cuis to text.

Paper (author not affiliated with this code): https://arxiv.org/pdf/1804.01486.pdf

Documentation: https://godoc.org/github.com/hscells/cui2vec

Example: See [cmd/cui2vec/main.go](cmd/cui2vec/main.go)

## Data

Pre-trained embeddings (model) can be downloaded from https://figshare.com/s/00d69861786cd0156d81.

A pre-computed distances version of these pre-trained embeddings is included in the [testdata](testdata) folder. 

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

```bash
Usage: cui2vec [--cui CUI] [--model MODEL] [--type TYPE] [--skipfirst] [--numcuis NUMCUIS] [--mapping MAPPING] [--verbose]

Options:
  --cui CUI
  --model MODEL
  --type TYPE
  --skipfirst
  --numcuis NUMCUIS, -n NUMCUIS
  --mapping MAPPING
  --verbose, -v
  --help, -h             display this help and exit
  --version              display version and exit
```

### Pre-computing distances

A tool that can be used to compress and increase the speed of computing similar CUIs is included in the form of `pcdvec`.

Install the tool with 

```bash
go get github.com/hscells/cui2vec/cmd/pcdvec
```

```bash
Usage: pcdvec --cui CUI [--output OUTPUT] [--skipfirst]

Options:
  --cui CUI
  --output OUTPUT, -o OUTPUT
  --skipfirst
  --help, -h             display this help and exit
  --version              display version and exit
```