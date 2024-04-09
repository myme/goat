// Description: Address search functionality.
//
// Search for addresses using the Geonorge APIs. Resulting JSON is parsed and
// only the relevant fields are extracted.

package goat

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
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
	return FetchAllPages(func(page int) chan Result[*AddressSearchResponse] {
		url := fmt.Sprintf(
			"%s?sok=%s&side=%d&treffPerSide=%d",
			ADDRESS_SEARCH_BASE,
			url.QueryEscape(query),
			page,
			hitsPerPage,
		)
		return GetJSON(url, func(data io.Reader) (*AddressSearchResponse, error) {
			return ParseAddress(hitsPerPage, data)
		})
	})
}

type PageFetcher func(int) chan Result[*AddressSearchResponse]

func FetchAllPages(fetchPage PageFetcher) chan Result[[]Address] {
	ch := make(chan Result[[]Address])

	go func() {
		// Fetch the first page
		result := <-fetchPage(0)
		if result.Err != nil {
			ch <- Result[[]Address]{Ok: nil, Err: result.Err}
			return
		}

		addresses := (*result.Ok).Addresses
		meta := (*result.Ok).Metadata

		// Are there more pages to fetch?
		if meta.To < meta.TotalHits {
			// Start fetching the rest of the pages in parallel
			pagesToFetch := int(math.Ceil(float64(meta.TotalHits)/float64(meta.HitsPerPage))) - 1
			channels := make([]chan Result[*AddressSearchResponse], pagesToFetch)
			for i := 0; i < len(channels); i++ {
				channels[i] = fetchPage(i + 1)
			}

			// Allocate space for the rest of the addresses (avoid re-alloc later)
			addresses = append(addresses, make([]Address, meta.TotalHits-meta.HitsPerPage)...)

			// Collect the results
			for i := 0; i < len(channels); i++ {
				result := <-channels[i]
				if result.Err != nil {
					ch <- Result[[]Address]{Ok: nil, Err: result.Err}
					return
				}
				for j, address := range (*result.Ok).Addresses {
					addresses[meta.HitsPerPage*(i+1)+j] = address
				}
			}
		}

		ch <- Result[[]Address]{Ok: &addresses, Err: nil}
	}()

	return ch
}

func ParseAddress(hitsPerPage int, data io.Reader) (*AddressSearchResponse, error) {
	var addressResponse AddressSearchResponse
	addressResponse.Addresses = make([]Address, hitsPerPage)

	err := json.NewDecoder(data).Decode(&addressResponse)
	if err != nil {
		return nil, err
	}

	return &addressResponse, nil
}

func (a Address) Format() string {
	return fmt.Sprintf("🏠 %s, %s %s", a.Text, a.PostCode, a.PostText)
}
