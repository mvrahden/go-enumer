package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/mvrahden/go-enumer/pkg/linter"
)

func main() {
	cfg := linter.Config{}

	singlechecker.Main(linter.New(&cfg))
}
