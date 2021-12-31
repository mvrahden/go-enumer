package pills

type Pill int

const PillPlacebo Pill = 0
const PillAspirin Pill = 1
const PillIbuprofen Pill = 2
const PillAcetaminophen Pill = 3

type PillSigned8 int8

const (
	PillSigned8Placebo PillSigned8 = iota
	PillSigned8Aspirin
	PillSigned8Ibuprofen
	PillSigned8Acetaminophen
)

type PillSigned16 int16

const (
	PillSigned16Placebo PillSigned16 = iota
	PillSigned16Aspirin
	PillSigned16Ibuprofen
	PillSigned16Acetaminophen
)

type PillSigned32 int32

const (
	PillSigned32Placebo PillSigned32 = iota
	PillSigned32Aspirin
	PillSigned32Ibuprofen
	PillSigned32Acetaminophen
)

type PillSigned64 int64

const (
	PillSigned64Placebo PillSigned64 = iota
	PillSigned64Aspirin
	PillSigned64Ibuprofen
	PillSigned64Acetaminophen
)

type PillUnsigned uint

const (
	PillUnsignedPlacebo PillUnsigned = iota
	PillUnsignedAspirin
	PillUnsignedIbuprofen
	PillUnsignedAcetaminophen
)

type PillUnsigned8 uint8

const (
	PillUnsigned8Placebo PillUnsigned8 = iota
	PillUnsigned8Aspirin
	PillUnsigned8Ibuprofen
	PillUnsigned8Acetaminophen
)

type PillUnsigned16 uint16

const (
	PillUnsigned16Placebo PillUnsigned16 = iota
	PillUnsigned16Aspirin
	PillUnsigned16Ibuprofen
	PillUnsigned16Acetaminophen
)

type PillUnsigned32 uint32

const (
	PillUnsigned32Placebo PillUnsigned32 = iota
	PillUnsigned32Aspirin
	PillUnsigned32Ibuprofen
	PillUnsigned32Acetaminophen
)

type PillUnsigned64 uint64

const (
	PillUnsigned64Placebo PillUnsigned64 = iota
	PillUnsigned64Aspirin
	PillUnsigned64Ibuprofen
	PillUnsigned64Acetaminophen
)

type PillRowed int

const PillRowedPlacebo, PillRowedAspirin, PillRowedIbuprofen, PillRowedAcetaminophen PillRowed = iota, 1, 2, 3

type PillAliased Pill

const PillAliasedPlacebo, PillAliasedAspirin, PillAliasedIbuprofen, PillAliasedAcetaminophen PillAliased = iota, 1, 2, 3

type PillUndefined int

const (
	PillUndefinedPlacebo PillUndefined = iota + 1
	PillUndefinedAspirin
	PillUndefinedIbuprofen
	PillUndefinedAcetaminophen
)
