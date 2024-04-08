package goat

import (
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

	parsed, _ := ParseAddress(strings.NewReader(jsonData))

	expected := Address{
		Text:     "Myrvollveien 5C",
		PostCode: "1415",
		PostText: "OPPEGÅRD",
		Loc:      Location{59.78502106569645, 10.799290993113777},
	}

	if len(parsed) != 1 || parsed[0] != expected {
		t.Errorf("Expected %v address, got %v", expected, parsed)
	}
}
