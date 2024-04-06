package goat

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

func SearchAddress(query string) ([]Address, error) {
	client := http.Client{}

	url := fmt.Sprintf("%s?sok=%s", ADDRESS_SEARCH_BASE, url.QueryEscape(query))
	res, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	return ParseAddress(res.Body)
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
