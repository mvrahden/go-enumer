// Code generated by "go-enumer (github.com/mvrahden/go-enumer)"; DO NOT EDIT.

package greetings

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"io"
	"strconv"
)

var (
	ErrNoValidEnum = errors.New("not a valid enum")
)

const (
	_GreetingString      = "Россия中國日本한국ČeskáRepublika𝜋"
	_GreetingLowerString = "россия中國日本한국českárepublika𝜋"
)

var (
	_GreetingValues  = [6]Greeting{1, 2, 3, 4, 5, 6}
	_GreetingStrings = [6]string{_GreetingString[0:12], _GreetingString[12:18], _GreetingString[18:24], _GreetingString[24:30], _GreetingString[30:46], _GreetingString[46:50]}
)

// _GreetingNoOp is a compile time assertion.
// An "invalid argument/out of bounds" compiler error signifies that the enum values have changed.
// Re-run the enumer command to generate an updated version of Greeting.
func _GreetingNoOp() {
	var x [1]struct{}
	_ = x[GreetingРоссия-(1)]
	_ = x[Greeting中國-(2)]
	_ = x[Greeting日本-(3)]
	_ = x[Greeting한국-(4)]
	_ = x[GreetingČeskáRepublika-(5)]
	_ = x[Greeting𝜋-(6)]
}

// GreetingValues returns all values of the enum.
func GreetingValues() []Greeting {
	cp := _GreetingValues
	return cp[:]
}

// GreetingStrings returns a slice of all String values of the enum.
func GreetingStrings() []string {
	cp := _GreetingStrings
	return cp[:]
}

// IsValid tests whether the value is a valid enum value.
func (_g Greeting) IsValid() bool {
	return _g >= 0 && _g <= 6
}

// Validate whether the value is within the range of enum values.
func (_g Greeting) Validate() error {
	if !_g.IsValid() {
		return fmt.Errorf("Greeting(%d) is %w", _g, ErrNoValidEnum)
	}
	return nil
}

// String returns the string of the enum value.
// If the enum value is invalid, it will produce a string
// of the following pattern Greeting(%d) instead.
func (_g Greeting) String() string {
	if !_g.IsValid() {
		return fmt.Sprintf("Greeting(%d)", _g)
	}
	if _g == 0 {
		return ""
	}
	idx := uint(_g) - 1
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
	if len(raw) == 0 {
		return Greeting(0), true
	}
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
	if err := _g.Validate(); err != nil {
		return nil, fmt.Errorf("Cannot marshal value %q as Greeting. %w", _g, err)
	}
	return []byte(_g.String()), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface for Greeting.
func (_g *Greeting) UnmarshalBinary(text []byte) error {
	str := string(text)

	var ok bool
	*_g, ok = GreetingFromString(str)
	if !ok {
		return fmt.Errorf("Value %q does not represent a Greeting", str)
	}
	return nil
}

// MarshalBSONValue implements the bson.ValueMarshaler interface for Greeting.
func (_g Greeting) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if err := _g.Validate(); err != nil {
		return 0, nil, fmt.Errorf("Cannot marshal value %q as Greeting. %w", _g, err)
	}
	if _g == 0 {
		return bsontype.Undefined, nil, nil
	}
	return bson.MarshalValue(_g.String())
}

// UnmarshalBSONValue implements the bson.ValueUnmarshaler interface for Greeting.
func (_g *Greeting) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t != bsontype.String && t != bsontype.Undefined {
		return fmt.Errorf("Greeting should be a string, got %q of Type %q", data, t)
	}
	str, data, ok := bsoncore.ReadString(data)
	if !ok {
		return fmt.Errorf("failed reading value as string, got %q", data)
	}

	*_g, ok = GreetingFromString(str)
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
	*_g, ok = GreetingFromString(str)
	if !ok {
		return fmt.Errorf("Value %q does not represent a Greeting", str)
	}
	return nil
}

