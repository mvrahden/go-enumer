// Code generated by "%s"; DO NOT EDIT.

package greetings

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

const (
	_GreetingString      = "WorldРоссия中國日本한국ČeskáRepublika𝜋"
	_GreetingLowerString = "worldроссия中國日本한국českárepublika𝜋"
)

var (
	_GreetingValueRange = [2]Greeting{0, 6}
	_GreetingValues     = []Greeting{0, 1, 2, 3, 4, 5, 6}
	_GreetingStrings    = []string{_GreetingString[0:5], _GreetingString[5:17], _GreetingString[17:23], _GreetingString[23:29], _GreetingString[29:35], _GreetingString[35:51], _GreetingString[51:55]}
)

// _GreetingNoOp is a compile time assertion.
// An "invalid argument/out of bounds" compiler error signifies that the enum values have changed.
// Re-run the enumer command to generate an updated version of Greeting.
func _GreetingNoOp() {
	var x [1]struct{}
	_ = x[GreetingWorld-(0)]
	_ = x[GreetingРоссия-(1)]
	_ = x[Greeting中國-(2)]
	_ = x[Greeting日本-(3)]
	_ = x[Greeting한국-(4)]
	_ = x[GreetingČeskáRepublika-(5)]
	_ = x[Greeting𝜋-(6)]
}

// GreetingValues returns all values of the enum.
func GreetingValues() []Greeting {
	strs := make([]Greeting, len(_GreetingValues))
	copy(strs, _GreetingValues)
	return _GreetingValues
}

// GreetingStrings returns a slice of all String values of the enum.
func GreetingStrings() []string {
	strs := make([]string, len(_GreetingStrings))
	copy(strs, _GreetingStrings)
	return strs
}

// IsValid inspects whether the value is valid enum value.
func (_g Greeting) IsValid() bool {
	return _g >= _GreetingValueRange[0] && _g <= _GreetingValueRange[1]
}

// String returns the string of the enum value.
// If the enum value is invalid, it will produce a string
// of the following pattern Greeting(%d) instead.
func (_g Greeting) String() string {
	if !_g.IsValid() {
		return fmt.Sprintf("Greeting(%d)", _g)
	}
	idx := int(_g)
	return _GreetingStrings[idx]
}

var (
	_GreetingStringToValueMap = map[string]Greeting{
		_GreetingString[0:5]:   GreetingWorld,
		_GreetingString[5:17]:  GreetingРоссия,
		_GreetingString[17:23]: Greeting中國,
		_GreetingString[23:29]: Greeting日本,
		_GreetingString[29:35]: Greeting한국,
		_GreetingString[35:51]: GreetingČeskáRepublika,
		_GreetingString[51:55]: Greeting𝜋,
	}
	_GreetingLowerStringToValueMap = map[string]Greeting{
		_GreetingLowerString[0:5]:   GreetingWorld,
		_GreetingLowerString[5:17]:  GreetingРоссия,
		_GreetingLowerString[17:23]: Greeting中國,
		_GreetingLowerString[23:29]: Greeting日本,
		_GreetingLowerString[29:35]: Greeting한국,
		_GreetingLowerString[35:51]: GreetingČeskáRepublika,
		_GreetingLowerString[51:55]: Greeting𝜋,
	}
)

// GreetingFromString determines the enum value with an exact case match.
func GreetingFromString(raw string) (Greeting, bool) {
	if len(raw) == 0 {
		return Greeting(0), true
	}
	v, ok := _GreetingStringToValueMap[raw]
	if !ok {
		return Greeting(0), false
	}
	return v, true
}

// GreetingFromStringIgnoreCase determines the enum value with a case-insensitive match.
func GreetingFromStringIgnoreCase(raw string) (Greeting, bool) {
	v, ok := GreetingFromString(raw)
	if ok {
		return v, ok
	}
	v, ok = _GreetingLowerStringToValueMap[raw]
	if !ok {
		return Greeting(0), false
	}
	return v, true
}

