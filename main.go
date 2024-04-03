package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
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

type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

func SearchAddress(query string) ([]Address, error) {
	client := http.Client{}

	url := fmt.Sprintf("%s?sok=%s", ADDRESS_SEARCH_BASE, url.QueryEscape(query))
	res, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	var addressResponse AdressSearchResponse
	err = json.NewDecoder(res.Body).Decode(&addressResponse)
	if err != nil {
		return nil, err
	}

	return addressResponse.Addresses, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: goat <address> <...address>")
		return
	}

	query := strings.Join(os.Args[1:], " ")
	res, err := SearchAddress(query)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(res)
}