// MarshalJSON implements the json.Marshaler interface for Greeting.
func (_g Greeting) MarshalJSON() ([]byte, error) {
	if err := _g.Validate(); err != nil {
		return nil, fmt.Errorf("Cannot marshal value %q as Greeting. %w", _g, err)
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
	*_g, ok = GreetingFromString(str)
	if !ok {
		return fmt.Errorf("Value %q does not represent a Greeting", str)
	}
	return nil
}

// Value implements the sql/driver.Valuer interface for Greeting.
func (_g Greeting) Value() (driver.Value, error) {
	if err := _g.Validate(); err != nil {
		return nil, fmt.Errorf("Cannot serialize value %q as Greeting. %w", _g, err)
	}
	if _g == 0 {
		return nil, nil
	}
	return _g.String(), nil
}

// Scan implements the sql/driver.Scanner interface for Greeting.
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
	*_g, ok = GreetingFromString(str)
	if !ok {
		return fmt.Errorf("Value %q does not represent a Greeting", str)
	}
	return nil
}

// MarshalText implements the encoding.TextMarshaler interface for Greeting.
func (_g Greeting) MarshalText() ([]byte, error) {
	if err := _g.Validate(); err != nil {
		return nil, fmt.Errorf("Cannot marshal value %q as Greeting. %w", _g, err)
	}
	return []byte(_g.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for Greeting.
func (_g *Greeting) UnmarshalText(text []byte) error {
	str := string(text)

	var ok bool
	*_g, ok = GreetingFromString(str)
	if !ok {
		return fmt.Errorf("Value %q does not represent a Greeting", str)
	}
	return nil
}

// MarshalYAML implements a YAML Marshaler for Greeting.
func (_g Greeting) MarshalYAML() (interface{}, error) {
	if err := _g.Validate(); err != nil {
		return nil, fmt.Errorf("Cannot marshal value %q as Greeting. %w", _g, err)
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
	*_g, ok = GreetingFromString(str)
	if !ok {
		return fmt.Errorf("Value %q does not represent a Greeting", str)
	}
	return nil
}

// Values returns a slice of all String values of the enum.
func (Greeting) Values() []string {
	return GreetingStrings()
}

const (
	_GreetingWithDefaultString      = "WorldРоссия中國日本한국ČeskáRepublika𝜋"
	_GreetingWithDefaultLowerString = "worldроссия中國日本한국českárepublika𝜋"
)

var (
	_GreetingWithDefaultValues  = [7]GreetingWithDefault{0, 1, 2, 3, 4, 5, 6}
	_GreetingWithDefaultStrings = [7]string{_GreetingWithDefaultString[0:5], _GreetingWithDefaultString[5:17], _GreetingWithDefaultString[17:23], _GreetingWithDefaultString[23:29], _GreetingWithDefaultString[29:35], _GreetingWithDefaultString[35:51], _GreetingWithDefaultString[51:55]}
)

// _GreetingWithDefaultNoOp is a compile time assertion.
// An "invalid argument/out of bounds" compiler error signifies that the enum values have changed.
// Re-run the enumer command to generate an updated version of GreetingWithDefault.
func _GreetingWithDefaultNoOp() {
	var x [1]struct{}
	_ = x[GreetingWithDefaultWorld-(0)]
	_ = x[GreetingWithDefaultРоссия-(1)]
	_ = x[GreetingWithDefault中國-(2)]
	_ = x[GreetingWithDefault日本-(3)]
	_ = x[GreetingWithDefault한국-(4)]
	_ = x[GreetingWithDefaultČeskáRepublika-(5)]
	_ = x[GreetingWithDefault𝜋-(6)]
}

// GreetingWithDefaultValues returns all values of the enum.
func GreetingWithDefaultValues() []GreetingWithDefault {
	cp := _GreetingWithDefaultValues
	return cp[:]
}

// GreetingWithDefaultStrings returns a slice of all String values of the enum.
func GreetingWithDefaultStrings() []string {
	cp := _GreetingWithDefaultStrings
	return cp[:]
}

// IsValid tests whether the value is a valid enum value.
func (_g GreetingWithDefault) IsValid() bool {
	return _g >= 0 && _g <= 6
}

// Validate whether the value is within the range of enum values.
func (_g GreetingWithDefault) Validate() error {
	if !_g.IsValid() {
		return fmt.Errorf("GreetingWithDefault(%d) is %w", _g, ErrNoValidEnum)
	}
	return nil
}

// String returns the string of the enum value.
// If the enum value is invalid, it will produce a string
// of the following pattern GreetingWithDefault(%d) instead.
func (_g GreetingWithDefault) String() string {
	if !_g.IsValid() {
		return fmt.Sprintf("GreetingWithDefault(%d)", _g)
	}
	idx := uint(_g)
	return _GreetingWithDefaultStrings[idx]
}

var (
	_GreetingWithDefaultStringToValueMap = map[string]GreetingWithDefault{
		_GreetingWithDefaultString[0:5]:   GreetingWithDefaultWorld,
		_GreetingWithDefaultString[5:17]:  GreetingWithDefaultРоссия,
		_GreetingWithDefaultString[17:23]: GreetingWithDefault中國,
		_GreetingWithDefaultString[23:29]: GreetingWithDefault日本,
		_GreetingWithDefaultString[29:35]: GreetingWithDefault한국,
		_GreetingWithDefaultString[35:51]: GreetingWithDefaultČeskáRepublika,
		_GreetingWithDefaultString[51:55]: GreetingWithDefault𝜋,
	}
	_GreetingWithDefaultLowerStringToValueMap = map[string]GreetingWithDefault{
		_GreetingWithDefaultLowerString[0:5]:   GreetingWithDefaultWorld,
		_GreetingWithDefaultLowerString[5:17]:  GreetingWithDefaultРоссия,
		_GreetingWithDefaultLowerString[17:23]: GreetingWithDefault中國,
		_GreetingWithDefaultLowerString[23:29]: GreetingWithDefault日本,
		_GreetingWithDefaultLowerString[29:35]: GreetingWithDefault한국,
		_GreetingWithDefaultLowerString[35:51]: GreetingWithDefaultČeskáRepublika,
		_GreetingWithDefaultLowerString[51:55]: GreetingWithDefault𝜋,
	}
)

// GreetingWithDefaultFromString determines the enum value with an exact case match.
func GreetingWithDefaultFromString(raw string) (GreetingWithDefault, bool) {
	if len(raw) == 0 {
		return GreetingWithDefault(0), true
	}
	v, ok := _GreetingWithDefaultStringToValueMap[raw]
	if !ok {
		return GreetingWithDefault(0), false
	}
	return v, true
}

// GreetingWithDefaultFromStringIgnoreCase determines the enum value with a case-insensitive match.
func GreetingWithDefaultFromStringIgnoreCase(raw string) (GreetingWithDefault, bool) {
	if len(raw) == 0 {
		return GreetingWithDefault(0), true
	}
	v, ok := GreetingWithDefaultFromString(raw)
	if ok {
		return v, ok
	}
	v, ok = _GreetingWithDefaultLowerStringToValueMap[raw]
	if !ok {
		return GreetingWithDefault(0), false
	}
	return v, true
}

// MarshalBinary implements the encoding.BinaryMarshaler interface for GreetingWithDefault.
func (_g GreetingWithDefault) MarshalBinary() ([]byte, error) {
	if err := _g.Validate(); err != nil {
		return nil, fmt.Errorf("Cannot marshal value %q as GreetingWithDefault. %w", _g, err)
	}
	return []byte(_g.String()), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface for GreetingWithDefault.
func (_g *GreetingWithDefault) UnmarshalBinary(text []byte) error {
	str := string(text)

	var ok bool
	*_g, ok = GreetingWithDefaultFromString(str)
	if !ok {
		return fmt.Errorf("Value %q does not represent a GreetingWithDefault", str)
	}
	return nil
}

// MarshalBSONValue implements the bson.ValueMarshaler interface for GreetingWithDefault.
func (_g GreetingWithDefault) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if err := _g.Validate(); err != nil {
		return 0, nil, fmt.Errorf("Cannot marshal value %q as GreetingWithDefault. %w", _g, err)
	}
	return bson.MarshalValue(_g.String())
}

// UnmarshalBSONValue implements the bson.ValueUnmarshaler interface for GreetingWithDefault.
func (_g *GreetingWithDefault) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t != bsontype.String && t != bsontype.Undefined {
		return fmt.Errorf("GreetingWithDefault should be a string, got %q of Type %q", data, t)
	}
	str, data, ok := bsoncore.ReadString(data)
	if !ok {
		return fmt.Errorf("failed reading value as string, got %q", data)
	}

	*_g, ok = GreetingWithDefaultFromString(str)
	if !ok {
		return fmt.Errorf("Value %q does not represent a GreetingWithDefault", str)
	}
	return nil
}

// MarshalGQL implements the graphql.Marshaler interface for GreetingWithDefault.
func (_g GreetingWithDefault) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(_g.String()))
}

