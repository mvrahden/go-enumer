// Code generated by "%s"; DO NOT EDIT.

package planets

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

const (
	_PlanetString      = "MarsPlutoVenusMercuryJupiterSaturnUranusNeptune"
	_PlanetLowerString = "marsplutovenusmercuryjupitersaturnuranusneptune"
)

var (
	_PlanetValueRange = [2]Planet{1, 8}
	_PlanetValues     = []Planet{1, 2, 3, 4, 5, 6, 7, 8}
	_PlanetStrings    = []string{_PlanetString[0:4], _PlanetString[4:9], _PlanetString[9:14], _PlanetString[14:21], _PlanetString[21:28], _PlanetString[28:34], _PlanetString[34:40], _PlanetString[40:47]}
)

// _PlanetNoOp is a compile time assertion.
// An "invalid argument/out of bounds" compiler error signifies that the enum values have changed.
// Re-run the enumer command to generate an updated version of Planet.
func _PlanetNoOp() {
	var x [1]struct{}
	_ = x[PlanetMars-(1)]
	_ = x[PlanetPluto-(2)]
	_ = x[PlanetVenus-(3)]
	_ = x[PlanetMercury-(4)]
	_ = x[PlanetJupiter-(5)]
	_ = x[PlanetSaturn-(6)]
	_ = x[PlanetUranus-(7)]
	_ = x[PlanetNeptune-(8)]
}

// PlanetValues returns all values of the enum.
func PlanetValues() []Planet {
	strs := make([]Planet, len(_PlanetValues))
	copy(strs, _PlanetValues)
	return _PlanetValues
}

// PlanetStrings returns a slice of all String values of the enum.
func PlanetStrings() []string {
	strs := make([]string, len(_PlanetStrings))
	copy(strs, _PlanetStrings)
	return strs
}

// IsValid inspects whether the value is valid enum value.
func (_p Planet) IsValid() bool {
	return _p >= _PlanetValueRange[0] && _p <= _PlanetValueRange[1]
}

// String returns the string of the enum value.
// If the enum value is invalid, it will produce a string
// of the following pattern Planet(%d) instead.
func (_p Planet) String() string {
	if !_p.IsValid() {
		return fmt.Sprintf("Planet(%d)", _p)
	}
	idx := uint(_p) - 1
	return _PlanetStrings[idx]
}

var (
	_PlanetStringToValueMap = map[string]Planet{
		_PlanetString[0:4]:   PlanetMars,
		_PlanetString[4:9]:   PlanetPluto,
		_PlanetString[9:14]:  PlanetVenus,
		_PlanetString[14:21]: PlanetMercury,
		_PlanetString[21:28]: PlanetJupiter,
		_PlanetString[28:34]: PlanetSaturn,
		_PlanetString[34:40]: PlanetUranus,
		_PlanetString[40:47]: PlanetNeptune,
	}
	_PlanetLowerStringToValueMap = map[string]Planet{
		_PlanetLowerString[0:4]:   PlanetMars,
		_PlanetLowerString[4:9]:   PlanetPluto,
		_PlanetLowerString[9:14]:  PlanetVenus,
		_PlanetLowerString[14:21]: PlanetMercury,
		_PlanetLowerString[21:28]: PlanetJupiter,
		_PlanetLowerString[28:34]: PlanetSaturn,
		_PlanetLowerString[34:40]: PlanetUranus,
		_PlanetLowerString[40:47]: PlanetNeptune,
	}
)

// PlanetFromString determines the enum value with an exact case match.
func PlanetFromString(raw string) (Planet, bool) {
	v, ok := _PlanetStringToValueMap[raw]
	if !ok {
		return Planet(0), false
	}
	return v, true
}

// PlanetFromStringIgnoreCase determines the enum value with a case-insensitive match.
func PlanetFromStringIgnoreCase(raw string) (Planet, bool) {
	v, ok := PlanetFromString(raw)
	if ok {
		return v, ok
	}
	v, ok = _PlanetLowerStringToValueMap[raw]
	if !ok {
		return Planet(0), false
	}
	return v, true
}

// MarshalBinary implements the encoding.BinaryMarshaler interface for Planet.
func (_p Planet) MarshalBinary() ([]byte, error) {
	if !_p.IsValid() {
		return nil, fmt.Errorf("Cannot marshal invalid value %q as Planet", _p)
	}
	return []byte(_p.String()), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface for Planet.
func (_p *Planet) UnmarshalBinary(text []byte) error {
	str := string(text)
	if len(str) == 0 {
		return fmt.Errorf("Planet cannot be derived from empty string")
	}

	var ok bool
	*_p, ok = PlanetFromString(str)
	if !ok {
		return fmt.Errorf("Value %q does not represent a Planet", str)
	}
	return nil
}

// MarshalGQL implements the graphql.Marshaler interface for Planet.
func (_p Planet) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(_p.String()))
}

// UnmarshalGQL implements the graphql.Unmarshaler interface for Planet.
func (_p *Planet) UnmarshalGQL(value interface{}) error {
	var str string
	switch v := value.(type) {
	case nil:
	case []byte:
		str = string(v)
	case string:
		str = v
	case fmt.Stringer:
		str = v.String()
	default:
		return fmt.Errorf("invalid value of Planet: %[1]T(%[1]v)", value)
	}
	if len(str) == 0 {
		return fmt.Errorf("Planet cannot be derived from empty string")
	}

	var ok bool
	*_p, ok = PlanetFromString(str)
	if !ok {
		return fmt.Errorf("Value %q does not represent a Planet", str)
	}
	return nil
}

