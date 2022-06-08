package linter

import (
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/mvrahden/go-enumer/config"
	"github.com/mvrahden/go-enumer/pkg/utils/slices"
)

type EnumTypeConfig struct {
	Node *ast.Comment

	*config.Options
	FromSource string
}

func DefaultConfig(cfg *config.Options) *EnumTypeConfig {
	return &EnumTypeConfig{Options: cfg}
}

func (e *EnumType) DetectMagicComment() (c *ast.Comment, err error) {
	mcIdx := slices.FindIndex(e.Node.Doc.List, func(c *ast.Comment, idx int) bool {
		return MAGIC_MARKER.MatchString(c.Text)
	})
	if mcIdx == -1 {
		return nil, errors.New("no magic comment")
	}
	return e.Node.Doc.List[mcIdx], nil
}

// ExtractCommentString is a Noop func, but allows to be intercepted during tests
// to clean comment texts from test artifacts while keeping production code unaffected.
var ExtractCommentString = func(c *ast.Comment) string {
	return c.Text
}

func (e *EnumType) ParseMagicComment(mc *ast.Comment, opts *config.Options) (err error) {
	doc := ExtractCommentString(mc)

	cfg := DefaultConfig(opts)

	if args := strings.Split(doc, " "); len(args) > 1 {
		var f flag.FlagSet
		f.SetOutput(io.Discard) // silence flagset StdErr output

		f.StringVar(&cfg.TransformStrategy, "transform", cfg.TransformStrategy, "")
		f.Var(&cfg.Serializers, "serializers", "")
		f.Var(&cfg.SupportedFeatures, "support", "")
		f.StringVar(&cfg.FromSource, "from", "", "")
		err = f.Parse(args[1:]) // hint: parse w/o magic marker
		if err != nil {
			if els := strings.SplitAfter(err.Error(), "not defined: -"); len(els) == 2 { // flag provided but not defined: -<unknown opt>
				return fmt.Errorf("unknown option %q", els[1])
			}
			return fmt.Errorf("failed parsing doc-string. err: %s", err)
		}
		if f.NArg() > 0 {
			// report non-flag arguments
			return fmt.Errorf("unknown args %v", f.Args())
		}
	}

	if len(cfg.FromSource) > 0 {
		if !strings.HasSuffix(cfg.FromSource, ".csv") {
			return errors.New("unsupported file extension")
		}
		if strings.Contains(cfg.FromSource, "../") {
			return errors.New("source path cannot contain path traversals")
		}
		if strings.HasPrefix(cfg.FromSource, "./") || strings.HasPrefix(cfg.FromSource, "/") {
			return errors.New("source path cannot start with \"./\" or \"/\"")
		}
		cfg.FromSource = filepath.Clean(cfg.FromSource)
	}

	e.Config = cfg
	e.Config.Node = mc

	return nil
}

func (e *EnumType) GetPkgFS(fset *token.FileSet) (fs.FS, bool) {
	if len(e.Config.FromSource) == 0 {
		return nil, false
	}
	dirPath := filepath.Dir(fset.Position(e.Node.Pos()).Filename)
	pkgFs := os.DirFS(dirPath)
	return pkgFs, true
}

func (e *EnumType) ValidateEnumTypeConfig(fset *token.FileSet) error {
	// valdidate simple enum options
	// TODO

	// validate filebased enum options
	pkgFS, ok := e.GetPkgFS(fset)
	if !ok {
		return nil
	}

	_, err := fs.Stat(pkgFS, e.Config.FromSource)
	if errors.Is(err, os.ErrNotExist) {
		return errors.New("no such file")
	} else if err != nil {
		return fmt.Errorf("please verify the file. err: %w", err)
	}
	return nil
}

func detectBlockAndType(node ast.Node, typesInfo *types.Info, genFile *ast.File, enumTypes []*EnumType) (token.Pos, error) {
	decl, ok := node.(*ast.GenDecl)
	if !ok || decl.Tok != token.CONST {
		// we only care about const declarations
		return -1, nil
	}
	if len(decl.Specs) == 0 {
		// we only care about const declarations with types
		return -1, nil
	}
	{ // assert we are not inspecting anything from the generated files
		if genFile != nil && decl.Pos() >= genFile.Pos() && decl.Pos() <= genFile.End() {
			return -1, nil
		}
	}

	specs := slices.Filter(decl.Specs, func(v ast.Spec, idx int) bool {
		_, ok := v.(*ast.ValueSpec)
		return ok
	})

	valueSpecs := slices.Map(specs, func(v ast.Spec, idx int) *EnumValueSpec {
		return &EnumValueSpec{Node: v.(*ast.ValueSpec)}
	})
	// first enum type found will be assumed relevant type for an entire block
	var relevantEnumType *EnumType
	{ // find relevant enum type and assert it is the first declaration within the block
		typedConstIdx := slices.FindIndex(valueSpecs, func(vs *EnumValueSpec, _ int) bool {
			vstyp := vs.GetTypeVia(typesInfo)
			typeIdx := slices.FindIndex(enumTypes, func(et *EnumType, _ int) bool {
				return types.IdenticalIgnoreTags(vstyp, et.GetTypeVia(typesInfo))
			})
			if typeIdx == -1 {
				return false
			}
			relevantEnumType = enumTypes[typeIdx]
			return true
		})
		if typedConstIdx == -1 {
			// no enum was found, that's ok
			return -1, nil
		}
		// enums must always be defined as a pure and coherent blocks
		// we need to assert that here
		if typedConstIdx != 0 {
			// therefore we expect the first declaration to define the block's relevant enum type
			// hence this case needs to be reported and the block will not be validated any further
			return valueSpecs[typedConstIdx].Node.Pos(), errors.New("enum constants must be defined in a block of their own")
		}
		if relevantEnumType.ConstBlock != nil {
			// another block was assigned already
			// we want them to be in
			return decl.Pos(), errors.New("enum constants must be defined in a common block")
		}
	}

	// assign block to enum type definition
	relevantEnumType.ConstBlock = &EnumConstBlock{Specs: valueSpecs, Node: decl}
	return -1, nil
}
