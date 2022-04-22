package project

// UserRole represents a set of user specific roles.
// UserRoles can be deserialized from "undefined"/empty values.
//go:enum -transform=lower -support=undefined
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
//go:enum -from=enums/timezones.csv
type Timezone uint

const (
	TimezoneEuropeLondon Timezone = 379
)

// Currency represents a set of 3-Letter codes of Currencies.
// Each Currency comes with its additional data columns to augment the enum.
// Note: The CSV has a 5-column layout. The columns 3-5 are added values.
//go:enum -from=enums/currencies.csv
type Currency uint

// CountryCode represents a set of 3-Letter codes of Countries.
// Each CountryCode comes with its additional data columns to augment the enum.
// Note: The CSV has an 8-column layout. The columns 3-8 are added values.
//go:enum -from=enums/countries.csv
type CountryCode uint
