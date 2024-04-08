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

type AdressSearchResponse struct {
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
	return GetJSON(url, ParseAddress)
}

func ParseAddress(data io.Reader) ([]Address, error) {
	var addressResponse struct {
		Addresses []Address `json:"adresser"`
	}

	err := json.NewDecoder(data).Decode(&addressResponse)
	if err != nil {
		return nil, err
	}

	return addressResponse.Addresses, nil
}

func (a Address) Format() string {
	return fmt.Sprintf("🏠 %s, %s %s", a.Text, a.PostCode, a.PostText)
}
