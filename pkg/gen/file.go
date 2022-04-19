package gen

import "github.com/mvrahden/go-enumer/config"

type File struct {
	Header    Header
	Imports   []*Import
	TypeSpecs []*TypeSpec
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

type TypeSpec struct {
	Index             int
	Name              string
	Docstring         string
	Type              GoType
	ValueSpecs        []*ValueSpec
	Filepath          string
	Config            *config.Options
	IsFromCsvSource   bool
	HasAdditionalData bool
	DataColumns       []string // contains additional column header
}

type ValueSpec struct {
	Value          uint64 // The numeric value of an enum constant
	ValueString    string // String representation of Value
	IdentifierName string
	EnumValue      string
	DataCells      []string // contains additional column values
}
