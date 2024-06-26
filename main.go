package main

import (
	"cmp"
	"flag"
	"fmt"
	"math"
	"os"
	"slices"
	"strings"

	goat "github.com/myme/goat/src"
)

// Calculate the pythagorean distance between two locations.
func Distance(l1 goat.Location, l2 goat.Location) float64 {
	x := l2.Lat - l1.Lat
	y := l2.Lon - l1.Lon
	return math.Sqrt(math.Pow(x, 2) + math.Pow(y, 2))
}

// Print a list of places with their type, name, distance and location.
// Only include places where the presence of goats are likely.
func PrintPlaces(places []goat.Place, maxResults int) {
	potentialGoatPlaces := 0
	for _, place := range places {
		if place.CouldHaveGoats() {
			potentialGoatPlaces++
			if potentialGoatPlaces > maxResults {
				break
			}
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

// Select an address from a list of addresses.
// If there is only one address, it is returned immediately.
func SelectAddress(addresses []goat.Address, autoSelect bool) (*goat.Address, error) {
	// Return only match
	if len(addresses) == 1 && autoSelect {
		return &addresses[0], nil
	}

	listItems := make([]goat.Item, len(addresses))
	for i, address := range addresses {
		listItems[i] = goat.Item{
			Index:  i,
			Text:   address.Text,
			Desc:   fmt.Sprintf("%s %s", address.PostCode, address.PostText),
			Filter: address.Format(),
		}
	}
	item, err := goat.SelectFromList("Select an address", listItems)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &addresses[item.Index], nil
}

type Options struct {
	maxResults   int
	noAutoSelect bool
}

func main() {

	// CLI args parser
	var options Options
	flag.IntVar(&options.maxResults, "max-results", 10, "Maximum number of search results")
	flag.BoolVar(&options.noAutoSelect, "no-auto-select", false, "Don't accept first search result automatically")
	flag.Usage = func() {
		out := flag.CommandLine.Output()
		fmt.Fprintf(out, "Usage of %s:\n\n", os.Args[0])
		fmt.Fprintf(out, "Explore places where the probability of goats is greater than 0.\n\n")
		fmt.Fprintf(out, "goat <address> <...address>\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Usage: goat <address> <...address>")
		return
	}

	// Run parallel queries
	query := strings.Join(args, " ")
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
	address, err := SelectAddress(*addresses.Ok, !options.noAutoSelect)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Selected address: %v\n\n", address.Format())

	// Find, sort & print places near selected address.
	// These are supposedly "farm" areas, where goats could be found.
	places := <-goat.Places(address.Loc)
	if places.Err != nil {
		fmt.Println(places.Err)
		return
	}
	goat.SortPlaces(*places.Ok)
	PrintPlaces(*places.Ok, options.maxResults)

}
