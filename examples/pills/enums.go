package pills

// Pill with a Doc string.
//go:enum
type PillUnsigned uint

const (
	PillUnsignedPlacebo PillUnsigned = iota
	PillUnsignedAspirin
	PillUnsignedIbuprofen
	PillUnsignedParacetamol
	PillUnsignedAcetaminophen PillUnsigned = iota - 1
	PillUnsignedVitaminC
)

//go:enum
type PillUnsigned8 uint8

const (
	PillUnsigned8Placebo       PillUnsigned8 = 0
	PillUnsigned8Aspirin       PillUnsigned8 = 1
	PillUnsigned8Ibuprofen     PillUnsigned8 = 2
	PillUnsigned8Paracetamol   PillUnsigned8 = 3
	PillUnsigned8Acetaminophen PillUnsigned8 = PillUnsigned8Paracetamol // hint: fallsback to Paracetamol
	PillUnsigned8VitaminC      PillUnsigned8 = 4
)

//go:enum
type PillUnsigned16 uint16

const (
	PillUnsigned16Placebo PillUnsigned16 = iota
	PillUnsigned16Aspirin
	PillUnsigned16Ibuprofen
	PillUnsigned16Paracetamol
	PillUnsigned16Acetaminophen PillUnsigned16 = iota - 1
	PillUnsigned16VitaminC
)

//go:enum
type PillUnsigned32 uint32

const (
	PillUnsigned32Placebo PillUnsigned32 = iota
	PillUnsigned32Aspirin
	PillUnsigned32Ibuprofen
	PillUnsigned32Paracetamol
	PillUnsigned32Acetaminophen PillUnsigned32 = iota - 1
	PillUnsigned32VitaminC
)

//go:enum
type PillUnsigned64 uint64

const (
	PillUnsigned64Placebo PillUnsigned64 = iota
	PillUnsigned64Aspirin
	PillUnsigned64Ibuprofen
	PillUnsigned64Paracetamol
	PillUnsigned64Acetaminophen PillUnsigned64 = iota - 1
	PillUnsigned64VitaminC
)

//go:enum
type PillRowed uint

const PillRowedPlacebo, PillRowedAspirin, PillRowedIbuprofen, PillRowedParacetamol, PillRowedAcetaminophen, PillRowedVitaminC PillRowed = 0, 1, 2, 3, 3, 4

//go:enum
type PillAliased PillUnsigned

const (
	PillAliasedPlacebo PillAliased = iota
	PillAliasedAspirin
	PillAliasedIbuprofen
	PillAliasedParacetamol
	PillAliasedAcetaminophen PillAliased = iota - 1
	PillAliasedVitaminC
)
