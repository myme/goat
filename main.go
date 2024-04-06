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

func PrintAddresses(addresses []goat.Address, loc goat.Location) {
	for i, address := range addresses {
		distance := Distance(loc, address.Loc)
		fmt.Printf("%d: %.2f: %s, %s %s\n", i + 1, distance, address.Text, address.PostCode, address.PostText)
		fmt.Printf("  %f, %f\n", address.Loc.Lat, address.Loc.Lon)
	}
}

func PrintPlaces(places []goat.Place) {
	for _, place := range places {
		if place.CouldHaveGoats() {
			fmt.Printf("%s %s\n", place.Type, place.Name)
			fmt.Printf("%10.1fm pos: %f,%f\n", place.Distance, place.Loc.Lat, place.Loc.Lon)
		}
	}
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
	if len(*addresses.Ok) == 0 {
		fmt.Println("No addresses found")
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
		distA := Distance((*ip.Ok).Loc, a.Loc)
		distB := Distance((*ip.Ok).Loc, b.Loc)
		return cmp.Compare(distA, distB)
	})

	// Select an address to search for goats nearby
	PrintAddresses(*addresses.Ok, (*ip.Ok).Loc)

	// Find places near selected address
	places := <-goat.Places((*addresses.Ok)[0].Loc)
	if places.Err != nil {
		fmt.Println(places.Err)
		return
	}

	// Sort places by distance
	goat.SortPlaces(*places.Ok)
	PrintPlaces(*places.Ok)
}
