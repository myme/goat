package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
)

func GetUrl(url string) (string, error) {
	client := http.Client{}

	res, err := client.Get(url)
	if err != nil {
		return "", err
	}

	lines := ""
	reader := bufio.NewReader(res.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return "", err
		}
		lines += line
	}

	return lines, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: goat <url>")
		return
	}

	res, err := GetUrl(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(res)
}