// MarshalJSON implements the json.Marshaler interface for Planet.
func (_p Planet) MarshalJSON() ([]byte, error) {
	if !_p.IsValid() {
		return nil, fmt.Errorf("Cannot marshal invalid value %q as Planet", _p)
	}
	return json.Marshal(_p.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for Planet.
func (_p *Planet) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return fmt.Errorf("Planet should be a string, got %q", data)
	}
	if len(str) == 0 {
		return fmt.Errorf("Planet cannot be derived from empty string")
	}

	var ok bool
	*_p, ok = PlanetFromString(str)
	if !ok {
		return fmt.Errorf("Value %q does not represent a Planet", str)
	}
	return nil
}

func (_p Planet) Value() (driver.Value, error) {
	if !_p.IsValid() {
		return nil, fmt.Errorf("Cannot serialize invalid value %q as Planet", _p)
	}
	return _p.String(), nil
}

func (_p *Planet) Scan(value interface{}) error {
	var str string
	switch v := value.(type) {
	case nil:
	case []byte:
		str = string(v)
	case string:
		str = v
	case fmt.Stringer:
		str = v.String()
	default:
		return fmt.Errorf("invalid value of Planet: %[1]T(%[1]v)", value)
	}
	if len(str) == 0 {
		return fmt.Errorf("Planet cannot be derived from empty string")
	}

	var ok bool
	*_p, ok = PlanetFromString(str)
	if !ok {
		return fmt.Errorf("Value %q does not represent a Planet", str)
	}
	return nil
}

