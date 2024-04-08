// Description: GeoIP search functionality.
//
// Search for the geolocation of an IP address using the https://ipinfo.io API.
// The request is not parameterized, so the IP address used is only the
// externally visible one.

package goat

import (
	"encoding/json"
	"io"
)

const IP_SEARCH_BASE = "https://ipinfo.io/"

type IPLocation struct {
	Ip  string
	Loc Location
}

func SearchGeoIP() chan Result[*IPLocation] {
	return GetJSON(IP_SEARCH_BASE, ParseGeoIP)
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
