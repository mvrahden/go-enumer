package planets

// Planet has NO default value here.
//go:enumer
type Planet uint8

const (
	PlanetMars Planet = iota + 1
	PlanetPluto
	PlanetVenus
	PlanetMercury
	PlanetJupiter
	PlanetSaturn
	PlanetUranus
	PlanetNeptune
)

// PlanetWithDefault has a default value Earth.
//go:enumer
type PlanetWithDefault uint8

const (
	PlanetWithDefaultEarth PlanetWithDefault = iota
	PlanetWithDefaultMars
	PlanetWithDefaultPluto
	PlanetWithDefaultVenus
	PlanetWithDefaultMercury
	PlanetWithDefaultJupiter
	PlanetWithDefaultSaturn
	PlanetWithDefaultUranus
	PlanetWithDefaultNeptune
)

// Planet has NO default value here.
// But it supports deserialization from "undefined"/zero values
// and serialization to an "" (empty string).
// For this scenario a special const will be generated "<type>Undefined"
//go:enumer -support=undefined
type PlanetSupportUndefined uint8

const (
	PlanetSupportUndefinedMars PlanetSupportUndefined = iota + 1
	PlanetSupportUndefinedPluto
	PlanetSupportUndefinedVenus
	PlanetSupportUndefinedMercury
	PlanetSupportUndefinedJupiter
	PlanetSupportUndefinedSaturn
	PlanetSupportUndefinedUranus
	PlanetSupportUndefinedNeptune
)

// PlanetSupportUndefinedWithDefault has a default value Earth
// and it supports deserialization from "undefined"/zero values.
//go:enumer -support=undefined
type PlanetSupportUndefinedWithDefault uint8

const (
	PlanetSupportUndefinedWithDefaultEarth PlanetSupportUndefinedWithDefault = iota
	PlanetSupportUndefinedWithDefaultMars
	PlanetSupportUndefinedWithDefaultPluto
	PlanetSupportUndefinedWithDefaultVenus
	PlanetSupportUndefinedWithDefaultMercury
	PlanetSupportUndefinedWithDefaultJupiter
	PlanetSupportUndefinedWithDefaultSaturn
	PlanetSupportUndefinedWithDefaultUranus
	PlanetSupportUndefinedWithDefaultNeptune
)
