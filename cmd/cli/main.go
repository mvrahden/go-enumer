package cli

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

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
	flag.StringVar(&args.TransformStrategy, "transform", "noop", "string transformation (camel|kebab|pascal|snake|upper-kebab|upper-snake|whitespace); defaults to \"noop\" which applies no transormation to the enum value.")
	flag.StringVar(&args.TypeAliasName, "typealias", "", "the type alias (or type name) to perform the scan against.")
	flag.Var(&args.Serializers, "serializers", "a list of opt-in serializers (binary|json|sql|text|yaml).")
	flag.Var(&args.SupportedFeatures, "support", "a list of opt-in supported features (undefined|ent).")
	flag.StringVar(&scanPath, "dir", "", "directory of target package; defaults to CWD.")
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

	f, err := os.Create(targetFilename(targetDir, cfg))
	if err != nil {
		log.Fatalf("failed opening %q. err: %s", targetFilename(targetDir, cfg), err)
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

func targetFilename(dir string, cfg *config.Options) string {
	filename := fmt.Sprintf("%s_enumer.go", strings.ToLower(cfg.TypeAliasName))
	return filepath.Join(dir, filename)
}
