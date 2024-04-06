package goat

import (
	"encoding/json"
	"io"
	"net/http"
)

const IP_SEARCH_BASE = "https://ipinfo.io/"

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

	return ParseGeoIP(res.Body)
}

func ParseGeoIP(data io.Reader) (*IPLocation, error) {
	var ipInfoResponse struct {
		Ip  string
		Loc string
	}

	err := json.NewDecoder(data).Decode(&ipInfoResponse)
	if err != nil {
		return nil, err
	}

	loc, err := ParseLocation(ipInfoResponse.Loc)
	if err != nil {
		return nil, err
	}

	return &IPLocation{Ip: ipInfoResponse.Ip, Loc: *loc}, nil
}
