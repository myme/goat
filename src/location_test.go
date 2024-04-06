package goat

import "testing"

func TestParseLocation(t *testing.T) {
	// Valid location, sunshine case
	{
		loc, _ := ParseLocation("59.9127,10.7461")
		expected := Location{Lat: 59.9127, Lon: 10.7461}
		if *loc != expected {
			t.Fatalf("expected %v, got %v", expected, loc)
		}
	}

	// Invalid location, empty
	{
		_, err := ParseLocation("")
		if err.Error() != "invalid location format: " {
			t.Fatalf("unexpected error, got: %s", err.Error())
		}
	}

	// Invalid location, single field
	{
		_, err := ParseLocation("59.9127")
		if err.Error() != "invalid location format: 59.9127" {
			t.Fatalf("unexpected error, got: %s", err.Error())
		}
	}

	// Invalid location, too many fields
	{
		_, err := ParseLocation("59.9127,10.7461,1337")
		if err.Error() != "invalid location format: 59.9127,10.7461,1337" {
			t.Fatalf("unexpected error, got: %s", err.Error())
		}
	}

	// Invalid location, alphbetical
	{
		_, err := ParseLocation("abcdef")
		if err.Error() != "invalid location format: abcdef" {
			t.Fatalf("unexpected error, got: %s", err.Error())
		}
	}
}