// UnmarshalGQL implements the graphql.Unmarshaler interface for GreetingWithDefault.
func (_g *GreetingWithDefault) UnmarshalGQL(value interface{}) error {
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
		return fmt.Errorf("invalid value of GreetingWithDefault: %[1]T(%[1]v)", value)
	}

	var ok bool
	*_g, ok = GreetingWithDefaultFromString(str)
	if !ok {
		return fmt.Errorf("Value %q does not represent a GreetingWithDefault", str)
	}
	return nil
}

// MarshalJSON implements the json.Marshaler interface for GreetingWithDefault.
func (_g GreetingWithDefault) MarshalJSON() ([]byte, error) {
	if err := _g.Validate(); err != nil {
		return nil, fmt.Errorf("Cannot marshal value %q as GreetingWithDefault. %w", _g, err)
	}
	return json.Marshal(_g.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for GreetingWithDefault.
func (_g *GreetingWithDefault) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return fmt.Errorf("GreetingWithDefault should be a string, got %q", data)
	}

	var ok bool
	*_g, ok = GreetingWithDefaultFromString(str)
	if !ok {
		return fmt.Errorf("Value %q does not represent a GreetingWithDefault", str)
	}
	return nil
}

// Value implements the sql/driver.Valuer interface for GreetingWithDefault.
func (_g GreetingWithDefault) Value() (driver.Value, error) {
	if err := _g.Validate(); err != nil {
		return nil, fmt.Errorf("Cannot serialize value %q as GreetingWithDefault. %w", _g, err)
	}
	return _g.String(), nil
}

