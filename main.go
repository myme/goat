package main

import (
	"bufio"
	"cmp"
	"fmt"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"

	goat "myme.no/goat/src"
)

func Distance(l1 goat.Location, l2 goat.Location) float64 {
	x := l2.Lat - l1.Lat
	y := l2.Lon - l1.Lon
	return math.Sqrt(math.Pow(x, 2) + math.Pow(y, 2))
}

// Print a list of places with their type, name, distance and location.
// Only include places where the presence of goats are likely.
func PrintPlaces(places []goat.Place) {
	for _, place := range places {
		if place.CouldHaveGoats() {
			url := fmt.Sprintf(
				"https://www.google.com/maps/place/?q=%f,%f&t=k",
				place.Loc.Lat,
				place.Loc.Lon,
			)
			fmt.Printf("🐐 [%s] %s\n", place.Type, place.Name)
			fmt.Printf("%10.1fm map: %s\n", place.Distance, url)
		}
	}
}

func SelectMatchingAddress(addresses []goat.Address) (*goat.Address, error) {
	if len(addresses) == 1 {
		return &addresses[0], nil
	}

	fmt.Print("Select an address: ")
	input, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("Error reading input: %v", err)
	}
	index, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil || index < 1 || index > len(addresses) {
		return nil, fmt.Errorf("Index must be a number between 1 and %d", len(addresses))
	}

	return &addresses[index-1], nil
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

	// Present selection list of addresses
	listItems := make([]goat.Item, len(*addresses.Ok))
	for i, address := range *addresses.Ok {
		listItems[i] = goat.Item{
			Index: i,
			Text: address.Text,
			Desc: fmt.Sprintf("%s %s", address.PostCode, address.PostText),
			Filter: address.Format(),
		}
	}
	item, err := goat.SelectFromList("Select an address", listItems)
	if err != nil {
		fmt.Println(err)
		return
	}

	address := (*addresses.Ok)[item.Index]
	fmt.Printf("Selected address: %v\n\n", address.Format())

	// Find, sort & print places near selected address.
	// These are supposedly "farm" areas, where goats could be found.
	places := <-goat.Places(address.Loc)
	if places.Err != nil {
		fmt.Println(places.Err)
		return
	}
	goat.SortPlaces(*places.Ok)
	PrintPlaces(*places.Ok)
}
