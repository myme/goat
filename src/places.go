// Description: Fetches places from the Geonorge API.
//
// Fetches "points of interest" (POI) surrounding a given location. The
// relevant fields are extracted from the JSON response:
//
// - Type: Enumeration of POI types, e.g. "Gard" (farm) or "Bruk" (small farm).
// - Distance: Distance in meters from the given location.
// - Name: The name of the POI.
// - Loc: The location of the POI.

package goat

import (
	"cmp"
	"encoding/json"
	"fmt"
	"io"
	"slices"
)

const PLACES_SEARCH_BASE = "https://ws.geonorge.no/stedsnavn/v1/punkt"

// A Place - the bits of data we care about from the places API call
type Place struct {
	Type     string
	Distance float64
	Name     string
	Loc      Location
}

func (place Place) CouldHaveGoats() bool {
	return place.Type == "Gard" || place.Type == "Bruk"
}

func SortPlaces(places []Place) {
	slices.SortFunc(places, func(a, b Place) int {
		if a.Type == b.Type {
			return cmp.Compare(a.Distance, b.Distance)
		}
		if a.Type == "Gard" {
			return -1
		}
		return 1
	})
}

func Places(loc Location) chan Result[[]Place] {
	// TODO Build query string using query string builder?
	page := 1
	hitsPerPage := 500
	// Upper limit. Consider to fetch multiple
	radius := 5000
	url := fmt.Sprintf(
		"%s?nord=%f&ost=%f&treffPerSide=%d&koordsys=4258&radius=%d&side=%d",
		PLACES_SEARCH_BASE,
		loc.Lat,
		loc.Lon,
		hitsPerPage,
		radius,
		page,
	)

	return GetJSON(url, ParsePlaces)
}

func ParsePlaces(data io.Reader) ([]Place, error) {
	var placesJson struct {
		// TODO: In case we need pagination
		Metadata struct {
			TotalHits int `json:"totaltAntallTreff"`
			Page      int `json:"side"`
		} `json:"metadata"`
		// The actual places
		Places []struct {
			Type     string  `json:"navneobjekttype"`
			Distance float64 `json:"meterFraPunkt"`
			Loc      struct {
				Lat float64 `json:"nord"`
				Lon float64 `json:"øst"`
			} `json:"representasjonspunkt"`
			Names []struct {
				Name string `json:"skrivemåte"`
			} `json:"stedsnavn"`
		} `json:"navn"`
	}

	err := json.NewDecoder(data).Decode(&placesJson)
	if err != nil {
		return nil, err
	}

	places := make([]Place, len(placesJson.Places))
	for i, place := range placesJson.Places {
		name := place.Names[0].Name
		places[i] = Place{
			Type:     place.Type,
			Distance: place.Distance,
			Name:     name,
			Loc:      Location{place.Loc.Lat, place.Loc.Lon},
		}
	}

	return places, nil
}