// Scan implements the sql/driver.Scanner interface for GreetingWithDefault.
func (_g *GreetingWithDefault) Scan(value interface{}) error {
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
		return fmt.Errorf("invalid value of GreetingWithDefault: %[1]T(%[1]v)", value)
	}

	var ok bool
	*_g, ok = GreetingWithDefaultFromString(str)
	if !ok {
		return fmt.Errorf("Value %q does not represent a GreetingWithDefault", str)
	}
	return nil
}

// MarshalText implements the encoding.TextMarshaler interface for GreetingWithDefault.
func (_g GreetingWithDefault) MarshalText() ([]byte, error) {
	if err := _g.Validate(); err != nil {
		return nil, fmt.Errorf("Cannot marshal value %q as GreetingWithDefault. %w", _g, err)
	}
	return []byte(_g.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for GreetingWithDefault.
func (_g *GreetingWithDefault) UnmarshalText(text []byte) error {
	str := string(text)

	var ok bool
	*_g, ok = GreetingWithDefaultFromString(str)
	if !ok {
		return fmt.Errorf("Value %q does not represent a GreetingWithDefault", str)
	}
	return nil
}

// MarshalYAML implements a YAML Marshaler for GreetingWithDefault.
func (_g GreetingWithDefault) MarshalYAML() (interface{}, error) {
	if err := _g.Validate(); err != nil {
		return nil, fmt.Errorf("Cannot marshal value %q as GreetingWithDefault. %w", _g, err)
	}
	return _g.String(), nil
}

// UnmarshalYAML implements a YAML Unmarshaler for GreetingWithDefault.
func (_g *GreetingWithDefault) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var str string
	if err := unmarshal(&str); err != nil {
		return err
	}

	var ok bool
	*_g, ok = GreetingWithDefaultFromString(str)
	if !ok {
		return fmt.Errorf("Value %q does not represent a GreetingWithDefault", str)
	}
	return nil
}

// Values returns a slice of all String values of the enum.
func (GreetingWithDefault) Values() []string {
	return GreetingWithDefaultStrings()
}
