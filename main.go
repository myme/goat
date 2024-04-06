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

	query := strings.Join(os.Args[1:], " ")
	addresses, err := goat.SearchAddress(query)
	if err != nil {
		fmt.Println(err)
		return
	}

	ip, err := goat.SearchGeoIP()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Sort addresses by distance to IP location
	slices.SortFunc(addresses, func(a, b goat.Address) int {
		return cmp.Compare(Distance(ip.Loc, a.Loc), Distance(ip.Loc, b.Loc))
	})

	// Print addresses
	for _, address := range addresses {
		distance := Distance(ip.Loc, address.Loc)
		fmt.Printf("%.2f: %s, %s %s\n", distance, address.Text, address.PostCode, address.PostText)
	}
}
