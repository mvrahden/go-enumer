package gen

type File struct {
	Header     Header
	Imports    []*Import
	ValueSpecs []*ValueSpec
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

type ValueSpec struct {
	Index                      int
	IdentifierName, EnumString string
	Type                       GoType
	Value                      uint64 // sign is infered by value/type combination
	ValueString                string
}