// MarshalText implements the encoding.TextMarshaler interface for Planet.
func (_p Planet) MarshalText() ([]byte, error) {
	if !_p.IsValid() {
		return nil, fmt.Errorf("Cannot marshal invalid value %q as Planet", _p)
	}
	return []byte(_p.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for Planet.
func (_p *Planet) UnmarshalText(text []byte) error {
	str := string(text)
	if len(str) == 0 {
		return fmt.Errorf("Planet cannot be derived from empty string")
	}

	var ok bool
	*_p, ok = PlanetFromString(str)
	if !ok {
		return fmt.Errorf("Value %q does not represent a Planet", str)
	}
	return nil
}

// MarshalYAML implements a YAML Marshaler for Planet.
func (_p Planet) MarshalYAML() (interface{}, error) {
	if !_p.IsValid() {
		return nil, fmt.Errorf("Cannot marshal invalid value %q as Planet", _p)
	}
	return _p.String(), nil
}

// UnmarshalYAML implements a YAML Unmarshaler for Planet.
func (_p *Planet) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var str string
	if err := unmarshal(&str); err != nil {
		return err
	}
	if len(str) == 0 {
		return fmt.Errorf("Planet cannot be derived from empty string")
	}

	var ok bool
	*_p, ok = PlanetFromString(str)
	if !ok {
		return fmt.Errorf("Value %q does not represent a Planet", str)
	}
	return nil
}

const (
	_PlanetWithDefaultString      = "EarthMarsPlutoVenusMercuryJupiterSaturnUranusNeptune"
	_PlanetWithDefaultLowerString = "earthmarsplutovenusmercuryjupitersaturnuranusneptune"
)

var (
	_PlanetWithDefaultValueRange = [2]PlanetWithDefault{0, 8}
	_PlanetWithDefaultValues     = []PlanetWithDefault{0, 1, 2, 3, 4, 5, 6, 7, 8}
	_PlanetWithDefaultStrings    = []string{_PlanetWithDefaultString[0:5], _PlanetWithDefaultString[5:9], _PlanetWithDefaultString[9:14], _PlanetWithDefaultString[14:19], _PlanetWithDefaultString[19:26], _PlanetWithDefaultString[26:33], _PlanetWithDefaultString[33:39], _PlanetWithDefaultString[39:45], _PlanetWithDefaultString[45:52]}
)

// _PlanetWithDefaultNoOp is a compile time assertion.
// An "invalid argument/out of bounds" compiler error signifies that the enum values have changed.
// Re-run the enumer command to generate an updated version of PlanetWithDefault.
func _PlanetWithDefaultNoOp() {
	var x [1]struct{}
	_ = x[PlanetWithDefaultEarth-(0)]
	_ = x[PlanetWithDefaultMars-(1)]
	_ = x[PlanetWithDefaultPluto-(2)]
	_ = x[PlanetWithDefaultVenus-(3)]
	_ = x[PlanetWithDefaultMercury-(4)]
	_ = x[PlanetWithDefaultJupiter-(5)]
	_ = x[PlanetWithDefaultSaturn-(6)]
	_ = x[PlanetWithDefaultUranus-(7)]
	_ = x[PlanetWithDefaultNeptune-(8)]
}

// PlanetWithDefaultValues returns all values of the enum.
func PlanetWithDefaultValues() []PlanetWithDefault {
	strs := make([]PlanetWithDefault, len(_PlanetWithDefaultValues))
	copy(strs, _PlanetWithDefaultValues)
	return _PlanetWithDefaultValues
}

// PlanetWithDefaultStrings returns a slice of all String values of the enum.
func PlanetWithDefaultStrings() []string {
	strs := make([]string, len(_PlanetWithDefaultStrings))
	copy(strs, _PlanetWithDefaultStrings)
	return strs
}

// IsValid inspects whether the value is valid enum value.
func (_p PlanetWithDefault) IsValid() bool {
	return _p >= _PlanetWithDefaultValueRange[0] && _p <= _PlanetWithDefaultValueRange[1]
}

// String returns the string of the enum value.
// If the enum value is invalid, it will produce a string
// of the following pattern PlanetWithDefault(%d) instead.
func (_p PlanetWithDefault) String() string {
	if !_p.IsValid() {
		return fmt.Sprintf("PlanetWithDefault(%d)", _p)
	}
	idx := uint(_p)
	return _PlanetWithDefaultStrings[idx]
}

var (
	_PlanetWithDefaultStringToValueMap = map[string]PlanetWithDefault{
		_PlanetWithDefaultString[0:5]:   PlanetWithDefaultEarth,
		_PlanetWithDefaultString[5:9]:   PlanetWithDefaultMars,
		_PlanetWithDefaultString[9:14]:  PlanetWithDefaultPluto,
		_PlanetWithDefaultString[14:19]: PlanetWithDefaultVenus,
		_PlanetWithDefaultString[19:26]: PlanetWithDefaultMercury,
		_PlanetWithDefaultString[26:33]: PlanetWithDefaultJupiter,
		_PlanetWithDefaultString[33:39]: PlanetWithDefaultSaturn,
		_PlanetWithDefaultString[39:45]: PlanetWithDefaultUranus,
		_PlanetWithDefaultString[45:52]: PlanetWithDefaultNeptune,
	}
	_PlanetWithDefaultLowerStringToValueMap = map[string]PlanetWithDefault{
		_PlanetWithDefaultLowerString[0:5]:   PlanetWithDefaultEarth,
		_PlanetWithDefaultLowerString[5:9]:   PlanetWithDefaultMars,
		_PlanetWithDefaultLowerString[9:14]:  PlanetWithDefaultPluto,
		_PlanetWithDefaultLowerString[14:19]: PlanetWithDefaultVenus,
		_PlanetWithDefaultLowerString[19:26]: PlanetWithDefaultMercury,
		_PlanetWithDefaultLowerString[26:33]: PlanetWithDefaultJupiter,
		_PlanetWithDefaultLowerString[33:39]: PlanetWithDefaultSaturn,
		_PlanetWithDefaultLowerString[39:45]: PlanetWithDefaultUranus,
		_PlanetWithDefaultLowerString[45:52]: PlanetWithDefaultNeptune,
	}
)

// PlanetWithDefaultFromString determines the enum value with an exact case match.
func PlanetWithDefaultFromString(raw string) (PlanetWithDefault, bool) {
	v, ok := _PlanetWithDefaultStringToValueMap[raw]
	if !ok {
		return PlanetWithDefault(0), false
	}
	return v, true
}

// PlanetWithDefaultFromStringIgnoreCase determines the enum value with a case-insensitive match.
func PlanetWithDefaultFromStringIgnoreCase(raw string) (PlanetWithDefault, bool) {
	v, ok := PlanetWithDefaultFromString(raw)
	if ok {
		return v, ok
	}
	v, ok = _PlanetWithDefaultLowerStringToValueMap[raw]
	if !ok {
		return PlanetWithDefault(0), false
	}
	return v, true
}

// MarshalBinary implements the encoding.BinaryMarshaler interface for PlanetWithDefault.
func (_p PlanetWithDefault) MarshalBinary() ([]byte, error) {
	if !_p.IsValid() {
		return nil, fmt.Errorf("Cannot marshal invalid value %q as PlanetWithDefault", _p)
	}
	return []byte(_p.String()), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface for PlanetWithDefault.
func (_p *PlanetWithDefault) UnmarshalBinary(text []byte) error {
	str := string(text)
	if len(str) == 0 {
		return fmt.Errorf("PlanetWithDefault cannot be derived from empty string")
	}

	var ok bool
	*_p, ok = PlanetWithDefaultFromString(str)
	if !ok {
		return fmt.Errorf("Value %q does not represent a PlanetWithDefault", str)
	}
	return nil
}

// MarshalGQL implements the graphql.Marshaler interface for PlanetWithDefault.
func (_p PlanetWithDefault) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(_p.String()))
}

// UnmarshalGQL implements the graphql.Unmarshaler interface for PlanetWithDefault.
func (_p *PlanetWithDefault) UnmarshalGQL(value interface{}) error {
	var str string
	switch v := value.(type) {
	case nil:
	case []byte:
		str = string(v)
	case string:
		str = v
	case fmt.Stringer:
		str = v.String()
	default:
		return fmt.Errorf("invalid value of PlanetWithDefault: %[1]T(%[1]v)", value)
	}
	if len(str) == 0 {
		return fmt.Errorf("PlanetWithDefault cannot be derived from empty string")
	}

	var ok bool
	*_p, ok = PlanetWithDefaultFromString(str)
	if !ok {
		return fmt.Errorf("Value %q does not represent a PlanetWithDefault", str)
	}
	return nil
}

