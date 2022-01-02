// Code generated by "%s"; DO NOT EDIT.

package pills

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
)

const (
	_PillString      = "PLACEBOASPIRINIBUPROFENPARACETAMOLACETAMINOPHENVITAMIN-C"
	_PillLowerString = "placeboaspirinibuprofenparacetamolacetaminophenvitamin-c"
)

var (
	_PillValueRange = [2]Pill{0, 4}

	_PillValues = []Pill{0, 1, 2, 3, 4}

	_PillStrings = []string{_PillString[0:7], _PillString[7:14], _PillString[14:23], _PillString[23:34], _PillString[47:56]}
)

// PillValues returns all values of the enum.
func PillValues() []Pill {
	strs := make([]Pill, len(_PillValues))
	copy(strs, _PillValues)
	return _PillValues
}

// PillStrings returns a slice of all String values of the enum.
func PillStrings() []string {
	strs := make([]string, len(_PillStrings))
	copy(strs, _PillStrings)
	return strs
}

// IsValid inspects whether the value is valid enum value.
func (_p Pill) IsValid() bool {
	return _p >= _PillValueRange[0] && _p <= _PillValueRange[1]
}

// String returns the string of the enum value.
// If the enum value is invalid.
func (_p Pill) String() string {
	if !_p.IsValid() {
		return fmt.Sprintf("Pill(%d)", _p)
	}
	idx := int(_p)
	return _PillStrings[idx]
}

var (
	_PillStringToValueMap = map[string]Pill{
		_PillString[0:7]:   PillPlacebo,
		_PillString[7:14]:  PillAspirin,
		_PillString[14:23]: PillIbuprofen,
		_PillString[23:34]: PillParacetamol,
		_PillString[34:47]: PillAcetaminophen,
		_PillString[47:56]: PillVitaminC,
	}
	_PillLowerStringToValueMap = map[string]Pill{
		_PillLowerString[0:7]:   PillPlacebo,
		_PillLowerString[7:14]:  PillAspirin,
		_PillLowerString[14:23]: PillIbuprofen,
		_PillLowerString[23:34]: PillParacetamol,
		_PillLowerString[34:47]: PillAcetaminophen,
		_PillLowerString[47:56]: PillVitaminC,
	}
)

// PillFromString determines the enum value with an exact case match.
func PillFromString(raw string) (Pill, bool) {
	v, ok := _PillStringToValueMap[raw]
	if !ok {
		return Pill(0), false
	}
	return v, true
}

// PillFromStringIgnoreCase determines the enum value with an case-insensitive match.
func PillFromStringIgnoreCase(raw string) (Pill, bool) {
	v, ok := PillFromString(raw)
	if ok {
		return v, ok
	}
	v, ok = _PillLowerStringToValueMap[raw]
	if !ok {
		return Pill(0), false
	}
	return v, true
}

// MarshalBinary implements the encoding.BinaryMarshaler interface for Pill.
func (_p Pill) MarshalBinary() ([]byte, error) {
	if !_p.IsValid() {
		return nil, fmt.Errorf("Cannot marshal invalid value %q as Pill", _p.String())
	}
	return []byte(_p.String()), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface for Pill.
func (_p *Pill) UnmarshalBinary(text []byte) error {
	str := string(text)
	if len(str) == 0 {
		return fmt.Errorf("Pill cannot be derived from empty string")
	}

	var ok bool
	*_p, ok = PillFromString(str)
	if !ok {
		return fmt.Errorf("Value %q does not represent a Pill", str)
	}
	return nil
}

// MarshalJSON implements the json.Marshaler interface for Pill.
func (_p Pill) MarshalJSON() ([]byte, error) {
	if !_p.IsValid() {
		return nil, fmt.Errorf("Cannot marshal invalid value %q as Pill", _p.String())
	}
	return json.Marshal(_p.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for Pill.
func (_p *Pill) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return fmt.Errorf("Pill should be a string, got %q", data)
	}
	if len(str) == 0 {
		return fmt.Errorf("Pill cannot be derived from empty string")
	}

	var ok bool
	*_p, ok = PillFromString(str)
	if !ok {
		return fmt.Errorf("Value %q does not represent a Pill", str)
	}
	return nil
}

func (_p Pill) Value() (driver.Value, error) {
	if !_p.IsValid() {
		return nil, fmt.Errorf("Cannot serialize invalid value %q as Pill", _p.String())
	}
	return _p.String(), nil
}

func (_p *Pill) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	var str string
	switch v := value.(type) {
	case []byte:
		str = string(v)
	case string:
		str = v
	case fmt.Stringer:
		str = v.String()
	default:
		return fmt.Errorf("invalid value of Pill: %[1]T(%[1]v)", value)
	}
	if len(str) == 0 {
		return fmt.Errorf("Pill cannot be derived from empty string")
	}

	var ok bool
	*_p, ok = PillFromString(str)
	if !ok {
		return fmt.Errorf("Value %q does not represent a Pill", str)
	}
	return nil
}

// MarshalText implements the encoding.TextMarshaler interface for Pill.
func (_p Pill) MarshalText() ([]byte, error) {
	if !_p.IsValid() {
		return nil, fmt.Errorf("Cannot marshal invalid value %q as Pill", _p.String())
	}
	return []byte(_p.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for Pill.
func (_p *Pill) UnmarshalText(text []byte) error {
	str := string(text)
	if len(str) == 0 {
		return fmt.Errorf("Pill cannot be derived from empty string")
	}

	var ok bool
	*_p, ok = PillFromString(str)
	if !ok {
		return fmt.Errorf("Value %q does not represent a Pill", str)
	}
	return nil
}

// MarshalYAML implements a YAML Marshaler for Pill
func (_p Pill) MarshalYAML() (interface{}, error) {
	if !_p.IsValid() {
		return nil, fmt.Errorf("Cannot marshal invalid value %q as Pill", _p.String())
	}
	return _p.String(), nil
}

// UnmarshalYAML implements a YAML Unmarshaler for Pill
func (_p *Pill) UnmarshalYAML(n *yaml.Node) error {
	const stringTag = "!!str"
	if n.ShortTag() != stringTag {
		return fmt.Errorf("Pill must be derived from a string node")
	}
	str := n.Value
	if len(str) == 0 {
		return fmt.Errorf("Pill cannot be derived from empty string")
	}

	var ok bool
	*_p, ok = PillFromString(str)
	if !ok {
		return fmt.Errorf("Value %q does not represent a Pill", str)
	}
	return nil
}
