package linter

import (
	"errors"
	"go/ast"
	"go/constant"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"

	"github.com/mvrahden/go-enumer/config"
	"github.com/mvrahden/go-enumer/pkg/common"
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

	genFile := common.DetermineGeneratedFile(pass.Files)

	enumTypes := determineEnumTypes(inspector, pass, genFile)
	if len(enumTypes) == 0 {
		// nothing to evaluate
		return nil, nil
	}
	enumTypes = validateEnumTypes(pass, enumTypes)

	determineEnumBlocksForTypes(inspector, pass, genFile, enumTypes)
	enumTypes = validateEnumBlocks(pass, enumTypes)

	// hint: marking existence of generated file is deferred to here
	// so that the enum blocks can be evaluated, even without the
	// existence of it.
	// However the subsequent checks are dependent on the generated file.
	if genFile == nil {
		pass.Reportf(enumTypes[0].Node.Pos(), "please generate enum file")
		return nil, nil
	}

	validateInOfSyncWithGeneratedFile(pass, enumTypes, genFile)

	return nil, nil
}

func validateGenerateCommand(inspector *inspector.Inspector, pass *analysis.Pass) {
	// TODO
}

func determineEnumTypes(inspector *inspector.Inspector, pass *analysis.Pass, genFile *ast.File) []*common.EnumType {
	var enumTypes []*common.EnumType
	inspector.Preorder([]ast.Node{(*ast.GenDecl)(nil)}, func(n ast.Node) {
		et, pos, err := common.DetermineEnumType(n, pass.TypesInfo, genFile)
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

func validateEnumTypes(pass *analysis.Pass, enumTypes []*common.EnumType) []*common.EnumType {
	enumTypes = slices.Filter(enumTypes, func(v *common.EnumType, idx int) bool {
		mc, err := v.DetectMagicComment()
		if err != nil {
			pass.Reportf(v.Node.Pos(), err.Error())
			return false
		}
		err = v.ParseMagicComment(mc, &config.Options{TransformStrategy: "noop"})
		if err != nil {
			pass.Reportf(mc.Pos(), err.Error())
			return false
		}
		return true
	})

	// validate redundant source
	enumTypes = slices.Filter(enumTypes, func(v *common.EnumType, idx int) bool {
		if len(v.Config.FromSource) == 0 {
			return true
		}
		sameSource := slices.Any(enumTypes[:idx], func(v2 *common.EnumType, _ int) bool {
			return v.Config.FromSource == v2.Config.FromSource
		})
		if sameSource {
			pass.Reportf(v.Config.Node.Pos(), "enum of same file already exists")
			return false
		}
		return true
	})

	// validate file source existence
	enumTypes = slices.Filter(enumTypes, func(v *common.EnumType, idx int) bool {
		err := v.ValidateEnumTypeConfig(pass.Fset)
		if err != nil {
			pass.Reportf(v.Config.Node.Pos(), err.Error())
			return false
		}
		return true
	})

	return enumTypes
}

func determineEnumBlocksForTypes(inspector *inspector.Inspector, pass *analysis.Pass, genFile *ast.File, enumTypes []*common.EnumType) {
	// Find relevant enum const blocks for enum types
	inspector.Preorder([]ast.Node{(*ast.GenDecl)(nil)}, func(n ast.Node) {
		pos, err := common.AssignEnumConstBlockToType(n, pass.TypesInfo, genFile, enumTypes)
		if err != nil {
			pass.Reportf(pos, err.Error())
			return
		}
	})
}

func validateEnumBlocks(pass *analysis.Pass, enumTypes []*common.EnumType) []*common.EnumType {
	bareEnums, blockEnums := slices.SplitBy(enumTypes, func(v *common.EnumType, idx int) bool {
		return v.ConstBlock == nil
	})
	// inspect those enums without blocks
	// enums that do not originate from a file require blocks or are obsolete otherwise
	bareEnums = slices.Filter(bareEnums, func(v *common.EnumType, idx int) bool {
		if len(v.Config.FromSource) == 0 {
			pass.Reportf(v.Node.Pos(), "enum types require a const block or a file source")
			return false
		}
		return true
	})

	// inspect const blocks for enum types
	// assert relative positioning of const blocks to their relevant types
	blockEnums = slices.Filter(blockEnums, func(v *common.EnumType, idx int) bool {
		// assert const block is in same file
		if pass.Fset.File(v.ConstBlock.Node.Pos()) != pass.Fset.File(v.Node.Pos()) {
			pass.Reportf(v.ConstBlock.Node.Pos(), "enum blocks must be in same file as their type definition")
			return false
		}
		// assert const block is after relevant type
		if v.ConstBlock.Node.Pos() < v.Node.Pos() {
			pass.Reportf(v.ConstBlock.Node.Pos(), "enum blocks must be defined after their type definition")
			return false
		}
		return true
	})

	// report all rowed up enum value assignments
	blockEnums = slices.Filter(blockEnums, func(v *common.EnumType, idx int) bool {
		ok := slices.None(v.ConstBlock.Specs, func(v *common.EnumValueSpec, idx int) bool {
			return len(v.Node.Names) > 1 || len(v.Node.Values) > 1
		})
		if !ok {
			pass.Reportf(v.ConstBlock.Node.Pos(), "enum blocks must not contain rowed declarations")
			return false
		}
		return true
	})

	// assert only enum values of relevant enum type within block
	blockEnums = slices.Filter(blockEnums, func(v *common.EnumType, idx int) bool {
		ok := slices.None(v.ConstBlock.Specs, func(curr *common.EnumValueSpec, idx int) bool {
			if idx == 0 { // assert that first type declaration is an enum
				isSameAsBlockType := types.IdenticalIgnoreTags(
					curr.GetTypeVia(pass.TypesInfo),
					v.GetTypeVia(pass.TypesInfo),
				)
				return !isSameAsBlockType
			}
			relevantEnumType := v.ConstBlock.Specs[0].GetTypeVia(pass.TypesInfo)
			isSameAsBlockType := types.IdenticalIgnoreTags(
				pass.TypesInfo.TypeOf(curr.Node.Names[0]),
				relevantEnumType,
			)
			return !isSameAsBlockType
		})
		if !ok {
			pass.Reportf(v.ConstBlock.Node.Pos(), "enum blocks must not contain unrelated type declarations")
			return false
		}
		return true
	})

	// assert order of values in block
	blockEnums = slices.Filter(blockEnums, func(v *common.EnumType, _ int) bool {
		badIdx, err := slices.RangeErr(v.ConstBlock.Specs, func(v *common.EnumValueSpec, idx int) error {
			val := v.GetObjectVia(pass.TypesInfo).(*types.Const).Val()
			{
				val, ok := constant.Int64Val(val)
				if !ok {
					return errors.New("invalid numerical format")
				}
				// hint: reflow otherwise overflown values.
				// that's ok as we're dealing with uint types exclusively
				v.Value = uint64(val)
			}
			return nil
		})
		if badIdx > -1 {
			pass.Reportf(v.ConstBlock.Node.Pos(), err.Error())
			return false
		}

		if isNotFileBased := len(v.Config.FromSource) == 0; isNotFileBased && v.ConstBlock.Specs[0].Value > 1 {
			// hint: file based enums can start with arbitrary numbers in const blocks
			// as they do not represent the SPEC in this case but merely refer to individual values.
			pass.Reportf(v.ConstBlock.Node.Pos(), "enum block sequences must start with either 0 or 1")
			return false
		}

		badIdx, err = slices.RangeErr(v.ConstBlock.Specs, func(vs *common.EnumValueSpec, idx int) error {
			if idx == 0 {
				return nil
			}
			prev := v.ConstBlock.Specs[idx-1].Value
			if prev > vs.Value {
				return errors.New("enum block sequences must be ordered")
			}
			return nil
		})
		if badIdx > -1 {
			pass.Reportf(v.ConstBlock.Node.Pos(), err.Error())
			return false
		}

		badIdx, err = slices.RangeErr(v.ConstBlock.Specs, func(vs *common.EnumValueSpec, idx int) error {
			if idx == 0 {
				return nil
			}
			prev := v.ConstBlock.Specs[idx-1].Value
			if prev+1 < vs.Value {
				return errors.New("enum block sequences must increment at most by one")
			}
			return nil
		})
		if badIdx > -1 {
			pass.Reportf(v.ConstBlock.Node.Pos(), err.Error())
			return false
		}
		return true
	})
	return append(bareEnums, blockEnums...)
}

func validateInOfSyncWithGeneratedFile(pass *analysis.Pass, enumTypes []*common.EnumType, genFile *ast.File) {
	filebasedEnums, simpleEnums := slices.SplitBy(enumTypes, func(v *common.EnumType, idx int) bool {
		return len(v.Config.FromSource) > 0
	})

	validateBlockConstantsOfSimpleEnums(pass, simpleEnums, genFile)

	filebasedEnumsWithBlocks := slices.Filter(filebasedEnums, func(v *common.EnumType, idx int) bool {
		return v.ConstBlock != nil
	})
	validateBlockConstantsOfFilebasedEnums(pass, filebasedEnumsWithBlocks, genFile)
}

func validateBlockConstantsOfSimpleEnums(pass *analysis.Pass, enumTypes []*common.EnumType, genFile *ast.File) {
	// TODO
}

func validateBlockConstantsOfFilebasedEnums(pass *analysis.Pass, enumTypes []*common.EnumType, genFile *ast.File) {
	slices.Range(enumTypes, func(v *common.EnumType, idx int) {
		pkgFS, ok := v.GetPkgFS(pass.Fset)
		if !ok {
			return
		}
		f, err := pkgFS.Open(v.Config.FromSource)
		if err != nil {
			pass.Reportf(v.Node.Pos(), "could not open file. err: %s", err)
			return
		}
		defer f.Close()
	})

	// detect duplicate constants
	// this is a special case for enum constants that are referring to generated values from a source file like CSV
	// it is not a good practice to have redundant definitions in these cases
	slices.Range(enumTypes, func(v *common.EnumType, idx int) {
		slices.Range(v.ConstBlock.Specs, func(s1 *common.EnumValueSpec, idx int) {
			if idx == len(v.ConstBlock.Specs) {
				return
			}
			slices.Range(v.ConstBlock.Specs[idx+1:], func(s2 *common.EnumValueSpec, idx int) {
				if s1.Value == s2.Value {
					pass.Reportf(s2.Node.Pos(), "redundant constant")
					return
				}
			})
		})
	})
}