// MarshalJSON implements the json.Marshaler interface for PlanetWithDefault.
func (_p PlanetWithDefault) MarshalJSON() ([]byte, error) {
	if !_p.IsValid() {
		return nil, fmt.Errorf("Cannot marshal invalid value %q as PlanetWithDefault", _p)
	}
	return json.Marshal(_p.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for PlanetWithDefault.
func (_p *PlanetWithDefault) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return fmt.Errorf("PlanetWithDefault should be a string, got %q", data)
	}
	if len(str) == 0 {
		return fmt.Errorf("PlanetWithDefault cannot be derived from empty string")
	}

	var ok bool
	*_p, ok = PlanetWithDefaultFromString(str)
	if !ok {
		return fmt.Errorf("Value %q does not represent a PlanetWithDefault", str)
	}
	return nil
}

func (_p PlanetWithDefault) Value() (driver.Value, error) {
	if !_p.IsValid() {
		return nil, fmt.Errorf("Cannot serialize invalid value %q as PlanetWithDefault", _p)
	}
	return _p.String(), nil
}

func (_p *PlanetWithDefault) Scan(value interface{}) error {
	var str string
	switch v := value.(type) {
	case nil:
	case []byte:
		str = string(v)
	case string:
		str = v
	case fmt.Stringer:
		str = v.String()
	default:
		return fmt.Errorf("invalid value of PlanetWithDefault: %[1]T(%[1]v)", value)
	}
	if len(str) == 0 {
		return fmt.Errorf("PlanetWithDefault cannot be derived from empty string")
	}

	var ok bool
	*_p, ok = PlanetWithDefaultFromString(str)
	if !ok {
		return fmt.Errorf("Value %q does not represent a PlanetWithDefault", str)
	}
	return nil
}

