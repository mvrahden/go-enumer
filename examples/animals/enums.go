package animals

//go:enum -transform=noop
type Animal uint8

const (
	AnimalDog Animal = iota
	AnimalCat
	AnimalSeal
	AnimalSeaLion
	AnimalIceBear
)

//go:enum -transform=upper-kebab
type Mammal uint8

const (
	MammalBumblebeeBat Mammal = iota
	MammalBlueWhale
	MammalBowheadWhale
)

//go:enum -transform=upper-snake
type Bird uint8

const (
	BirdAlbatross Bird = iota
	BirdHummingBird
	BirdDarwinsFinch
	BirdOstrich
	BirdKingFisher
)

//go:enum -transform=camel
type Reptile uint8

const (
	ReptileSaltwaterCrocodile Reptile = iota
	ReptileBeardedDragon
	ReptileChameleon
	ReptileComodoDragon
)

//go:enum -transform=snake
type Fish uint8

const (
	FishGiantGrouper Fish = iota
	FishHagfish
	FishReedfish
	FishBowfin
	FishCatfish
	FishHornShark
)