// MarshalBinary implements the encoding.BinaryMarshaler interface for Greeting.
func (_g Greeting) MarshalBinary() ([]byte, error) {
	if !_g.IsValid() {
		return nil, fmt.Errorf("Cannot marshal invalid value %q as Greeting", _g)
	}
	return []byte(_g.String()), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface for Greeting.
func (_g *Greeting) UnmarshalBinary(text []byte) error {
	str := string(text)

	var ok bool
	*_g, ok = GreetingFromStringIgnoreCase(str)
	if !ok {
		return fmt.Errorf("Value %q does not represent a Greeting", str)
	}
	return nil
}

// MarshalGQL implements the graphql.Marshaler interface for Greeting.
func (_g Greeting) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(_g.String()))
}

// UnmarshalGQL implements the graphql.Unmarshaler interface for Greeting.
func (_g *Greeting) UnmarshalGQL(value interface{}) error {
	var str string
	switch v := value.(type) {
	case []byte:
		str = string(v)
	case string:
		str = v
	case fmt.Stringer:
		str = v.String()
	default:
		return fmt.Errorf("invalid value of Greeting: %[1]T(%[1]v)", value)
	}

	var ok bool
	*_g, ok = GreetingFromStringIgnoreCase(str)
	if !ok {
		return fmt.Errorf("Value %q does not represent a Greeting", str)
	}
	return nil
}

// MarshalJSON implements the json.Marshaler interface for Greeting.
func (_g Greeting) MarshalJSON() ([]byte, error) {
	if !_g.IsValid() {
		return nil, fmt.Errorf("Cannot marshal invalid value %q as Greeting", _g)
	}
	return json.Marshal(_g.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for Greeting.
func (_g *Greeting) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return fmt.Errorf("Greeting should be a string, got %q", data)
	}

	var ok bool
	*_g, ok = GreetingFromStringIgnoreCase(str)
	if !ok {
		return fmt.Errorf("Value %q does not represent a Greeting", str)
	}
	return nil
}

func (_g Greeting) Value() (driver.Value, error) {
	if !_g.IsValid() {
		return nil, fmt.Errorf("Cannot serialize invalid value %q as Greeting", _g)
	}
	return _g.String(), nil
}

func (_g *Greeting) Scan(value interface{}) error {
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
		return fmt.Errorf("invalid value of Greeting: %[1]T(%[1]v)", value)
	}

	var ok bool
	*_g, ok = GreetingFromStringIgnoreCase(str)
	if !ok {
		return fmt.Errorf("Value %q does not represent a Greeting", str)
	}
	return nil
}

// MarshalText implements the encoding.TextMarshaler interface for Greeting.
func (_g Greeting) MarshalText() ([]byte, error) {
	if !_g.IsValid() {
		return nil, fmt.Errorf("Cannot marshal invalid value %q as Greeting", _g)
	}
	return []byte(_g.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for Greeting.
func (_g *Greeting) UnmarshalText(text []byte) error {
	str := string(text)

	var ok bool
	*_g, ok = GreetingFromStringIgnoreCase(str)
	if !ok {
		return fmt.Errorf("Value %q does not represent a Greeting", str)
	}
	return nil
}

// MarshalYAML implements a YAML Marshaler for Greeting.
func (_g Greeting) MarshalYAML() (interface{}, error) {
	if !_g.IsValid() {
		return nil, fmt.Errorf("Cannot marshal invalid value %q as Greeting", _g)
	}
	return _g.String(), nil
}

// UnmarshalYAML implements a YAML Unmarshaler for Greeting.
func (_g *Greeting) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var str string
	if err := unmarshal(&str); err != nil {
		return err
	}

	var ok bool
	*_g, ok = GreetingFromStringIgnoreCase(str)
	if !ok {
		return fmt.Errorf("Value %q does not represent a Greeting", str)
	}
	return nil
}

// Values returns a slice of all String values of the enum.
func (Greeting) Values() []string {
	return GreetingStrings()
}
