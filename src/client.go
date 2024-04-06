package goat

import (
	// "fmt"
	"io"
	"net/http"
)

type Result[t any] struct {
	Ok  *t
	Err error
}

func GetJSON[t any](url string, parseJSON func(r io.Reader) (t, error)) chan Result[t] {
	// fmt.Printf("GetJSON: Fetching %s\n", url)
	ch := make(chan Result[t])

	go func() {
		client := http.Client{}
		res, err := client.Get(url)
		// fmt.Printf("GetJSON: Status %v\n", res.Status)
		if err != nil {
			ch <- Result[t]{Ok: nil, Err: err}
		}
		json, err := parseJSON(res.Body)
		ch <- Result[t]{Ok: &json, Err: err}
	}()

	return ch
}
