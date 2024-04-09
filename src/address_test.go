package goat

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestParseAddress(t *testing.T) {
	// API JSON response data
	jsonData := `
		{
		  "adresser": [
			{
			  "adressekode": 25440,
			  "adressenavn": "Myrvollveien",
			  "adressetekst": "Myrvollveien 5C",
			  "adressetekstutenadressetilleggsnavn": "Myrvollveien 5C",
			  "adressetilleggsnavn": null,
			  "bokstav": "C",
			  "bruksenhetsnummer": [],
			  "bruksnummer": 898,
			  "festenummer": 0,
			  "gardsnummer": 243,
			  "kommunenavn": "NORDRE FOLLO",
			  "kommunenummer": "3207",
			  "nummer": 5,
			  "objtype": "Vegadresse",
			  "oppdateringsdato": "2024-01-01T00:00:00",
			  "postnummer": "1415",
			  "poststed": "OPPEGÅRD",
			  "representasjonspunkt": {
				"epsg": "EPSG:4258",
				"lat": 59.78502106569645,
				"lon": 10.799290993113777
			  },
			  "stedfestingverifisert": true,
			  "undernummer": null
			}
		  ],
		  "metadata": {
			"asciiKompatibel": true,
			"side": 0,
			"sokeStreng": "sok=myrvollveien%205c",
			"totaltAntallTreff": 1,
			"treffPerSide": 10,
			"viserFra": 0,
			"viserTil": 10
		  }
		}
	`

	parsed, _ := ParseAddress(1, strings.NewReader(jsonData))

	expected := Address{
		Text:     "Myrvollveien 5C",
		PostCode: "1415",
		PostText: "OPPEGÅRD",
		Loc:      Location{59.78502106569645, 10.799290993113777},
	}

	if len(parsed.Addresses) != 1 || parsed.Addresses[0] != expected {
		t.Errorf("Expected %v address, got %v", expected, parsed)
	}
}

func TestFetchAllPages(t *testing.T) {
	hitsPerPage := 1

	// API JSON response data
	makeJsonData := func(page, totalHits int) string {
		from := page * hitsPerPage
		to := from + hitsPerPage
		return fmt.Sprintf(`
			{
			  "adresser": [
				{
				  "adressetekst": "Myrvollveien %d",
				  "postnummer": "1415",
				  "poststed": "OPPEGÅRD",
				  "representasjonspunkt": {
					"epsg": "EPSG:4258",
					"lat": 59.78502106569645,
					"lon": 10.799290993113777
				  }
				}
			  ],
			  "metadata": {
				"asciiKompatibel": true,
				"side": %d,
				"sokeStreng": "sok=myrvollveien%%205c",
				"totaltAntallTreff": %d,
				"treffPerSide": %d,
				"viserFra": %d,
				"viserTil": %d
			  }
			}
		`, page+1, page, totalHits, hitsPerPage, from, to)
	}

	// Synthesize fetching all pages
	fetchAllPages := func(totalHits int, err error) chan Result[[]Address] {
		return FetchAllPages(func(page int) chan Result[*AddressSearchResponse] {
			ch := make(chan Result[*AddressSearchResponse])

			go func() {
				if err != nil {
					ch <- Result[*AddressSearchResponse]{Ok: nil, Err: err}
					return
				}
				res, err := ParseAddress(hitsPerPage, strings.NewReader(makeJsonData(page, totalHits)))
				ch <- Result[*AddressSearchResponse]{Ok: &res, Err: err}
			}()

			return ch
		})
	}

	t.Run("Error", func(t *testing.T) {
		err := fmt.Errorf("Something wrong")
		totalHits := 1
		result := <-fetchAllPages(totalHits, err)

		if result.Err != err {
			t.Errorf("Expected %v address, got %v", err, result.Err)
		}
	})

	t.Run("Single page", func(t *testing.T) {
		totalHits := 1
		result := <-fetchAllPages(totalHits, nil)

		fetched := *result.Ok
		expected := []Address{
			{
				Text:     "Myrvollveien 1",
				PostCode: "1415",
				PostText: "OPPEGÅRD",
				Loc:      Location{59.78502106569645, 10.799290993113777},
			},
		}

		if !reflect.DeepEqual(fetched, expected) {
			t.Errorf("Expected %v address, got %v", expected, fetched)
		}
	})

	t.Run("Three pages", func(t *testing.T) {
		totalHits := 3
		result := <-fetchAllPages(totalHits, nil)

		fetched := *result.Ok
		expected := []Address{
			{
				Text:     "Myrvollveien 1",
				PostCode: "1415",
				PostText: "OPPEGÅRD",
				Loc:      Location{59.78502106569645, 10.799290993113777},
			},
			{
				Text:     "Myrvollveien 2",
				PostCode: "1415",
				PostText: "OPPEGÅRD",
				Loc:      Location{59.78502106569645, 10.799290993113777},
			},
			{
				Text:     "Myrvollveien 3",
				PostCode: "1415",
				PostText: "OPPEGÅRD",
				Loc:      Location{59.78502106569645, 10.799290993113777},
			},
		}

		if !reflect.DeepEqual(fetched, expected) {
			t.Errorf("Expected %v address, got %v", expected, fetched)
		}
	})
}
