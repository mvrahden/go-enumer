// Code generated by "%s"; DO NOT EDIT.

package booking

import (
	"encoding/json"
	"errors"
	"fmt"
)

var (
	ErrNoValidEnum = errors.New("not a valid enum")
)

const (
	_BookingStateString      = "CreatedUnavailableFailedCanceledNotFoundDeleted"
	_BookingStateLowerString = "createdunavailablefailedcancelednotfounddeleted"
)

var (
	_BookingStateValueRange     = [2]BookingState{0, 5}
	_BookingStateValues         = [6]BookingState{0, 1, 2, 3, 4, 5}
	_BookingStateStrings        = [6]string{_BookingStateString[0:7], _BookingStateString[7:18], _BookingStateString[18:24], _BookingStateString[24:32], _BookingStateString[32:40], _BookingStateString[40:47]}
	_BookingStateAdditionalData = [6]struct{ Description string }{
		{"The booking was created successfully"},
		{"The booking was not available"},
		{"The booking failed"},
		{"The booking was canceled"},
		{"The booking was not found"},
		{"The booking was deleted"},
	}
)

// BookingStateValues returns all values of the enum.
func BookingStateValues() []BookingState {
	cp := _BookingStateValues
	return cp[:]
}

// BookingStateStrings returns a slice of all String values of the enum.
func BookingStateStrings() []string {
	cp := _BookingStateStrings
	return cp[:]
}

// IsValid inspects whether the value is valid enum value.
func (_b BookingState) IsValid() bool {
	return _b >= _BookingStateValueRange[0] && _b <= _BookingStateValueRange[1]
}

// Validate whether the value is within the range of enum values.
func (_b BookingState) Validate() error {
	if !_b.IsValid() {
		return fmt.Errorf("BookingState(%d) is %w", _b, ErrNoValidEnum)
	}
	return nil
}

// String returns the string of the enum value.
// If the enum value is invalid, it will produce a string
// of the following pattern BookingState(%d) instead.
func (_b BookingState) String() string {
	if !_b.IsValid() {
		return fmt.Sprintf("BookingState(%d)", _b)
	}
	idx := uint(_b)
	return _BookingStateStrings[idx]
}

// GetDescription returns the "description" of the enum value as string
// if the enum is valid.
func (_b BookingState) GetDescription() (string, bool) {
	if !_b.IsValid() {
		return "", false
	}
	idx := uint(_b)
	d := _BookingStateAdditionalData[idx]
	return d.Description, true
}

var (
	_BookingStateStringToValueMap = map[string]BookingState{
		_BookingStateString[0:7]:   0,
		_BookingStateString[7:18]:  1,
		_BookingStateString[18:24]: 2,
		_BookingStateString[24:32]: 3,
		_BookingStateString[32:40]: 4,
		_BookingStateString[40:47]: 5,
	}
	_BookingStateLowerStringToValueMap = map[string]BookingState{
		_BookingStateLowerString[0:7]:   0,
		_BookingStateLowerString[7:18]:  1,
		_BookingStateLowerString[18:24]: 2,
		_BookingStateLowerString[24:32]: 3,
		_BookingStateLowerString[32:40]: 4,
		_BookingStateLowerString[40:47]: 5,
	}
)

// BookingStateFromString determines the enum value with an exact case match.
func BookingStateFromString(raw string) (BookingState, bool) {
	v, ok := _BookingStateStringToValueMap[raw]
	if !ok {
		return BookingState(0), false
	}
	return v, true
}

// BookingStateFromStringIgnoreCase determines the enum value with a case-insensitive match.
func BookingStateFromStringIgnoreCase(raw string) (BookingState, bool) {
	v, ok := BookingStateFromString(raw)
	if ok {
		return v, ok
	}
	v, ok = _BookingStateLowerStringToValueMap[raw]
	if !ok {
		return BookingState(0), false
	}
	return v, true
}

// MarshalYAML implements a YAML Marshaler for BookingState.
func (_b BookingState) MarshalYAML() (interface{}, error) {
	if err := _b.Validate(); err != nil {
		return nil, fmt.Errorf("Cannot marshal value %q as BookingState. %w", _b, err)
	}
	return _b.String(), nil
}

