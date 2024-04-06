package main

import (
	"cmp"
	"fmt"
	"math"
	"os"
	"slices"
	"strings"

	goat "myme.no/goat/src"
)

func Distance(l1 goat.Location, l2 goat.Location) float64 {
	x := l2.Lat - l1.Lat
	y := l2.Lon - l1.Lon
	return math.Sqrt(math.Pow(x, 2) + math.Pow(y, 2))
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: goat <address> <...address>")
		return
	}

	// Run parallel queries
	query := strings.Join(os.Args[1:], " ")
	addrCh := goat.SearchAddress(query)
	geoIPCh := goat.SearchGeoIP()

	// Wait for address results
	addresses := <-addrCh
	if addresses.Err != nil {
		fmt.Println("Error fetching addresses:", addresses.Err)
		return
	}

	// Wait for GeoIP results
	ip := <-geoIPCh
	if ip.Err != nil {
		fmt.Println("Error fetching IP location:", ip.Err)
		return
	}

	// Sort addresses by distance to IP location
	slices.SortFunc(*addresses.Ok, func(a, b goat.Address) int {
		return cmp.Compare(Distance((*ip.Ok).Loc, a.Loc), Distance((*ip.Ok).Loc, b.Loc))
	})

	// Print addresses
	for _, address := range *addresses.Ok {
		distance := Distance((*ip.Ok).Loc, address.Loc)
		fmt.Printf("%.2f: %s, %s %s\n", distance, address.Text, address.PostCode, address.PostText)
	}
}
