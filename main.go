package main

import (
	"log"

	"github.com/mvrahden/go-enumer/cmd/cli"
)

func main() {
	err := cli.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