// UnmarshalYAML implements a YAML Unmarshaler for BookingState.
func (_b *BookingState) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var str string
	if err := unmarshal(&str); err != nil {
		return err
	}
	if len(str) == 0 {
		return fmt.Errorf("BookingState cannot be derived from empty string")
	}

	var ok bool
	*_b, ok = BookingStateFromStringIgnoreCase(str)
	if !ok {
		return fmt.Errorf("Value %q does not represent a BookingState", str)
	}
	return nil
}

// Values returns a slice of all String values of the enum.
func (BookingState) Values() []string {
	return BookingStateStrings()
}

const (
	_BookingStateWithConfigString      = "CreatedUnavailableFailedCanceledNotFoundDeleted"
	_BookingStateWithConfigLowerString = "createdunavailablefailedcancelednotfounddeleted"
)

var (
	_BookingStateWithConfigValueRange     = [2]BookingStateWithConfig{0, 5}
	_BookingStateWithConfigValues         = [6]BookingStateWithConfig{0, 1, 2, 3, 4, 5}
	_BookingStateWithConfigStrings        = [6]string{_BookingStateWithConfigString[0:7], _BookingStateWithConfigString[7:18], _BookingStateWithConfigString[18:24], _BookingStateWithConfigString[24:32], _BookingStateWithConfigString[32:40], _BookingStateWithConfigString[40:47]}
	_BookingStateWithConfigAdditionalData = [6]struct{ Description string }{
		{"The booking was created successfully"},
		{"The booking was not available"},
		{"The booking failed"},
		{"The booking was canceled"},
		{"The booking was not found"},
		{"The booking was deleted"},
	}
)

// BookingStateWithConfigValues returns all values of the enum.
func BookingStateWithConfigValues() []BookingStateWithConfig {
	cp := _BookingStateWithConfigValues
	return cp[:]
}

// BookingStateWithConfigStrings returns a slice of all String values of the enum.
func BookingStateWithConfigStrings() []string {
	cp := _BookingStateWithConfigStrings
	return cp[:]
}

// IsValid inspects whether the value is valid enum value.
func (_b BookingStateWithConfig) IsValid() bool {
	return _b >= _BookingStateWithConfigValueRange[0] && _b <= _BookingStateWithConfigValueRange[1]
}

// Validate whether the value is within the range of enum values.
func (_b BookingStateWithConfig) Validate() error {
	if !_b.IsValid() {
		return fmt.Errorf("BookingStateWithConfig(%d) is %w", _b, ErrNoValidEnum)
	}
	return nil
}

// String returns the string of the enum value.
// If the enum value is invalid, it will produce a string
// of the following pattern BookingStateWithConfig(%d) instead.
func (_b BookingStateWithConfig) String() string {
	if !_b.IsValid() {
		return fmt.Sprintf("BookingStateWithConfig(%d)", _b)
	}
	idx := uint(_b)
	return _BookingStateWithConfigStrings[idx]
}

// GetDescription returns the "description" of the enum value as string
// if the enum is valid.
func (_b BookingStateWithConfig) GetDescription() (string, bool) {
	if !_b.IsValid() {
		return "", false
	}
	idx := uint(_b)
	d := _BookingStateWithConfigAdditionalData[idx]
	return d.Description, true
}

var (
	_BookingStateWithConfigStringToValueMap = map[string]BookingStateWithConfig{
		_BookingStateWithConfigString[0:7]:   0,
		_BookingStateWithConfigString[7:18]:  1,
		_BookingStateWithConfigString[18:24]: 2,
		_BookingStateWithConfigString[24:32]: 3,
		_BookingStateWithConfigString[32:40]: 4,
		_BookingStateWithConfigString[40:47]: 5,
	}
	_BookingStateWithConfigLowerStringToValueMap = map[string]BookingStateWithConfig{
		_BookingStateWithConfigLowerString[0:7]:   0,
		_BookingStateWithConfigLowerString[7:18]:  1,
		_BookingStateWithConfigLowerString[18:24]: 2,
		_BookingStateWithConfigLowerString[24:32]: 3,
		_BookingStateWithConfigLowerString[32:40]: 4,
		_BookingStateWithConfigLowerString[40:47]: 5,
	}
)

