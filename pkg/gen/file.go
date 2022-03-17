package gen

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
	Index         int
	Name          string
	Docstring     string
	Type          GoType
	ValueSpecs    []*ValueSpec
	Filepath      string
	VirtualValues []*VirtualValue
}

type VirtualValue struct {
	Index          int
	Value          uint64
	ValueString    string
	EnumValue      string
	CanonicalValue string
}

type ValueSpec struct {
	IdentifierName, EnumString string
	Type                       GoType
	Value                      uint64 // sign is infered by value/type combination
	ValueString                string
}
