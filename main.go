package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const IP_SEARCH_BASE = "https://ipinfo.io/"
const ADDRESS_SEARCH_BASE = "https://ws.geonorge.no/adresser/v1/sok"

type IPLocation struct {
	Ip  string
	Loc Location
}

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

func SearchGeoIP() (*IPLocation, error) {
	client := http.Client{}

	res, err := client.Get(IP_SEARCH_BASE)
	if err != nil {
		return nil, err
	}

	var ipInfoResponse struct {
		Ip  string
		Loc string
	}
	err = json.NewDecoder(res.Body).Decode(&ipInfoResponse)
	if err != nil {
		return nil, err
	}

	locParts := strings.FieldsFunc(ipInfoResponse.Loc, func(r rune) bool {
		return r == ','
	})
	if len(locParts) != 2 {
		return nil, fmt.Errorf("invalid location format: %s", ipInfoResponse.Loc)
	}

	lat, err := strconv.ParseFloat(locParts[0], 64)
	if err != nil {
		return nil, err
	}

	lon, err := strconv.ParseFloat(locParts[1], 64)
	if err != nil {
		return nil, err
	}

	loc := Location{Lat: lat, Lon: lon}
	return &IPLocation{Ip: ipInfoResponse.Ip, Loc: loc}, nil
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
	addresses, err := SearchAddress(query)
	if err != nil {
		fmt.Println(err)
		return
	}

	ip, err := SearchGeoIP()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(ip)
	fmt.Println(addresses)
}