// MarshalText implements the encoding.TextMarshaler interface for PlanetWithDefault.
func (_p PlanetWithDefault) MarshalText() ([]byte, error) {
	if !_p.IsValid() {
		return nil, fmt.Errorf("Cannot marshal invalid value %q as PlanetWithDefault", _p)
	}
	return []byte(_p.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for PlanetWithDefault.
func (_p *PlanetWithDefault) UnmarshalText(text []byte) error {
	str := string(text)
	if len(str) == 0 {
		return fmt.Errorf("PlanetWithDefault cannot be derived from empty string")
	}

	var ok bool
	*_p, ok = PlanetWithDefaultFromString(str)
	if !ok {
		return fmt.Errorf("Value %q does not represent a PlanetWithDefault", str)
	}
	return nil
}

// MarshalYAML implements a YAML Marshaler for PlanetWithDefault.
func (_p PlanetWithDefault) MarshalYAML() (interface{}, error) {
	if !_p.IsValid() {
		return nil, fmt.Errorf("Cannot marshal invalid value %q as PlanetWithDefault", _p)
	}
	return _p.String(), nil
}

// UnmarshalYAML implements a YAML Unmarshaler for PlanetWithDefault.
func (_p *PlanetWithDefault) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var str string
	if err := unmarshal(&str); err != nil {
		return err
	}
	if len(str) == 0 {
		return fmt.Errorf("PlanetWithDefault cannot be derived from empty string")
	}

	var ok bool
	*_p, ok = PlanetWithDefaultFromString(str)
	if !ok {
		return fmt.Errorf("Value %q does not represent a PlanetWithDefault", str)
	}
	return nil
}

const (
	_PlanetSupportUndefinedString      = "MarsPlutoVenusMercuryJupiterSaturnUranusNeptune"
	_PlanetSupportUndefinedLowerString = "marsplutovenusmercuryjupitersaturnuranusneptune"
)

const (
	// PlanetSupportUndefinedUndefined is the generated zero value of the PlanetSupportUndefined enum.
	PlanetSupportUndefinedUndefined PlanetSupportUndefined = 0
)

var (
	_PlanetSupportUndefinedValueRange = [2]PlanetSupportUndefined{0, 8}
	_PlanetSupportUndefinedValues     = []PlanetSupportUndefined{1, 2, 3, 4, 5, 6, 7, 8}
	_PlanetSupportUndefinedStrings    = []string{_PlanetSupportUndefinedString[0:4], _PlanetSupportUndefinedString[4:9], _PlanetSupportUndefinedString[9:14], _PlanetSupportUndefinedString[14:21], _PlanetSupportUndefinedString[21:28], _PlanetSupportUndefinedString[28:34], _PlanetSupportUndefinedString[34:40], _PlanetSupportUndefinedString[40:47]}
)

// _PlanetSupportUndefinedNoOp is a compile time assertion.
// An "invalid argument/out of bounds" compiler error signifies that the enum values have changed.
// Re-run the enumer command to generate an updated version of PlanetSupportUndefined.
func _PlanetSupportUndefinedNoOp() {
	var x [1]struct{}
	_ = x[PlanetSupportUndefinedUndefined-(0)]
	_ = x[PlanetSupportUndefinedMars-(1)]
	_ = x[PlanetSupportUndefinedPluto-(2)]
	_ = x[PlanetSupportUndefinedVenus-(3)]
	_ = x[PlanetSupportUndefinedMercury-(4)]
	_ = x[PlanetSupportUndefinedJupiter-(5)]
	_ = x[PlanetSupportUndefinedSaturn-(6)]
	_ = x[PlanetSupportUndefinedUranus-(7)]
	_ = x[PlanetSupportUndefinedNeptune-(8)]
}

// PlanetSupportUndefinedValues returns all values of the enum.
func PlanetSupportUndefinedValues() []PlanetSupportUndefined {
	strs := make([]PlanetSupportUndefined, len(_PlanetSupportUndefinedValues))
	copy(strs, _PlanetSupportUndefinedValues)
	return _PlanetSupportUndefinedValues
}

// PlanetSupportUndefinedStrings returns a slice of all String values of the enum.
func PlanetSupportUndefinedStrings() []string {
	strs := make([]string, len(_PlanetSupportUndefinedStrings))
	copy(strs, _PlanetSupportUndefinedStrings)
	return strs
}

// IsValid inspects whether the value is valid enum value.
func (_p PlanetSupportUndefined) IsValid() bool {
	return _p >= _PlanetSupportUndefinedValueRange[0] && _p <= _PlanetSupportUndefinedValueRange[1]
}

// String returns the string of the enum value.
// If the enum value is invalid, it will produce a string
// of the following pattern PlanetSupportUndefined(%d) instead.
func (_p PlanetSupportUndefined) String() string {
	if !_p.IsValid() {
		return fmt.Sprintf("PlanetSupportUndefined(%d)", _p)
	}
	if _p == PlanetSupportUndefinedUndefined {
		return ""
	}
	idx := uint(_p) - 1
	return _PlanetSupportUndefinedStrings[idx]
}

var (
	_PlanetSupportUndefinedStringToValueMap = map[string]PlanetSupportUndefined{
		_PlanetSupportUndefinedString[0:4]:   PlanetSupportUndefinedMars,
		_PlanetSupportUndefinedString[4:9]:   PlanetSupportUndefinedPluto,
		_PlanetSupportUndefinedString[9:14]:  PlanetSupportUndefinedVenus,
		_PlanetSupportUndefinedString[14:21]: PlanetSupportUndefinedMercury,
		_PlanetSupportUndefinedString[21:28]: PlanetSupportUndefinedJupiter,
		_PlanetSupportUndefinedString[28:34]: PlanetSupportUndefinedSaturn,
		_PlanetSupportUndefinedString[34:40]: PlanetSupportUndefinedUranus,
		_PlanetSupportUndefinedString[40:47]: PlanetSupportUndefinedNeptune,
	}
	_PlanetSupportUndefinedLowerStringToValueMap = map[string]PlanetSupportUndefined{
		_PlanetSupportUndefinedLowerString[0:4]:   PlanetSupportUndefinedMars,
		_PlanetSupportUndefinedLowerString[4:9]:   PlanetSupportUndefinedPluto,
		_PlanetSupportUndefinedLowerString[9:14]:  PlanetSupportUndefinedVenus,
		_PlanetSupportUndefinedLowerString[14:21]: PlanetSupportUndefinedMercury,
		_PlanetSupportUndefinedLowerString[21:28]: PlanetSupportUndefinedJupiter,
		_PlanetSupportUndefinedLowerString[28:34]: PlanetSupportUndefinedSaturn,
		_PlanetSupportUndefinedLowerString[34:40]: PlanetSupportUndefinedUranus,
		_PlanetSupportUndefinedLowerString[40:47]: PlanetSupportUndefinedNeptune,
	}
)

// PlanetSupportUndefinedFromString determines the enum value with an exact case match.
func PlanetSupportUndefinedFromString(raw string) (PlanetSupportUndefined, bool) {
	if len(raw) == 0 {
		return PlanetSupportUndefined(0), true
	}
	v, ok := _PlanetSupportUndefinedStringToValueMap[raw]
	if !ok {
		return PlanetSupportUndefined(0), false
	}
	return v, true
}

// PlanetSupportUndefinedFromStringIgnoreCase determines the enum value with a case-insensitive match.
func PlanetSupportUndefinedFromStringIgnoreCase(raw string) (PlanetSupportUndefined, bool) {
	v, ok := PlanetSupportUndefinedFromString(raw)
	if ok {
		return v, ok
	}
	v, ok = _PlanetSupportUndefinedLowerStringToValueMap[raw]
	if !ok {
		return PlanetSupportUndefined(0), false
	}
	return v, true
}

// MarshalBinary implements the encoding.BinaryMarshaler interface for PlanetSupportUndefined.
func (_p PlanetSupportUndefined) MarshalBinary() ([]byte, error) {
	if !_p.IsValid() {
		return nil, fmt.Errorf("Cannot marshal invalid value %q as PlanetSupportUndefined", _p)
	}
	return []byte(_p.String()), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface for PlanetSupportUndefined.
func (_p *PlanetSupportUndefined) UnmarshalBinary(text []byte) error {
	str := string(text)

	var ok bool
	*_p, ok = PlanetSupportUndefinedFromString(str)
	if !ok {
		return fmt.Errorf("Value %q does not represent a PlanetSupportUndefined", str)
	}
	return nil
}

// MarshalGQL implements the graphql.Marshaler interface for PlanetSupportUndefined.
func (_p PlanetSupportUndefined) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(_p.String()))
}

// UnmarshalGQL implements the graphql.Unmarshaler interface for PlanetSupportUndefined.
func (_p *PlanetSupportUndefined) UnmarshalGQL(value interface{}) error {
	var str string
	switch v := value.(type) {
	case nil:
	case []byte:
		str = string(v)
	case string:
		str = v
	case fmt.Stringer:
		str = v.String()
	default:
		return fmt.Errorf("invalid value of PlanetSupportUndefined: %[1]T(%[1]v)", value)
	}

	var ok bool
	*_p, ok = PlanetSupportUndefinedFromString(str)
	if !ok {
		return fmt.Errorf("Value %q does not represent a PlanetSupportUndefined", str)
	}
	return nil
}

