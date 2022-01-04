package main

import (
	"log"
	"os"

	"github.com/mvrahden/go-enumer/cmd/cli"
)

func main() {
	err := cli.Execute(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
}
