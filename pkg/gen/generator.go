package gen

import (
	"fmt"
	"go/format"

	"golang.org/x/tools/go/packages"
)

type gen struct {
	i Inspector
	r Renderer
}

type Inspector interface {
	Inspect(pkg *packages.Package) (*File, error)
}

type Renderer interface {
	Render(f *File) ([]byte, error)
}

func NewGenerator(i Inspector, r Renderer) *gen {
	return &gen{i, r}
}

func loadPackage(targetPkg string) (*packages.Package, error) {
	p, err := packages.Load(&packages.Config{
		Mode:  packages.NeedSyntax | packages.NeedName | packages.NeedTypes | packages.NeedTypesInfo,
		Tests: false,
	}, targetPkg)
	if err != nil {
		return nil, fmt.Errorf("failed loading packages. err: %w", err)
	}
	if len(p) != 1 {
		return nil, fmt.Errorf("loaded unexpected amount of packages. want: 1, got: %d", len(p))
	}
	return p[0], nil
}

func (g *gen) Generate(targetPkg string) ([]byte, error) {
	pkg, err := loadPackage(targetPkg)
	if err != nil {
		return nil, err
	}
	out, err := g.i.Inspect(pkg)
	if err != nil {
		return nil, err
	}
	buf, err := g.r.Render(out)
	if err != nil {
		return nil, err
	}
	return g.formatOutput(buf)
}

func (gen) formatOutput(buf []byte) ([]byte, error) {
	src, err := format.Source(buf)
	if err != nil {
		return nil, fmt.Errorf("failed formatting the generated sources. err: %w", err)
	}
	return src, nil
}
