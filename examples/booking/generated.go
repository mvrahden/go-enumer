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
	_BookingStateDataColumn0 = "The booking was created successfullyThe booking was not availableThe booking failedThe booking was canceledThe booking was not foundThe booking was deleted"
)

var (
	_BookingStateValueRange     = [2]BookingState{0, 5}
	_BookingStateValues         = []BookingState{0, 1, 2, 3, 4, 5}
	_BookingStateStrings        = []string{_BookingStateString[0:7], _BookingStateString[7:18], _BookingStateString[18:24], _BookingStateString[24:32], _BookingStateString[32:40], _BookingStateString[40:47]}
	_BookingStateAdditionalData = map[uint8]map[BookingState]string{
		0: {0: _BookingStateDataColumn0[0:36], 1: _BookingStateDataColumn0[36:65], 2: _BookingStateDataColumn0[65:83], 3: _BookingStateDataColumn0[83:107], 4: _BookingStateDataColumn0[107:132], 5: _BookingStateDataColumn0[132:155]},
	}
)

// BookingStateValues returns all values of the enum.
func BookingStateValues() []BookingState {
	cp := make([]BookingState, len(_BookingStateValues))
	copy(cp, _BookingStateValues)
	return cp
}

// BookingStateStrings returns a slice of all String values of the enum.
func BookingStateStrings() []string {
	cp := make([]string, len(_BookingStateStrings))
	copy(cp, _BookingStateStrings)
	return cp
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

// GetDescription returns the "description" string of the enum value.
// If the enum value is invalid, it will produce a string
// of the following pattern BookingState(%d).Description instead.
func (_b BookingState) GetDescription() (string, bool) {
	if !_b.IsValid() {
		return fmt.Sprintf("BookingState(%d).Description", _b), false
	}
	return _b._fromColumn(0)
}

// _fromColumn looks up additional data of the enum.
func (_b BookingState) _fromColumn(colId uint8) (cellValue string, ok bool) {
	if ok := _b.IsValid(); !ok {
		return "", false
	}
	col, ok := _BookingStateAdditionalData[colId]
	if !ok {
		return "", false
	}
	v, ok := col[_b]
	if !ok {
		return "", false
	}
	return v, ok
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
	_BookingStateWithConfigDataColumn0 = "The booking was created successfullyThe booking was not availableThe booking failedThe booking was canceledThe booking was not foundThe booking was deleted"
)

var (
	_BookingStateWithConfigValueRange     = [2]BookingStateWithConfig{0, 5}
	_BookingStateWithConfigValues         = []BookingStateWithConfig{0, 1, 2, 3, 4, 5}
	_BookingStateWithConfigStrings        = []string{_BookingStateWithConfigString[0:7], _BookingStateWithConfigString[7:18], _BookingStateWithConfigString[18:24], _BookingStateWithConfigString[24:32], _BookingStateWithConfigString[32:40], _BookingStateWithConfigString[40:47]}
	_BookingStateWithConfigAdditionalData = map[uint8]map[BookingStateWithConfig]string{
		0: {0: _BookingStateWithConfigDataColumn0[0:36], 1: _BookingStateWithConfigDataColumn0[36:65], 2: _BookingStateWithConfigDataColumn0[65:83], 3: _BookingStateWithConfigDataColumn0[83:107], 4: _BookingStateWithConfigDataColumn0[107:132], 5: _BookingStateWithConfigDataColumn0[132:155]},
	}
)

// BookingStateWithConfigValues returns all values of the enum.
func BookingStateWithConfigValues() []BookingStateWithConfig {
	cp := make([]BookingStateWithConfig, len(_BookingStateWithConfigValues))
	copy(cp, _BookingStateWithConfigValues)
	return cp
}

// BookingStateWithConfigStrings returns a slice of all String values of the enum.
func BookingStateWithConfigStrings() []string {
	cp := make([]string, len(_BookingStateWithConfigStrings))
	copy(cp, _BookingStateWithConfigStrings)
	return cp
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

// GetDescription returns the "description" string of the enum value.
// If the enum value is invalid, it will produce a string
// of the following pattern BookingStateWithConfig(%d).Description instead.
func (_b BookingStateWithConfig) GetDescription() (string, bool) {
	if !_b.IsValid() {
		return fmt.Sprintf("BookingStateWithConfig(%d).Description", _b), false
	}
	return _b._fromColumn(0)
}

// _fromColumn looks up additional data of the enum.
func (_b BookingStateWithConfig) _fromColumn(colId uint8) (cellValue string, ok bool) {
	if ok := _b.IsValid(); !ok {
		return "", false
	}
	col, ok := _BookingStateWithConfigAdditionalData[colId]
	if !ok {
		return "", false
	}
	v, ok := col[_b]
	if !ok {
		return "", false
	}
	return v, ok
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
	_BookingStateWithConstantsDataColumn0 = "The booking was created successfullyThe booking was not availableThe booking failedThe booking was canceledThe booking was not foundThe booking was deleted"
)

var (
	_BookingStateWithConstantsValueRange     = [2]BookingStateWithConstants{0, 5}
	_BookingStateWithConstantsValues         = []BookingStateWithConstants{0, 1, 2, 3, 4, 5}
	_BookingStateWithConstantsStrings        = []string{_BookingStateWithConstantsString[0:7], _BookingStateWithConstantsString[7:18], _BookingStateWithConstantsString[18:24], _BookingStateWithConstantsString[24:32], _BookingStateWithConstantsString[32:40], _BookingStateWithConstantsString[40:47]}
	_BookingStateWithConstantsAdditionalData = map[uint8]map[BookingStateWithConstants]string{
		0: {0: _BookingStateWithConstantsDataColumn0[0:36], 1: _BookingStateWithConstantsDataColumn0[36:65], 2: _BookingStateWithConstantsDataColumn0[65:83], 3: _BookingStateWithConstantsDataColumn0[83:107], 4: _BookingStateWithConstantsDataColumn0[107:132], 5: _BookingStateWithConstantsDataColumn0[132:155]},
	}
)

// BookingStateWithConstantsValues returns all values of the enum.
func BookingStateWithConstantsValues() []BookingStateWithConstants {
	cp := make([]BookingStateWithConstants, len(_BookingStateWithConstantsValues))
	copy(cp, _BookingStateWithConstantsValues)
	return cp
}

// BookingStateWithConstantsStrings returns a slice of all String values of the enum.
func BookingStateWithConstantsStrings() []string {
	cp := make([]string, len(_BookingStateWithConstantsStrings))
	copy(cp, _BookingStateWithConstantsStrings)
	return cp
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

// GetDescription returns the "description" string of the enum value.
// If the enum value is invalid, it will produce a string
// of the following pattern BookingStateWithConstants(%d).Description instead.
func (_b BookingStateWithConstants) GetDescription() (string, bool) {
	if !_b.IsValid() {
		return fmt.Sprintf("BookingStateWithConstants(%d).Description", _b), false
	}
	return _b._fromColumn(0)
}

// _fromColumn looks up additional data of the enum.
func (_b BookingStateWithConstants) _fromColumn(colId uint8) (cellValue string, ok bool) {
	if ok := _b.IsValid(); !ok {
		return "", false
	}
	col, ok := _BookingStateWithConstantsAdditionalData[colId]
	if !ok {
		return "", false
	}
	v, ok := col[_b]
	if !ok {
		return "", false
	}
	return v, ok
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
