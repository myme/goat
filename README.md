# A basic Go project

## Setup

```sh
$ nix flake new -t github:myme/nix-templates#go goat
$ cd goat
$ go mod init myme.no/goat
```

## Create a main module

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

## Run it

```sh
go run .
```
