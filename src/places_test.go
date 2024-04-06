package goat

import (
	"reflect"
	"strings"
	"testing"
)

func TestParsePlaces(t *testing.T) {
	// Single-entry API response
	jsonData := `
		{
			"metadata": {
				"side": 1,
				"sokeStreng": "nord=59.900056&ost=10.590543&treffPerSide=1&koordsys=4258&radius=5000&side=1",
				"totaltAntallTreff": 1345,
				"treffPerSide": 1,
				"viserFra": 1,
				"viserTil": 1
			},
			"navn": [
				{
					"meterFraPunkt": 1991,
					"navneobjekttype": "Holdeplass",
					"representasjonspunkt": {
						"koordsys": 4258,
						"nord": 59.91168,
						"øst": 10.61755
					},
					"stedsnavn": [
						{
							"navnestatus": "historisk",
							"skrivemåte": "Myra holdeplass",
							"skrivemåtestatus": "godkjent og prioritert",
							"språk": "Norsk",
							"stedsnavnnummer": 1
						}
					],
					"stedsnummer": 550905,
					"stedstatus": "relikt"
				}
			]
		}
	`

	places, _ := ParsePlaces(strings.NewReader(jsonData))

	expected := Place{
		Type:     "Holdeplass",
		Distance: 1991,
		Name:     "Myra holdeplass",
		Loc:      Location{59.91168, 10.61755},
	}

	if len(places) != 1 || places[0] != expected {
		t.Errorf("Expected %v, got %v", expected, places)
	}
}

func TestSortPlaces(t *testing.T) {
	places := []Place{
		{Distance: 100, Type: "Bruk"},
		{Distance: 10, Type: "Bruk"},
		{Distance: 1000, Type: "Gard"},
		{Distance: 100, Type: "Gard"},
	}

	SortPlaces(places)

	expected := []Place{
		{Distance: 100, Type: "Gard"},
		{Distance: 1000, Type: "Gard"},
		{Distance: 10, Type: "Bruk"},
		{Distance: 100, Type: "Bruk"},
	}

	if !reflect.DeepEqual(places, expected) {
		t.Errorf("Expected %v, got %v", expected, places)
	}
}
