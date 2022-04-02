package project

// UserRole represents a set of user specific roles.
// UserRoles can be deserialized from "undefined"/empty values.
//go:enumer -transform=lower -support=undefined
type UserRole uint

const (
	UserRoleStandard UserRole = iota // note: as default
	UserRoleEditor
	UserRoleReviewer
	UserRoleAdmin
)

// Timezone represents a set of 424 Timezones from TimeZoneDB.
// Timezone has no default value, meaning it can only be deserialized from explicit values.
// Note: The CSV has a 2-column layout.
//go:enumer -from=enums/timezones.csv
type Timezone uint

const (
	TimezoneEuropeLondon Timezone = 379
)

// Currency represents a set of 3-Letter codes of Currencies.
// Each Currency comes with its canonical value.
// Note: The CSV has a 3-column layout. The 3rd column will add the Canonical Values.
//go:enumer -from=enums/currencies.csv
type Currency uint