// BookingStateWithConfigFromString determines the enum value with an exact case match.
func BookingStateWithConfigFromString(raw string) (BookingStateWithConfig, bool) {
	if len(raw) == 0 {
		return BookingStateWithConfig(0), true
	}
	v, ok := _BookingStateWithConfigStringToValueMap[raw]
	if !ok {
		return BookingStateWithConfig(0), false
	}
	return v, true
}

// BookingStateWithConfigFromStringIgnoreCase determines the enum value with a case-insensitive match.
func BookingStateWithConfigFromStringIgnoreCase(raw string) (BookingStateWithConfig, bool) {
	v, ok := BookingStateWithConfigFromString(raw)
	if ok {
		return v, ok
	}
	v, ok = _BookingStateWithConfigLowerStringToValueMap[raw]
	if !ok {
		return BookingStateWithConfig(0), false
	}
	return v, true
}

// MarshalJSON implements the json.Marshaler interface for BookingStateWithConfig.
func (_b BookingStateWithConfig) MarshalJSON() ([]byte, error) {
	if err := _b.Validate(); err != nil {
		return nil, fmt.Errorf("Cannot marshal value %q as BookingStateWithConfig. %w", _b, err)
	}
	return json.Marshal(_b.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for BookingStateWithConfig.
func (_b *BookingStateWithConfig) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return fmt.Errorf("BookingStateWithConfig should be a string, got %q", data)
	}

	var ok bool
	*_b, ok = BookingStateWithConfigFromString(str)
	if !ok {
		return fmt.Errorf("Value %q does not represent a BookingStateWithConfig", str)
	}
	return nil
}

// MarshalYAML implements a YAML Marshaler for BookingStateWithConfig.
func (_b BookingStateWithConfig) MarshalYAML() (interface{}, error) {
	if err := _b.Validate(); err != nil {
		return nil, fmt.Errorf("Cannot marshal value %q as BookingStateWithConfig. %w", _b, err)
	}
	return _b.String(), nil
}

// UnmarshalYAML implements a YAML Unmarshaler for BookingStateWithConfig.
func (_b *BookingStateWithConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var str string
	if err := unmarshal(&str); err != nil {
		return err
	}

	var ok bool
	*_b, ok = BookingStateWithConfigFromString(str)
	if !ok {
		return fmt.Errorf("Value %q does not represent a BookingStateWithConfig", str)
	}
	return nil
}

const (
	_BookingStateWithConstantsString      = "CreatedUnavailableFailedCanceledNotFoundDeleted"
	_BookingStateWithConstantsLowerString = "createdunavailablefailedcancelednotfounddeleted"
)

var (
	_BookingStateWithConstantsValueRange     = [2]BookingStateWithConstants{0, 5}
	_BookingStateWithConstantsValues         = [6]BookingStateWithConstants{0, 1, 2, 3, 4, 5}
	_BookingStateWithConstantsStrings        = [6]string{_BookingStateWithConstantsString[0:7], _BookingStateWithConstantsString[7:18], _BookingStateWithConstantsString[18:24], _BookingStateWithConstantsString[24:32], _BookingStateWithConstantsString[32:40], _BookingStateWithConstantsString[40:47]}
	_BookingStateWithConstantsAdditionalData = [6]struct{ Description string }{
		{"The booking was created successfully"},
		{"The booking was not available"},
		{"The booking failed"},
		{"The booking was canceled"},
		{"The booking was not found"},
		{"The booking was deleted"},
	}
)

// BookingStateWithConstantsValues returns all values of the enum.
func BookingStateWithConstantsValues() []BookingStateWithConstants {
	cp := _BookingStateWithConstantsValues
	return cp[:]
}

// BookingStateWithConstantsStrings returns a slice of all String values of the enum.
func BookingStateWithConstantsStrings() []string {
	cp := _BookingStateWithConstantsStrings
	return cp[:]
}

