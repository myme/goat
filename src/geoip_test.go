package goat

import (
	"strings"
	"testing"
)

func TestParseGeoIP(t *testing.T) {
	data := `
		{
          "city": "Oslo",
          "country": "NO",
          "hostname": "40.51-175-185.customer.lyse.net",
          "ip": "51.175.185.40",
          "loc": "59.9127,10.7461",
          "org": "AS29695 Altibox AS",
          "postal": "0001",
          "readme": "https://ipinfo.io/missingauth",
          "region": "Oslo",
          "timezone": "Europe/Oslo"
		}
	`

	loc, _ := ParseGeoIP(strings.NewReader(data))

	expected := IPLocation{
		Ip:  "51.175.185.40",
		Loc: Location{Lat: 59.9127, Lon: 10.7461},
	}

	if *loc != expected {
		t.Fatalf("expected %v, got %v", expected, loc)
	}
}
