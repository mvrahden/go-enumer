package cli

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/mvrahden/go-enumer/config"
	"github.com/mvrahden/go-enumer/pkg/gen"
)

var (
	args     config.Args
	scanPath string
)

func init() {
	// flag.StringVar(&args.Output, "output", "", "the filename of the generated file; defaults to \"<typealias|snake>_enumer.go\".")
	// flag.StringVar(&args.AddPrefix, "addprefix", "", "add given prefix to string values of enum.")
	flag.StringVar(&args.TransformStrategy, "transform", "noop", "string transformation (camel|kebab|snake|title|kebab-upper|snake-upper|whitespace); defaults to \"noop\" which applies no transormation to the enum value.")
	flag.StringVar(&args.TypeAliasName, "typealias", "", "the type alias (or type name) to perform the scan against.")
}

type Generator interface {
	Generate(targetPkg string) ([]byte, error)
}

func Execute() {
	flag.Parse()
	cfg := config.LoadWith(&args)

	targetDir, _ := os.Getwd()
	if len(scanPath) > 0 {
		targetDir = filepath.Clean(scanPath)
	}

	f, err := os.Create(filepath.Join(targetDir, fmt.Sprintf("%s_enumer.go", cfg.TypeAliasName)))
	if err != nil {
		log.Fatalf("failed opening %q. err: %s", filepath.Join(targetDir, fmt.Sprintf("%s_enumer.go", cfg.TypeAliasName)), err)
	}
	defer f.Close()

	g := gen.NewGenerator(
		gen.NewInspector(cfg),
		gen.NewRenderer(cfg),
	)
	buf, err := g.Generate(targetDir)
	if err != nil {
		log.Fatalf("failed generating code. err: %s", err)
	}

	_, err = f.Write(buf)
	if err != nil {
		log.Fatalf("failed writing output to file. err: %s", err)
	}
}
