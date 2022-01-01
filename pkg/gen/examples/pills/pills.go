package pills

type Pill int

const PillPlacebo Pill = 0
const PillAspirin Pill = 1
const PillIbuprofen Pill = 2
const PillParacetamol Pill = 3
const PillAcetaminophen Pill = 3
const PillVitaminC Pill = 4

type PillSigned8 int8

const (
	PillSigned8Placebo PillSigned8 = iota
	PillSigned8Aspirin
	PillSigned8Ibuprofen
	PillSigned8Paracetamol
	PillSigned8Acetaminophen PillSigned8 = iota - 1
	PillSigned8VitaminC
)

type PillSigned16 int16

const (
	PillSigned16Placebo PillSigned16 = iota
	PillSigned16Aspirin
	PillSigned16Ibuprofen
	PillSigned16Paracetamol
	PillSigned16Acetaminophen PillSigned16 = iota - 1
	PillSigned16VitaminC
)

type PillSigned32 int32

const (
	PillSigned32Placebo PillSigned32 = iota
	PillSigned32Aspirin
	PillSigned32Ibuprofen
	PillSigned32Paracetamol
	PillSigned32Acetaminophen PillSigned32 = iota - 1
	PillSigned32VitaminC
)

type PillSigned64 int64

const (
	PillSigned64Placebo PillSigned64 = iota
	PillSigned64Aspirin
	PillSigned64Ibuprofen
	PillSigned64Paracetamol
	PillSigned64Acetaminophen PillSigned64 = iota - 1
	PillSigned64VitaminC
)

type PillUnsigned uint

const (
	PillUnsignedPlacebo PillUnsigned = iota
	PillUnsignedAspirin
	PillUnsignedIbuprofen
	PillUnsignedParacetamol
	PillUnsignedAcetaminophen PillUnsigned = iota - 1
	PillUnsignedVitaminC
)

type PillUnsigned8 uint8

const (
	PillUnsigned8Placebo PillUnsigned8 = iota
	PillUnsigned8Aspirin
	PillUnsigned8Ibuprofen
	PillUnsigned8Paracetamol
	PillUnsigned8Acetaminophen PillUnsigned8 = iota - 1
	PillUnsigned8VitaminC
)

type PillUnsigned16 uint16

const (
	PillUnsigned16Placebo PillUnsigned16 = iota
	PillUnsigned16Aspirin
	PillUnsigned16Ibuprofen
	PillUnsigned16Paracetamol
	PillUnsigned16Acetaminophen PillUnsigned16 = iota - 1
	PillUnsigned16VitaminC
)

type PillUnsigned32 uint32

const (
	PillUnsigned32Placebo PillUnsigned32 = iota
	PillUnsigned32Aspirin
	PillUnsigned32Ibuprofen
	PillUnsigned32Paracetamol
	PillUnsigned32Acetaminophen PillUnsigned32 = iota - 1
	PillUnsigned32VitaminC
)

type PillUnsigned64 uint64

const (
	PillUnsigned64Placebo PillUnsigned64 = iota
	PillUnsigned64Aspirin
	PillUnsigned64Ibuprofen
	PillUnsigned64Paracetamol
	PillUnsigned64Acetaminophen PillUnsigned64 = iota - 1
	PillUnsigned64VitaminC
)

type PillRowed int

const PillRowedPlacebo, PillRowedAspirin, PillRowedIbuprofen, PillRowedParacetamol, PillRowedAcetaminophen, PillRowedVitaminC PillRowed = 0, 1, 2, 3, 3, 4

type PillAliased Pill

const PillAliasedPlacebo, PillAliasedAspirin, PillAliasedIbuprofen, PillAliasedParacetamol, PillAliasedAcetaminophen, PillAliasedVitaminC PillAliased = 0, 1, 2, 3, 3, 4

type PillUndefined int

const (
	PillUndefinedPlacebo PillUndefined = iota + 1
	PillUndefinedAspirin
	PillUndefinedIbuprofen
	PillUndefinedParacetamol
	PillUndefinedAcetaminophen PillUndefined = iota
	PillUndefinedVitaminC
)
