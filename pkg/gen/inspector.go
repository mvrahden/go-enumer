package gen

import (
	"errors"
	"fmt"
	"go/ast"
	"go/types"
	"regexp"
	"strings"

	goinspect "golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/packages"

	"github.com/mvrahden/go-enumer/config"
	"github.com/mvrahden/go-enumer/pkg/enumer"
	"github.com/mvrahden/go-enumer/pkg/utils/slices"
)

var (
	matchGeneratedFileRegex = regexp.MustCompile(`^// Code generated .* DO NOT EDIT.$`)
	matchNumericValueRegex  = regexp.MustCompile(`^\-?\d+`)
	matchTypedHeaderRegex   = regexp.MustCompile(`^.+(\(.+\))$`)
)

type inspector struct {
	cfg *config.Options
}

func NewInspector(cfg *config.Options) *inspector {
	return &inspector{cfg}
}

func (i inspector) Inspect(pkg *packages.Package) (*File, error) {
	out := &File{Imports: []*Import{}}

	i.loadHeader(pkg, out)

	err := i.loadEnumTypes(pkg, out)
	if err != nil {
		return nil, err
	}

	i.determineImports(out)

	i.sortTypeSpecs(out)

	return out, nil
}

func (i inspector) loadEnumTypes(pkg *packages.Package, out *File) error {
	insp := goinspect.New(pkg.Syntax)
	genFile := enumer.DetectGeneratedFile(pkg.Syntax) // hint: get the generated enumer file

	enumTypes, err := i.detectTypeSpecs(insp, pkg.TypesInfo, genFile)
	if err != nil {
		return err
	}

	idx, err := slices.RangeErr(enumTypes, func(v *enumer.EnumType, _ int) error {
		mc := v.DetectMagicComment()
		if mc == nil {
			return errors.New("no magic comment") // hint: this should never happen
		}
		err = v.ParseMagicComment(mc, i.cfg)
		if err != nil {
			return err
		}
		return v.ValidateEnumTypeConfig(pkg.Fset)
	})
	if err != nil {
		goto SPEC_IS_INVALID
	}

	err = i.detectConstBlocks(insp, pkg.TypesInfo, genFile, enumTypes)
	if err != nil {
		return err
	}

	idx, err = slices.RangeErr(enumTypes, func(v *enumer.EnumType, _ int) error {
		return v.ValidateConstBlock(pkg.Fset, pkg.TypesInfo)
	})
	if err != nil {
		goto SPEC_IS_INVALID
	}

	idx, err = slices.RangeErr(enumTypes, func(v *enumer.EnumType, _ int) error {
		return v.LoadSpec(pkg.Fset)
	})
	if err != nil {
		goto SPEC_IS_INVALID
	}
	idx, err = slices.RangeErr(enumTypes, func(v *enumer.EnumType, _ int) error {
		return v.ValidateSpec(pkg.Fset, pkg.TypesInfo)
	})
	if err != nil {
		goto SPEC_IS_INVALID
	}
	idx, err = slices.RangeErr(enumTypes, func(v *enumer.EnumType, _ int) error {
		return v.CrossValidateConstBlockWithSpec(pkg.Fset, pkg.TypesInfo)
	})
SPEC_IS_INVALID:
	if err != nil {
		return fmt.Errorf("%q type specification is invalid. err: %w", enumTypes[idx].Name(), err)
	}

	out.TypeSpecs = enumTypes
	return nil
}

func (inspector) detectTypeSpecs(insp *goinspect.Inspector, typesInfo *types.Info, genFile *ast.File) ([]*enumer.EnumType, error) {
	var errs []error
	var enumTypes []*enumer.EnumType
	insp.Preorder([]ast.Node{(*ast.GenDecl)(nil)}, func(n ast.Node) {
		et, _, err := enumer.DetermineEnumType(n, typesInfo, genFile)
		if err != nil {
			typ := n.(*ast.GenDecl).Specs[0].(*ast.TypeSpec)
			errs = append(errs, fmt.Errorf("%q type specification is invalid. err: %s", typ.Name, err))
			return
		}
		if et != nil {
			enumTypes = append(enumTypes, et)
		}
	})
	if len(errs) > 0 {
		return nil, errs[0]
	}
	return enumTypes, nil
}

func (inspector) detectConstBlocks(insp *goinspect.Inspector, typesInfo *types.Info, genFile *ast.File, enumTypes []*enumer.EnumType) error {
	var errs []error
	insp.Preorder([]ast.Node{(*ast.GenDecl)(nil)}, func(n ast.Node) {
		_, err := enumer.AssignEnumConstBlockToType(n, typesInfo, genFile, enumTypes)
		if err != nil {
			errs = append(errs, err)
		}
	})
	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}

func (inspector) loadHeader(pkg *packages.Package, out *File) {
	out.Header.Package = Package{
		Name: pkg.Name,
		Path: pkg.PkgPath,
	}
	if pkg.Module == nil {
		return
	}
	out.Header.Module = Module{
		Module:       pkg.Module.Version,
		Path:         pkg.Module.Dir,
		GoModPath:    pkg.Module.GoMod,
		GoModVersion: pkg.Module.Version,
		GoVersion:    pkg.Module.GoVersion,
	}
}

func (i inspector) determineImports(f *File) {
	f.Imports = append(f.Imports, &Import{Path: "errors"})
	f.Imports = append(f.Imports, &Import{Path: "fmt"})

	// we add all imports (also duplicates)
	for _, ts := range f.TypeSpecs {
		for _, v := range ts.Config.Serializers {
			switch v {
			case "gql":
				f.Imports = append(f.Imports, &Import{Path: "io"})
				f.Imports = append(f.Imports, &Import{Path: "strconv"})
			case "json":
				f.Imports = append(f.Imports, &Import{Path: "encoding/json"})
			case "sql":
				f.Imports = append(f.Imports, &Import{Path: "database/sql/driver"})
			case "yaml.v3":
				f.Imports = append(f.Imports, &Import{Path: "gopkg.in/yaml.v3"})
			}
		}
	}
	// sort alphabetically
	f.Imports = slices.SortStable(f.Imports, func(s []*Import, i, j int) bool {
		return strings.Compare(s[i].Path, s[j].Path) < 0 && strings.Compare(s[i].Name, s[j].Name) < 0
	})
	// filter unique import values
	f.Imports = slices.Filter(f.Imports, func(v *Import, idx int) bool {
		if idx == 0 {
			return true
		}
		return v.Path != f.Imports[idx-1].Path || v.Name != f.Imports[idx-1].Name
	})
}

func (i inspector) sortTypeSpecs(f *File) {
	// sort all enums
	f.TypeSpecs = slices.SortStable(f.TypeSpecs, func(s []*enumer.EnumType, i, j int) bool {
		return strings.Compare(s[i].Name().Name, s[j].Name().Name) < 0
	})
}
