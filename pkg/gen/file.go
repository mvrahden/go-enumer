package gen

import (
	"github.com/mvrahden/go-enumer/pkg/enumer"
)

type File struct {
	Header    Header
	Imports   []*Import
	TypeSpecs []*enumer.EnumType
}

type Header struct {
	Package Package
	Module  Module
}

type Package struct {
	Name string
	Path string
}

type Module struct {
	Module       string
	Path         string
	GoModPath    string
	GoModVersion string
	GoVersion    string
}

type Import struct {
	Name string // selector
	Path string
}
