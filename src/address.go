// Description: Address search functionality.
//
// Search for addresses using the Geonorge APIs. Resulting JSON is parsed and
// only the relevant fields are extracted.

package goat

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
)

const ADDRESS_SEARCH_BASE = "https://ws.geonorge.no/adresser/v1/sok"

type AddressSearchResponse struct {
	Metadata struct {
		Page        int `json:"side"`
		TotalHits   int `json:"totaltAntallTreff"`
		HitsPerPage int `json:"treffPerSide"`
		From        int `json:"viserFra"`
		To          int `json:"viserTil"`
	} `json:"metadata"`
	Addresses []Address `json:"adresser"`
}

type Address struct {
	Text     string   `json:"adressetekst"`
	PostCode string   `json:"postnummer"`
	PostText string   `json:"poststed"`
	Loc      Location `json:"representasjonspunkt"`
}

func SearchAddress(query string) chan Result[[]Address] {
	hitsPerPage := 100
	url := fmt.Sprintf(
		"%s?sok=%s&treffPerSide=%d",
		ADDRESS_SEARCH_BASE,
		url.QueryEscape(query),
		hitsPerPage,
	)
	return FetchAllPages(func(page int) chan Result[*AddressSearchResponse] {
		return GetJSON(url, ParseAddress)
	})
}

func FetchAllPages(fetch func(int) chan Result[*AddressSearchResponse]) chan Result[[]Address] {
	ch := make(chan Result[[]Address])

	go func() {
		result := <-fetch(0)
		if result.Err != nil {
			ch <- Result[[]Address]{Ok: nil, Err: result.Err}
		} else {
			ch <- Result[[]Address]{Ok: &(*result.Ok).Addresses, Err: nil}
		}
	}()

	return ch
}

func ParseAddress(data io.Reader) (*AddressSearchResponse, error) {
	var addressResponse AddressSearchResponse

	err := json.NewDecoder(data).Decode(&addressResponse)
	if err != nil {
		return nil, err
	}

	return &addressResponse, nil
}

func (a Address) Format() string {
	return fmt.Sprintf("🏠 %s, %s %s", a.Text, a.PostCode, a.PostText)
}
