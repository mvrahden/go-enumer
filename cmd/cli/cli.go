package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/mvrahden/go-enumer/config"
	"github.com/mvrahden/go-enumer/pkg/gen"
)

const (
	ArgumentKeySupport           = "support"
	ArgumentKeySerializers       = "serializers"
	ArgumentKeyTransformStrategy = "transform"
	ArgumentKeyScanDirectory     = "dir"
	ArgumentKeyOutputFile        = "out"
)

func parseFlags(args []string, cArgs *config.Args, scanPath, outputFile *string) error {
	// setup flags
	flags := flag.NewFlagSet("", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	flags.StringVar(outputFile, ArgumentKeyOutputFile, "types_enumr", "the filename of the generated file; defaults to \"types_enumer\" which results in \"types_enumer.go\".")
	flags.StringVar(&cArgs.TransformStrategy, ArgumentKeyTransformStrategy, "noop", "string transformation (camel|pascal|kebab|snake|... see README.md); defaults to \"noop\" which applies no transormation to the enum values.")
	flags.Var(&cArgs.Serializers, ArgumentKeySerializers, "a list of opt-in serializers (binary|json|sql|text|yaml).")
	flags.Var(&cArgs.SupportedFeatures, ArgumentKeySupport, "a list of opt-in supported features (undefined|ignore-case|ent).")
	flags.StringVar(scanPath, ArgumentKeyScanDirectory, "", "directory of target package; defaults to CWD.")
	return flags.Parse(args)
}

type Generator interface {
	Generate(targetPkg string) ([]byte, error)
}

func Execute(args []string) error {
	var cArgs config.Args
	var scanPath, outputFile string
	err := parseFlags(args, &cArgs, &scanPath, &outputFile)
	if err != nil {
		return fmt.Errorf("failed parsing arguments. err: %s", err)
	}
	cfg := config.LoadWith(&cArgs)

	if err := validate(outputFile, cfg); err != nil {
		return fmt.Errorf("invalid arguments. err: %s", err)
	}

	targetDir, _ := os.Getwd()
	if len(scanPath) > 0 {
		if filepath.IsAbs(scanPath) {
			targetDir = filepath.Clean(scanPath)
		} else {
			targetDir = filepath.Join(targetDir, scanPath)
		}
	}

	g := gen.NewGenerator(
		gen.NewInspector(cfg),
		gen.NewRenderer(cfg),
	)
	buf, err := g.Generate(targetDir)
	if err != nil {
		return fmt.Errorf("failed generating code. err: %s", err)
	}

	f, err := os.Create(targetFilename(targetDir, outputFile, cfg))
	if err != nil {
		return fmt.Errorf("failed opening %q. err: %s", targetFilename(targetDir, outputFile, cfg), err)
	}
	defer f.Close()

	_, err = f.Write(buf)
	if err != nil {
		return fmt.Errorf("failed writing output to file. err: %s", err)
	}
	return nil
}

var targetFilename = func(dir, filename string, cfg *config.Options) string {
	filename = fmt.Sprintf("%s.go", filename)
	return filepath.Join(dir, filename)
}

func validate(filename string, cfg *config.Options) error {
	if len(filename) == 0 {
		return errors.New("output file name cannot be empty")
	}
	if strings.ContainsAny(filename, " ") {
		return errors.New("output file name contains spaces")
	}
	if strings.ContainsAny(filename, "\"") {
		return errors.New("output file name contains forbidden characters")
	}
	if cfg.Serializers.Contains("yaml") && cfg.Serializers.Contains("yaml.v3") {
		return fmt.Errorf("serializers %q and %q cannot be applied together", "yaml", "yaml.v3")
	}
	return nil
}
