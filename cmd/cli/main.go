package cli

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/mvrahden/go-enumer/config"
)

var (
	args     config.Args
	scanPath string
)

func init() {
	flag.StringVar(&args.Output, "output", "", "the filename of the generated file; defaults to \"<typealias|snake>_enumer.go\".")
	flag.StringVar(&args.AddPrefix, "addprefix", "", "add given prefix to string values of enum.")
	flag.StringVar(&args.TransformStrategy, "transform", "noop", "string transformation (camel|snake|kebab|title); defaults to \"noop\" which applies no transormation to the enum value.")
	flag.StringVar(&args.TypeAliasName, "typealias", "", "the type alias (or type name) to perform the scan against.")
}

type Generator interface {
	Generate(targetPkg string) ([]byte, error)
}

func Execute(g Generator) {
	flag.Parse()
	cfg := config.LoadWith(&args)

	targetDir, _ := os.Getwd()
	if len(scanPath) > 0 {
		targetDir = filepath.Clean(scanPath)
	}

	f, err := os.Create(filepath.Join(targetDir, "types_jsoner.go"))
	if err != nil {
		log.Fatalf("failed opening %q. err: %s", filepath.Join(targetDir, "types_jsoner.go"), err)
	}
	defer f.Close()

	buf, err := g.Generate(targetDir)
	if err != nil {
		log.Fatalf("failed generating code. err: %s", err)
	}

	_, err = f.Write(buf)
	if err != nil {
		log.Fatalf("failed writing output to file. err: %s", err)
	}
}
