package goat

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const IP_SEARCH_BASE = "https://ipinfo.io/"

type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type IPLocation struct {
	Ip  string
	Loc Location
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
