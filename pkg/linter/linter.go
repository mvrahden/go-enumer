package linter

import (
	"errors"
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"

	"github.com/mvrahden/go-enumer/config"
	"github.com/mvrahden/go-enumer/pkg/enumer"
	"github.com/mvrahden/go-enumer/pkg/utils/slices"
)

// Config the enum linter configuration.
type Config struct {
	// TODO:
	// - skip out-of-sync checks (with generated files)
	// - skip file contents inspection (e.g. CSV files)
}

// New creates an analyzer.
func New(c *Config) *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "enums",
		Doc:  "Checks the usage of Go enums.",
		Run: func(p *analysis.Pass) (interface{}, error) {
			if c == nil {
				return nil, nil
			}
			return run(p, c)
		},
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
}

func run(pass *analysis.Pass, c *Config) (_ interface{}, err error) {
	inspector, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, errors.New("missing inspect analyser")
	}

	validateGenerateCommand(inspector, pass)

	genFile := enumer.DetectGeneratedFile(pass.Files)

	enumTypes := determineEnumTypes(inspector, pass, genFile)
	if len(enumTypes) == 0 {
		// nothing to evaluate
		return nil, nil
	}
	enumTypes = validateEnumTypes(pass, enumTypes)

	enumTypes = detectAndValidateEnumConstBlocksForTypes(inspector, pass, genFile, enumTypes)

	enumTypes = loadAndValidateSpec(pass, enumTypes)

	// hint: marking existence of generated file is deferred to here
	// so that the enum blocks can be evaluated, even without the
	// existence of it.
	// However the subsequent checks are dependent on the generated file.
	if genFile == nil {
		pass.Reportf(enumTypes[0].Node.Pos(), "please generate enum file")
		return nil, nil
	}

	return nil, nil
}

func validateGenerateCommand(inspector *inspector.Inspector, pass *analysis.Pass) {
	// TODO
}

func determineEnumTypes(inspector *inspector.Inspector, pass *analysis.Pass, genFile *ast.File) []*enumer.EnumType {
	var enumTypes []*enumer.EnumType
	inspector.Preorder([]ast.Node{(*ast.GenDecl)(nil)}, func(n ast.Node) {
		et, pos, err := enumer.DetermineEnumType(n, pass.TypesInfo, genFile)
		if err != nil {
			pass.Reportf(pos, err.Error())
			return
		}
		if et == nil {
			return
		}
		enumTypes = append(enumTypes, et)
	})
	return enumTypes
}

func validateEnumTypes(pass *analysis.Pass, enumTypes []*enumer.EnumType) []*enumer.EnumType {
	enumTypes = slices.Filter(enumTypes, func(v *enumer.EnumType, idx int) bool {
		mc := v.DetectMagicComment()
		err := v.ParseMagicComment(mc, &config.Options{TransformStrategy: "noop"})
		if err != nil {
			pass.Reportf(mc.Pos(), err.Error())
			return false
		}
		err = v.ValidateEnumTypeConfig(pass.Fset)
		if err != nil {
			pass.Reportf(mc.Pos(), err.Error())
			return false
		}
		return true
	})

	// validate redundant source
	enumTypes = slices.Filter(enumTypes, func(v *enumer.EnumType, idx int) bool {
		if len(v.Config.FromSource) == 0 {
			return true
		}
		sameSource := slices.Any(enumTypes[:idx], func(v2 *enumer.EnumType, _ int) bool {
			return v.Config.FromSource == v2.Config.FromSource
		})
		if sameSource {
			pass.Reportf(v.Config.Node.Pos(), "enum of same file already exists")
			return false
		}
		return true
	})

	// validate file source existence
	enumTypes = slices.Filter(enumTypes, func(v *enumer.EnumType, idx int) bool {
		err := v.ValidateEnumTypeConfig(pass.Fset)
		if err != nil {
			pass.Reportf(v.Config.Node.Pos(), err.Error())
			return false
		}
		return true
	})

	return enumTypes
}

func detectAndValidateEnumConstBlocksForTypes(inspector *inspector.Inspector, pass *analysis.Pass, genFile *ast.File, enumTypes []*enumer.EnumType) []*enumer.EnumType {
	// Find relevant enum const blocks for enum types
	inspector.Preorder([]ast.Node{(*ast.GenDecl)(nil)}, func(n ast.Node) {
		pos, err := enumer.AssignEnumConstBlockToType(n, pass.TypesInfo, genFile, enumTypes)
		if err != nil {
			pass.Reportf(pos, err.Error())
			return
		}
	})

	return slices.Filter(enumTypes, func(v *enumer.EnumType, idx int) bool {
		err := v.ValidateConstBlock(pass.Fset, pass.TypesInfo)
		if err != nil {
			pass.Reportf(v.Node.Pos(), err.Error())
			return false
		}
		return true
	})
}

func loadAndValidateSpec(pass *analysis.Pass, enumTypes []*enumer.EnumType) []*enumer.EnumType {
	enumTypes = slices.Filter(enumTypes, func(v *enumer.EnumType, idx int) bool {
		err := v.LoadSpec(pass.Fset)
		if err != nil {
			pass.Reportf(v.Node.Pos(), err.Error())
			return false
		}
		return true
	})
	enumTypes = slices.Filter(enumTypes, func(v *enumer.EnumType, idx int) bool {
		err := v.ValidateSpec(pass.Fset, pass.TypesInfo)
		if err != nil {
			pass.Reportf(v.Node.Pos(), err.Error())
			return false
		}
		return true
	})
	return slices.Filter(enumTypes, func(v *enumer.EnumType, idx int) bool {
		err := v.CrossValidateConstBlockWithSpec(pass.Fset, pass.TypesInfo)
		if err != nil {
			pass.Reportf(v.Node.Pos(), err.Error())
			return false
		}
		return true
	})
}