// MarshalJSON implements the json.Marshaler interface for PlanetSupportUndefined.
func (_p PlanetSupportUndefined) MarshalJSON() ([]byte, error) {
	if !_p.IsValid() {
		return nil, fmt.Errorf("Cannot marshal invalid value %q as PlanetSupportUndefined", _p)
	}
	return json.Marshal(_p.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for PlanetSupportUndefined.
func (_p *PlanetSupportUndefined) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return fmt.Errorf("PlanetSupportUndefined should be a string, got %q", data)
	}

	var ok bool
	*_p, ok = PlanetSupportUndefinedFromString(str)
	if !ok {
		return fmt.Errorf("Value %q does not represent a PlanetSupportUndefined", str)
	}
	return nil
}

func (_p PlanetSupportUndefined) Value() (driver.Value, error) {
	if !_p.IsValid() {
		return nil, fmt.Errorf("Cannot serialize invalid value %q as PlanetSupportUndefined", _p)
	}
	return _p.String(), nil
}

func (_p *PlanetSupportUndefined) Scan(value interface{}) error {
	var str string
	switch v := value.(type) {
	case nil:
	case []byte:
		str = string(v)
	case string:
		str = v
	case fmt.Stringer:
		str = v.String()
	default:
		return fmt.Errorf("invalid value of PlanetSupportUndefined: %[1]T(%[1]v)", value)
	}

	var ok bool
	*_p, ok = PlanetSupportUndefinedFromString(str)
	if !ok {
		return fmt.Errorf("Value %q does not represent a PlanetSupportUndefined", str)
	}
	return nil
}

