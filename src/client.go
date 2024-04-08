package goat

import (
	"fmt"
	"io"
	"net/http"
)

type Result[t any] struct {
	Ok  *t
	Err error
}

// GetJSON fetches a JSON resource from a URL and parses it into a value of type t.
// TODO: Consider restricting t to only types that can be unmarshalled from JSON.
func GetJSON[t any](url string, parseJSON func(r io.Reader) (t, error)) chan Result[t] {
	// fmt.Printf("GetJSON: Fetching %s\n", url)
	ch := make(chan Result[t])

	go func() {
		client := http.Client{}
		res, err := client.Get(url)
		// fmt.Printf("GetJSON: Status %v\n", res.Status)
		if err != nil {
			ch <- Result[t]{Ok: nil, Err: err}
			return
		}
		if res.StatusCode != http.StatusOK {
			err := fmt.Errorf("HTTP error: %s", res.Status)
			ch <- Result[t]{Ok: nil, Err: err}
		}
		json, err := parseJSON(res.Body)
		ch <- Result[t]{Ok: &json, Err: err}
	}()

	return ch
}
