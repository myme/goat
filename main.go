package main

import (
	"cmp"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"os"
	"slices"
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

	loc, err := ParseLocation(ipInfoResponse.Loc)
	if err != nil {
		return nil, err
	}

	return &IPLocation{Ip: ipInfoResponse.Ip, Loc: *loc}, nil
}

func ParseLocation(loc string) (*Location, error) {
	locParts := strings.FieldsFunc(loc, func(r rune) bool {
		return r == ','
	})
	if len(locParts) != 2 {
		return nil, fmt.Errorf("invalid location format: %s", loc)
	}

	lat, err := strconv.ParseFloat(locParts[0], 64)
	if err != nil {
		return nil, err
	}

	lon, err := strconv.ParseFloat(locParts[1], 64)
	if err != nil {
		return nil, err
	}

	return &Location{Lat: lat, Lon: lon}, nil
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

func Distance(l1 Location, l2 Location) float64 {
	x := l2.Lat - l1.Lat
	y := l2.Lon - l1.Lon
	return math.Sqrt(math.Pow(x, 2) + math.Pow(y, 2))
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

	// Sort addresses by distance to IP location
	slices.SortFunc(addresses, func(a, b Address) int {
		return cmp.Compare(Distance(ip.Loc, a.Loc), Distance(ip.Loc, b.Loc))
	})

	// Print addresses
	for _, address := range addresses {
		distance := Distance(ip.Loc, address.Loc)
		fmt.Printf("%.2f: %s, %s %s\n", distance, address.Text, address.PostCode, address.PostText)
	}
}
