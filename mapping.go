package cui2vec

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Mapping map[string]string

type AliasMapping map[string][]string

type frequency struct {
	cui       string
	term      string
	frequency int
}

// LoadCUIMapping loads a mapping of cui to most common title.
//
// Mapping of cuis->title is constructed as per:
// Jimmy, Zuccon G., Koopman B. (2018) Choices in Knowledge-Base Retrieval for Consumer Health Search.
// In: Pasi G., Piwowarski B., Azzopardi L., Hanbury A. (eds) Advances in Information Retrieval. ECIR 2018.
// Lecture Notes in Computer Science, vol 10772. Springer, Cham
//
// File must reflect this.
func LoadCUIMapping(path string) (Mapping, error) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	mapping := make(Mapping)
	for _, record := range records {
		if len(record) < 2 {
			continue
		}
		cui, title := record[0], record[1]

		for len(cui) < 7 {
			cui = "0" + cui
		}
		cui = "C" + cui

		mapping[cui] = title
	}
	return mapping, nil
}

func LoadCUIFrequencyMapping(path string) (Mapping, error) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	frequencies := make([]frequency, len(records))
	for i, record := range records {
		if len(record) < 3 {
			continue
		}
		cui, title, f := record[0], record[1], record[3]

		freq, err := strconv.Atoi(strings.Replace(f, `"`, "", -1))
		if err != nil {
			return nil, err
		}

		for len(cui) < 7 {
			cui = "0" + cui
		}
		cui = "C" + cui

		frequencies[i] = frequency{
			cui:       cui,
			term:      title,
			frequency: freq,
		}
	}

	sort.Slice(frequencies, func(i, j int) bool {
		return frequencies[i].frequency > frequencies[j].frequency
	})

	fmt.Println(frequencies[0:10])

	mapping := make(Mapping)
	for _, f := range frequencies {
		if _, ok := mapping[f.cui]; !ok {
			mapping[f.cui] = f.term
		}
	}

	return mapping, nil
}

func LoadCUIAliasMapping(path string) (AliasMapping, error) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	mapping := make(AliasMapping)
	for _, record := range records {
		if len(record) < 3 {
			continue
		}
		cui, term := record[0], record[1]

		for len(cui) < 7 {
			cui = "0" + cui
		}
		cui = "C" + cui

		mapping[cui] = append(mapping[cui], term)
	}

	return mapping, nil
}

func (m Mapping) Invert() Mapping {
	i := make(Mapping)
	for k, v := range m {
		i[v] = k
	}
	return i
}