// IsValid inspects whether the value is valid enum value.
func (_b BookingStateWithConstants) IsValid() bool {
	return _b >= _BookingStateWithConstantsValueRange[0] && _b <= _BookingStateWithConstantsValueRange[1]
}

// Validate whether the value is within the range of enum values.
func (_b BookingStateWithConstants) Validate() error {
	if !_b.IsValid() {
		return fmt.Errorf("BookingStateWithConstants(%d) is %w", _b, ErrNoValidEnum)
	}
	return nil
}

// String returns the string of the enum value.
// If the enum value is invalid, it will produce a string
// of the following pattern BookingStateWithConstants(%d) instead.
func (_b BookingStateWithConstants) String() string {
	if !_b.IsValid() {
		return fmt.Sprintf("BookingStateWithConstants(%d)", _b)
	}
	idx := uint(_b)
	return _BookingStateWithConstantsStrings[idx]
}

// GetDescription returns the "description" of the enum value as string
// if the enum is valid.
func (_b BookingStateWithConstants) GetDescription() (string, bool) {
	if !_b.IsValid() {
		return "", false
	}
	idx := uint(_b)
	d := _BookingStateWithConstantsAdditionalData[idx]
	return d.Description, true
}

var (
	_BookingStateWithConstantsStringToValueMap = map[string]BookingStateWithConstants{
		_BookingStateWithConstantsString[0:7]:   0,
		_BookingStateWithConstantsString[7:18]:  1,
		_BookingStateWithConstantsString[18:24]: 2,
		_BookingStateWithConstantsString[24:32]: 3,
		_BookingStateWithConstantsString[32:40]: 4,
		_BookingStateWithConstantsString[40:47]: 5,
	}
	_BookingStateWithConstantsLowerStringToValueMap = map[string]BookingStateWithConstants{
		_BookingStateWithConstantsLowerString[0:7]:   0,
		_BookingStateWithConstantsLowerString[7:18]:  1,
		_BookingStateWithConstantsLowerString[18:24]: 2,
		_BookingStateWithConstantsLowerString[24:32]: 3,
		_BookingStateWithConstantsLowerString[32:40]: 4,
		_BookingStateWithConstantsLowerString[40:47]: 5,
	}
)

// BookingStateWithConstantsFromString determines the enum value with an exact case match.
func BookingStateWithConstantsFromString(raw string) (BookingStateWithConstants, bool) {
	v, ok := _BookingStateWithConstantsStringToValueMap[raw]
	if !ok {
		return BookingStateWithConstants(0), false
	}
	return v, true
}

// BookingStateWithConstantsFromStringIgnoreCase determines the enum value with a case-insensitive match.
func BookingStateWithConstantsFromStringIgnoreCase(raw string) (BookingStateWithConstants, bool) {
	v, ok := BookingStateWithConstantsFromString(raw)
	if ok {
		return v, ok
	}
	v, ok = _BookingStateWithConstantsLowerStringToValueMap[raw]
	if !ok {
		return BookingStateWithConstants(0), false
	}
	return v, true
}

// MarshalYAML implements a YAML Marshaler for BookingStateWithConstants.
func (_b BookingStateWithConstants) MarshalYAML() (interface{}, error) {
	if err := _b.Validate(); err != nil {
		return nil, fmt.Errorf("Cannot marshal value %q as BookingStateWithConstants. %w", _b, err)
	}
	return _b.String(), nil
}

// UnmarshalYAML implements a YAML Unmarshaler for BookingStateWithConstants.
func (_b *BookingStateWithConstants) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var str string
	if err := unmarshal(&str); err != nil {
		return err
	}
	if len(str) == 0 {
		return fmt.Errorf("BookingStateWithConstants cannot be derived from empty string")
	}

	var ok bool
	*_b, ok = BookingStateWithConstantsFromStringIgnoreCase(str)
	if !ok {
		return fmt.Errorf("Value %q does not represent a BookingStateWithConstants", str)
	}
	return nil
}

// Values returns a slice of all String values of the enum.
func (BookingStateWithConstants) Values() []string {
	return BookingStateWithConstantsStrings()
}
