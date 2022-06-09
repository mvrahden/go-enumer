package gen

import (
	"github.com/mvrahden/go-enumer/pkg/common"
)

type File struct {
	Header    Header
	Imports   []*Import
	TypeSpecs []*common.EnumType
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
