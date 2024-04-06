package goat

import (
	"fmt"
	"strconv"
	"strings"
)

type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
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
