package cui2vec

import (
	"encoding/csv"
	"os"
)

type Mapping map[string]string

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
			cui = "0"+cui
		}
		cui = "C" + cui

		mapping[cui] = title
	}
	return mapping, nil
}
