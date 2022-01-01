// Code generated by "%s"; DO NOT EDIT.

package greetings

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

const (
	_GreetingString      = "Россия中國日本한국ČeskáRepublika𝜋"
	_GreetingLowerString = "россия中國日本한국českárepublika𝜋"
)

var (
	_GreetingValueRange = [2]Greeting{1, 6}

	_GreetingValues = []Greeting{1, 2, 3, 4, 5, 6}

	_GreetingStrings = []string{_GreetingString[0:12], _GreetingString[12:18], _GreetingString[18:24], _GreetingString[24:30], _GreetingString[30:46], _GreetingString[46:50]}
)

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
// If the enum value is invalid.
func (_g Greeting) String() string {
	if !_g.IsValid() {
		return fmt.Sprintf("Greeting(%d)", _g)
	}
	idx := int(_g) - 1
	return _GreetingStrings[idx]
}

var (
	_GreetingStringToValueMap = map[string]Greeting{
		_GreetingString[0:12]:  GreetingРоссия,
		_GreetingString[12:18]: Greeting中國,
		_GreetingString[18:24]: Greeting日本,
		_GreetingString[24:30]: Greeting한국,
		_GreetingString[30:46]: GreetingČeskáRepublika,
		_GreetingString[46:50]: Greeting𝜋,
	}
	_GreetingLowerStringToValueMap = map[string]Greeting{
		_GreetingLowerString[0:12]:  GreetingРоссия,
		_GreetingLowerString[12:18]: Greeting中國,
		_GreetingLowerString[18:24]: Greeting日本,
		_GreetingLowerString[24:30]: Greeting한국,
		_GreetingLowerString[30:46]: GreetingČeskáRepublika,
		_GreetingLowerString[46:50]: Greeting𝜋,
	}
)

// GreetingFromString determines the enum value with an exact case match.
func GreetingFromString(raw string) (Greeting, bool) {
	v, ok := _GreetingStringToValueMap[raw]
	if !ok {
		return Greeting(0), false
	}
	return v, true
}

// GreetingFromStringIgnoreCase determines the enum value with an case-insensitive match.
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
		return nil, fmt.Errorf("Cannot marshal invalid value %q as Greeting", _g.String())
	}
	return []byte(_g.String()), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface for Greeting.
func (_g *Greeting) UnmarshalBinary(text []byte) error {
	str := string(text)
	if len(str) == 0 {
		return fmt.Errorf("Greeting cannot be derived from empty string")
	}

	var ok bool
	*_g, ok = GreetingFromString(str)
	if !ok {
		return fmt.Errorf("Value %q does not represent a Greeting", str)
	}
	return nil
}

// MarshalJSON implements the json.Marshaler interface for Greeting.
func (_g Greeting) MarshalJSON() ([]byte, error) {
	if !_g.IsValid() {
		return nil, fmt.Errorf("Cannot marshal invalid value %q as Greeting", _g.String())
	}
	return json.Marshal(_g.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for Greeting.
func (_g *Greeting) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return fmt.Errorf("Greeting should be a string, got %q", data)
	}
	if len(str) == 0 {
		return fmt.Errorf("Greeting cannot be derived from empty string")
	}

	var ok bool
	*_g, ok = GreetingFromString(str)
	if !ok {
		return fmt.Errorf("Value %q does not represent a Greeting", str)
	}
	return nil
}

func (_g Greeting) Value() (driver.Value, error) {
	if !_g.IsValid() {
		return nil, fmt.Errorf("Cannot serialize invalid value %q as Greeting", _g.String())
	}
	return _g.String(), nil
}

func (_g *Greeting) Scan(value interface{}) error {
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
		return fmt.Errorf("invalid value of Greeting: %[1]T(%[1]v)", value)
	}
	if len(str) == 0 {
		return fmt.Errorf("Greeting cannot be derived from empty string")
	}

	var ok bool
	*_g, ok = GreetingFromString(str)
	if !ok {
		return fmt.Errorf("Value %q does not represent a Greeting", str)
	}
	return nil
}

// MarshalText implements the encoding.TextMarshaler interface for Greeting.
func (_g Greeting) MarshalText() ([]byte, error) {
	if !_g.IsValid() {
		return nil, fmt.Errorf("Cannot marshal invalid value %q as Greeting", _g.String())
	}
	return []byte(_g.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for Greeting.
func (_g *Greeting) UnmarshalText(text []byte) error {
	str := string(text)
	if len(str) == 0 {
		return fmt.Errorf("Greeting cannot be derived from empty string")
	}

	var ok bool
	*_g, ok = GreetingFromString(str)
	if !ok {
		return fmt.Errorf("Value %q does not represent a Greeting", str)
	}
	return nil
}

// MarshalYAML implements a YAML Marshaler for Greeting
func (_g Greeting) MarshalYAML() (interface{}, error) {
	if !_g.IsValid() {
		return nil, fmt.Errorf("Cannot marshal invalid value %q as Greeting", _g.String())
	}
	return _g.String(), nil
}

// UnmarshalYAML implements a YAML Unmarshaler for Greeting
func (_g *Greeting) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var str string
	if err := unmarshal(&str); err != nil {
		return err
	}
	if len(str) == 0 {
		return fmt.Errorf("Greeting cannot be derived from empty string")
	}

	var ok bool
	*_g, ok = GreetingFromString(str)
	if !ok {
		return fmt.Errorf("Value %q does not represent a Greeting", str)
	}
	return nil
}
