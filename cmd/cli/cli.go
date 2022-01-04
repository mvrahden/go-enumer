package cli

import (
	"flag"
	"fmt"
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

const (
	ArgumentKeyTypeAlias         = "typealias"
	ArgumentKeySupport           = "support"
	ArgumentKeySerializers       = "serializers"
	ArgumentKeyTransformStrategy = "transform"
	ArgumentKeyScanDirectory     = "dir"
)

func init() {
	// flag.StringVar(&args.Output, "output", "", "the filename of the generated file; defaults to \"<typealias|snake>_enumer.go\".")
	// flag.StringVar(&args.AddPrefix, "addprefix", "", "add given prefix to string values of enum.")
	flag.StringVar(&args.TransformStrategy, ArgumentKeyTransformStrategy, "noop", "string transformation (camel|pascal|kebab|snake|... see README.md); defaults to \"noop\" which applies no transormation to the enum values.")
	flag.StringVar(&args.TypeAliasName, ArgumentKeyTypeAlias, "", "the type alias (or type name) to perform the scan against.")
	flag.Var(&args.Serializers, ArgumentKeySerializers, "a list of opt-in serializers (binary|json|sql|text|yaml).")
	flag.Var(&args.SupportedFeatures, ArgumentKeySupport, "a list of opt-in supported features (undefined|ent).")
	flag.StringVar(&scanPath, ArgumentKeyScanDirectory, "", "directory of target package; defaults to CWD.")
}

type Generator interface {
	Generate(targetPkg string) ([]byte, error)
}

func Execute() error {
	flag.Parse()
	cfg := config.LoadWith(&args)

	if err := validate(cfg); err != nil {
		return fmt.Errorf("invalid arguments. err: %s", err)
	}

	targetDir, _ := os.Getwd()
	if len(scanPath) > 0 {
		targetDir = filepath.Clean(scanPath)
	}

	g := gen.NewGenerator(
		gen.NewInspector(cfg),
		gen.NewRenderer(cfg),
	)
	buf, err := g.Generate(targetDir)
	if err != nil {
		return fmt.Errorf("failed generating code. err: %s", err)
	}

	f, err := os.Create(targetFilename(targetDir, cfg))
	if err != nil {
		return fmt.Errorf("failed opening %q. err: %s", targetFilename(targetDir, cfg), err)
	}
	defer f.Close()

	_, err = f.Write(buf)
	if err != nil {
		return fmt.Errorf("failed writing output to file. err: %s", err)
	}
	return nil
}

func targetFilename(dir string, cfg *config.Options) string {
	filename := fmt.Sprintf("%s_enumer.go", strings.ToLower(cfg.TypeAliasName))
	return filepath.Join(dir, filename)
}

func validate(cfg *config.Options) error {
	if len(cfg.TypeAliasName) == 0 {
		return fmt.Errorf("argument %q cannot be empty.", ArgumentKeyTypeAlias)
	}
	if cfg.Serializers.Contains("yaml") && cfg.Serializers.Contains("yaml.v3") {
		return fmt.Errorf("serializers %q and %q are cannot be applied together.", "yaml", "yaml.v3")
	}
	return nil
}
