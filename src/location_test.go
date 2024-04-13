package goat

import (
	"fmt"
	"testing"
)

func TestParseLocation(t *testing.T) {
	t.Run("ParseLocation() valid location", func (t *testing.T) {
		loc, _ := ParseLocation("59.9127,10.7461")
		expected := Location{Lat: 59.9127, Lon: 10.7461}
		if *loc != expected {
			t.Fatalf("expected %v, got %v", expected, loc)
		}
	})

	// Invalid tests cases
	invalidLocations := []string{
		"",
		"59.9127",
		"59.9127,10.7461,1337",
		"abcdef",
	}

	for _, loc := range invalidLocations {
		t.Run(fmt.Sprintf("ParseLocation() invalid location %s", loc), func (t *testing.T) {
			_, err := ParseLocation(loc)
			if err.Error() != fmt.Sprintf("invalid location format: %s", loc) {
				t.Fatalf("unexpected error, got: %s", err.Error())
			}
		})
	}
}