// MarshalText implements the encoding.TextMarshaler interface for PlanetSupportUndefined.
func (_p PlanetSupportUndefined) MarshalText() ([]byte, error) {
	if !_p.IsValid() {
		return nil, fmt.Errorf("Cannot marshal invalid value %q as PlanetSupportUndefined", _p)
	}
	return []byte(_p.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for PlanetSupportUndefined.
func (_p *PlanetSupportUndefined) UnmarshalText(text []byte) error {
	str := string(text)

	var ok bool
	*_p, ok = PlanetSupportUndefinedFromString(str)
	if !ok {
		return fmt.Errorf("Value %q does not represent a PlanetSupportUndefined", str)
	}
	return nil
}

// MarshalYAML implements a YAML Marshaler for PlanetSupportUndefined.
func (_p PlanetSupportUndefined) MarshalYAML() (interface{}, error) {
	if !_p.IsValid() {
		return nil, fmt.Errorf("Cannot marshal invalid value %q as PlanetSupportUndefined", _p)
	}
	return _p.String(), nil
}

// UnmarshalYAML implements a YAML Unmarshaler for PlanetSupportUndefined.
func (_p *PlanetSupportUndefined) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var str string
	if err := unmarshal(&str); err != nil {
		return err
	}

	var ok bool
	*_p, ok = PlanetSupportUndefinedFromString(str)
	if !ok {
		return fmt.Errorf("Value %q does not represent a PlanetSupportUndefined", str)
	}
	return nil
}

const (
	_PlanetSupportUndefinedWithDefaultString      = "EarthMarsPlutoVenusMercuryJupiterSaturnUranusNeptune"
	_PlanetSupportUndefinedWithDefaultLowerString = "earthmarsplutovenusmercuryjupitersaturnuranusneptune"
)

var (
	_PlanetSupportUndefinedWithDefaultValueRange = [2]PlanetSupportUndefinedWithDefault{0, 8}
	_PlanetSupportUndefinedWithDefaultValues     = []PlanetSupportUndefinedWithDefault{0, 1, 2, 3, 4, 5, 6, 7, 8}
	_PlanetSupportUndefinedWithDefaultStrings    = []string{_PlanetSupportUndefinedWithDefaultString[0:5], _PlanetSupportUndefinedWithDefaultString[5:9], _PlanetSupportUndefinedWithDefaultString[9:14], _PlanetSupportUndefinedWithDefaultString[14:19], _PlanetSupportUndefinedWithDefaultString[19:26], _PlanetSupportUndefinedWithDefaultString[26:33], _PlanetSupportUndefinedWithDefaultString[33:39], _PlanetSupportUndefinedWithDefaultString[39:45], _PlanetSupportUndefinedWithDefaultString[45:52]}
)

// _PlanetSupportUndefinedWithDefaultNoOp is a compile time assertion.
// An "invalid argument/out of bounds" compiler error signifies that the enum values have changed.
// Re-run the enumer command to generate an updated version of PlanetSupportUndefinedWithDefault.
func _PlanetSupportUndefinedWithDefaultNoOp() {
	var x [1]struct{}
	_ = x[PlanetSupportUndefinedWithDefaultEarth-(0)]
	_ = x[PlanetSupportUndefinedWithDefaultMars-(1)]
	_ = x[PlanetSupportUndefinedWithDefaultPluto-(2)]
	_ = x[PlanetSupportUndefinedWithDefaultVenus-(3)]
	_ = x[PlanetSupportUndefinedWithDefaultMercury-(4)]
	_ = x[PlanetSupportUndefinedWithDefaultJupiter-(5)]
	_ = x[PlanetSupportUndefinedWithDefaultSaturn-(6)]
	_ = x[PlanetSupportUndefinedWithDefaultUranus-(7)]
	_ = x[PlanetSupportUndefinedWithDefaultNeptune-(8)]
}

// PlanetSupportUndefinedWithDefaultValues returns all values of the enum.
func PlanetSupportUndefinedWithDefaultValues() []PlanetSupportUndefinedWithDefault {
	strs := make([]PlanetSupportUndefinedWithDefault, len(_PlanetSupportUndefinedWithDefaultValues))
	copy(strs, _PlanetSupportUndefinedWithDefaultValues)
	return _PlanetSupportUndefinedWithDefaultValues
}

// PlanetSupportUndefinedWithDefaultStrings returns a slice of all String values of the enum.
func PlanetSupportUndefinedWithDefaultStrings() []string {
	strs := make([]string, len(_PlanetSupportUndefinedWithDefaultStrings))
	copy(strs, _PlanetSupportUndefinedWithDefaultStrings)
	return strs
}

// IsValid inspects whether the value is valid enum value.
func (_p PlanetSupportUndefinedWithDefault) IsValid() bool {
	return _p >= _PlanetSupportUndefinedWithDefaultValueRange[0] && _p <= _PlanetSupportUndefinedWithDefaultValueRange[1]
}

// String returns the string of the enum value.
// If the enum value is invalid, it will produce a string
// of the following pattern PlanetSupportUndefinedWithDefault(%d) instead.
func (_p PlanetSupportUndefinedWithDefault) String() string {
	if !_p.IsValid() {
		return fmt.Sprintf("PlanetSupportUndefinedWithDefault(%d)", _p)
	}
	idx := uint(_p)
	return _PlanetSupportUndefinedWithDefaultStrings[idx]
}

var (
	_PlanetSupportUndefinedWithDefaultStringToValueMap = map[string]PlanetSupportUndefinedWithDefault{
		_PlanetSupportUndefinedWithDefaultString[0:5]:   PlanetSupportUndefinedWithDefaultEarth,
		_PlanetSupportUndefinedWithDefaultString[5:9]:   PlanetSupportUndefinedWithDefaultMars,
		_PlanetSupportUndefinedWithDefaultString[9:14]:  PlanetSupportUndefinedWithDefaultPluto,
		_PlanetSupportUndefinedWithDefaultString[14:19]: PlanetSupportUndefinedWithDefaultVenus,
		_PlanetSupportUndefinedWithDefaultString[19:26]: PlanetSupportUndefinedWithDefaultMercury,
		_PlanetSupportUndefinedWithDefaultString[26:33]: PlanetSupportUndefinedWithDefaultJupiter,
		_PlanetSupportUndefinedWithDefaultString[33:39]: PlanetSupportUndefinedWithDefaultSaturn,
		_PlanetSupportUndefinedWithDefaultString[39:45]: PlanetSupportUndefinedWithDefaultUranus,
		_PlanetSupportUndefinedWithDefaultString[45:52]: PlanetSupportUndefinedWithDefaultNeptune,
	}
	_PlanetSupportUndefinedWithDefaultLowerStringToValueMap = map[string]PlanetSupportUndefinedWithDefault{
		_PlanetSupportUndefinedWithDefaultLowerString[0:5]:   PlanetSupportUndefinedWithDefaultEarth,
		_PlanetSupportUndefinedWithDefaultLowerString[5:9]:   PlanetSupportUndefinedWithDefaultMars,
		_PlanetSupportUndefinedWithDefaultLowerString[9:14]:  PlanetSupportUndefinedWithDefaultPluto,
		_PlanetSupportUndefinedWithDefaultLowerString[14:19]: PlanetSupportUndefinedWithDefaultVenus,
		_PlanetSupportUndefinedWithDefaultLowerString[19:26]: PlanetSupportUndefinedWithDefaultMercury,
		_PlanetSupportUndefinedWithDefaultLowerString[26:33]: PlanetSupportUndefinedWithDefaultJupiter,
		_PlanetSupportUndefinedWithDefaultLowerString[33:39]: PlanetSupportUndefinedWithDefaultSaturn,
		_PlanetSupportUndefinedWithDefaultLowerString[39:45]: PlanetSupportUndefinedWithDefaultUranus,
		_PlanetSupportUndefinedWithDefaultLowerString[45:52]: PlanetSupportUndefinedWithDefaultNeptune,
	}
)

// PlanetSupportUndefinedWithDefaultFromString determines the enum value with an exact case match.
func PlanetSupportUndefinedWithDefaultFromString(raw string) (PlanetSupportUndefinedWithDefault, bool) {
	if len(raw) == 0 {
		return PlanetSupportUndefinedWithDefault(0), true
	}
	v, ok := _PlanetSupportUndefinedWithDefaultStringToValueMap[raw]
	if !ok {
		return PlanetSupportUndefinedWithDefault(0), false
	}
	return v, true
}

// PlanetSupportUndefinedWithDefaultFromStringIgnoreCase determines the enum value with a case-insensitive match.
func PlanetSupportUndefinedWithDefaultFromStringIgnoreCase(raw string) (PlanetSupportUndefinedWithDefault, bool) {
	v, ok := PlanetSupportUndefinedWithDefaultFromString(raw)
	if ok {
		return v, ok
	}
	v, ok = _PlanetSupportUndefinedWithDefaultLowerStringToValueMap[raw]
	if !ok {
		return PlanetSupportUndefinedWithDefault(0), false
	}
	return v, true
}

// MarshalBinary implements the encoding.BinaryMarshaler interface for PlanetSupportUndefinedWithDefault.
func (_p PlanetSupportUndefinedWithDefault) MarshalBinary() ([]byte, error) {
	if !_p.IsValid() {
		return nil, fmt.Errorf("Cannot marshal invalid value %q as PlanetSupportUndefinedWithDefault", _p)
	}
	return []byte(_p.String()), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface for PlanetSupportUndefinedWithDefault.
func (_p *PlanetSupportUndefinedWithDefault) UnmarshalBinary(text []byte) error {
	str := string(text)

	var ok bool
	*_p, ok = PlanetSupportUndefinedWithDefaultFromString(str)
	if !ok {
		return fmt.Errorf("Value %q does not represent a PlanetSupportUndefinedWithDefault", str)
	}
	return nil
}

// MarshalGQL implements the graphql.Marshaler interface for PlanetSupportUndefinedWithDefault.
func (_p PlanetSupportUndefinedWithDefault) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(_p.String()))
}

