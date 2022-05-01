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
	Name           string
	Type           GoType
	Meta           *MetaTypeSpec
	ValueSpecs     []*ValueSpec
	AdditionalData *AdditionalData
}

type MetaTypeSpec struct {
	Docstring         string
	Config            *config.Options
	Filepath          string
	IsFromCsvSource   bool
	HasAdditionalData bool
}

type ValueSpec struct {
	ConstName          string
	Value              uint64 // hint: numeric value of an enum value/constant
	String             string // hint: the enum's actual value string
	IsAlternativeValue bool   // hint: is the enum an alternative value
}

type AdditionalData struct {
	Headers []AdditionalDataHeader
	Rows    [][]AdditionalDataCell
}

type AdditionalDataHeader struct {
	Name string // hint: the column name as-is (from CSV)
	Type GoType // hint: the type infered by type syntax
}
type AdditionalDataCell struct {
	LiteralValue string // hint: source representation of the value, e.g. literal strings are quoted
	RawValue     any    // hint: parsed value; actual type depends on header type
}