// UnmarshalGQL implements the graphql.Unmarshaler interface for PlanetSupportUndefinedWithDefault.
func (_p *PlanetSupportUndefinedWithDefault) UnmarshalGQL(value interface{}) error {
	var str string
	switch v := value.(type) {
	case nil:
	case []byte:
		str = string(v)
	case string:
		str = v
	case fmt.Stringer:
		str = v.String()
	default:
		return fmt.Errorf("invalid value of PlanetSupportUndefinedWithDefault: %[1]T(%[1]v)", value)
	}

	var ok bool
	*_p, ok = PlanetSupportUndefinedWithDefaultFromString(str)
	if !ok {
		return fmt.Errorf("Value %q does not represent a PlanetSupportUndefinedWithDefault", str)
	}
	return nil
}

// MarshalJSON implements the json.Marshaler interface for PlanetSupportUndefinedWithDefault.
func (_p PlanetSupportUndefinedWithDefault) MarshalJSON() ([]byte, error) {
	if !_p.IsValid() {
		return nil, fmt.Errorf("Cannot marshal invalid value %q as PlanetSupportUndefinedWithDefault", _p)
	}
	return json.Marshal(_p.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for PlanetSupportUndefinedWithDefault.
func (_p *PlanetSupportUndefinedWithDefault) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return fmt.Errorf("PlanetSupportUndefinedWithDefault should be a string, got %q", data)
	}

	var ok bool
	*_p, ok = PlanetSupportUndefinedWithDefaultFromString(str)
	if !ok {
		return fmt.Errorf("Value %q does not represent a PlanetSupportUndefinedWithDefault", str)
	}
	return nil
}

func (_p PlanetSupportUndefinedWithDefault) Value() (driver.Value, error) {
	if !_p.IsValid() {
		return nil, fmt.Errorf("Cannot serialize invalid value %q as PlanetSupportUndefinedWithDefault", _p)
	}
	return _p.String(), nil
}

func (_p *PlanetSupportUndefinedWithDefault) Scan(value interface{}) error {
	var str string
	switch v := value.(type) {
	case nil:
	case []byte:
		str = string(v)
	case string:
		str = v
	case fmt.Stringer:
		str = v.String()
	default:
		return fmt.Errorf("invalid value of PlanetSupportUndefinedWithDefault: %[1]T(%[1]v)", value)
	}

	var ok bool
	*_p, ok = PlanetSupportUndefinedWithDefaultFromString(str)
	if !ok {
		return fmt.Errorf("Value %q does not represent a PlanetSupportUndefinedWithDefault", str)
	}
	return nil
}

// MarshalText implements the encoding.TextMarshaler interface for PlanetSupportUndefinedWithDefault.
func (_p PlanetSupportUndefinedWithDefault) MarshalText() ([]byte, error) {
	if !_p.IsValid() {
		return nil, fmt.Errorf("Cannot marshal invalid value %q as PlanetSupportUndefinedWithDefault", _p)
	}
	return []byte(_p.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for PlanetSupportUndefinedWithDefault.
func (_p *PlanetSupportUndefinedWithDefault) UnmarshalText(text []byte) error {
	str := string(text)

	var ok bool
	*_p, ok = PlanetSupportUndefinedWithDefaultFromString(str)
	if !ok {
		return fmt.Errorf("Value %q does not represent a PlanetSupportUndefinedWithDefault", str)
	}
	return nil
}

// MarshalYAML implements a YAML Marshaler for PlanetSupportUndefinedWithDefault.
func (_p PlanetSupportUndefinedWithDefault) MarshalYAML() (interface{}, error) {
	if !_p.IsValid() {
		return nil, fmt.Errorf("Cannot marshal invalid value %q as PlanetSupportUndefinedWithDefault", _p)
	}
	return _p.String(), nil
}

// UnmarshalYAML implements a YAML Unmarshaler for PlanetSupportUndefinedWithDefault.
func (_p *PlanetSupportUndefinedWithDefault) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var str string
	if err := unmarshal(&str); err != nil {
		return err
	}

	var ok bool
	*_p, ok = PlanetSupportUndefinedWithDefaultFromString(str)
	if !ok {
		return fmt.Errorf("Value %q does not represent a PlanetSupportUndefinedWithDefault", str)
	}
	return nil
}