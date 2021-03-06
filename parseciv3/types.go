package parseciv3

import (
	"encoding/binary"
	"io"
)

/*
The save file seems to be a simple binary dump of C++ data structures
packed with no byte padding. Generally speaking, most of the data is in
classes inherited from two basic classes. Both start with a 4-byte string
which appears to be a class name closely related to its function. One class
then has a 32-bit integer expressing the length in bytes of the data structure
following. The other has a 32-bit integer as a count of records. Each record
begins with a 32-bit length in bytes followed by the data. Before I knew this
I called each labeled length a "section", so I'll sometimes use that term
even now.

Some non-conformers appear to be the inital CIV3 section, but it's at least a
consistent length. The FLAV section is a list of lists. The second GAME section
in the SAV (which is the first GAME section of the non-BIC info) has an
apparently meaningless integer after the header followed by some predictable
data and then some as-yet unpredictable data which may be integer arrays, but
I haven't yet found the count. The length in bytes from GAME to DATE seems to
always be odd, so there must be a lone byte or a byte array in there somewhere.
I found a couple of stray apparent int32s after one of the DATE sections.

Later after the map tile data I have yet to figure out, too.

My strategy in the 2013-2015 Python version of this parser and my strategy so
far in Go is to parse the header, length/count and the data and then interpret
it. But several of the sections repeat with different lengths and data, especially
TILE but also WRLD and some others. I am presuming this is due to successive
versions of the game inheriting classes from the earlier game and adding to them,
and it shows up in the SAV file as the inheritance chain with data from each
generation. Mechanically parsing lenghts and counts works, but there really is
no advantage in meaning.

So I'm going to instead start making Go structs that will capture the entire
inheritance chain in one read which should make more sense programatically.

As I type this, I am mechanically parsing to WRLD but extracting little meaning
so far. During transition I'll be reading with two different methods.
*/

// ParsedData is the structure of the parsed data
type ParsedData map[string]Section

// Civ3Data contains the game data
type Civ3Data struct {
	FileName      string
	Compressed    bool
	Data          ParsedData
	Civ3          Civ3Header
	BicResources  BicResources
	BicFileHeader [4]byte
	VerNum        []VerNum
	Bldg          []Bldg
	Ctzn          []Ctzn
	Cult          []Cult
	Diff          []Difficulty
	Eras          []Era
	Espn          []Espn
	Expr          []Expr
	Flav          [][]Flavor
	Good          []Good
	//
	Wrld           Wrld
	Tile           []Tile
	Cont           []Continent
	ResourceCounts []int32
	Next           string
}

// Section is the inteface for the various structs decoded from the data files
type Section interface{}

// ListItem are the structs in a list
type ListItem interface{}

// Civ3Header is the SAV file header
// The Gobbledygook values appear to be 16 bytes of uncorrelated data, perhaps a hash or checksum?
type Civ3Header struct {
	Name                                                       [4]byte
	Always0x1a00                                               int16
	MaybeVersionMinor, MaybeVersionMajor                       int32
	Gobbeldygook1, Gobbeldygook2, Gobbeldygook3, Gobbeldygook4 uint32
}

// BicResources is part of the second SAV file section. Guessing at the alignment
type BicResources struct {
	Name         [4]byte
	Length       int32
	A            int32
	ResourcePath [0x100]byte
	B            int32
	BicPath      [0x100]byte
	C            int32
}

// ListHeader ...
type ListHeader struct {
	Name  [4]byte
	Count int32
}

// VerNum ...
type VerNum struct {
	Length          int32
	A, B            uint32
	BicMajorVersion int32
	BicMinorVersion int32
	Description     [640]byte
	Title           [64]byte
}

// Bldg ...
type Bldg struct {
	Length                     int32
	Description                [64]byte
	Name                       [32]byte
	CivilopediaEntry           [32]byte
	DoublesHappinessOf4        int32
	GainInEveryCity            int32
	GainInEveryCityOnContinent int32
	RequiredBuilding           int32
	Cost                       int32
	Culture                    int32
	BombardmentDefense         int32
	NavalBombardmentDefense    int32
	DefenseBonus               int32
	NavalDefenseBonus          int32
	MaintenanceCost            int32
	HappyFacesAllCities        int32
	HappyFaces                 int32
	UnhappyFacesAllCities      int32
	UnhappyFaces               int32
	NumberOfRequiredBuildings  int32
	AirPower                   int32
	NavalPower                 int32
	Pollution                  int32
	Production                 int32
	RequiredGovernment         int32
	SpaceshipPart              int32
	RequiredAdvance            int32
	RenderedObsoleteBy         int32
	RequiredResource1          int32
	RequiredResource2          int32
	ImprovementsBitMap         int32
	OtherCharacteristicsBitMap int32
	SmallWondersBitMap         int32
	WondersBitMap              int32
	NumberOfArmiesRequired     int32
	FlavorsBitMap              int32
	A                          int32
	UnitProducedPRTORef        int32
	UnitFrequency              int32
}

// Ctzn ...
type Ctzn struct {
	Length               int32
	DefaultCitizen       int32
	CitizensSingularName [32]byte
	CivilopediaEntry     [32]byte
	PluralName           [32]byte
	Prerequisite         int32
	Luxuries             int32
	Research             int32
	Taxes                int32
	Corruption           int32
	Construction         int32
}

// Cult ...
type Cult struct {
	Length                       int32
	CultureOpinionName           [64]byte
	ChanceOfSuccessfulPropaganda int32
	CultureRatioPercentage       int32
	CultureRatioDenominator      int32
	CultureRatioNumerator        int32
	InitialResistanceChance      int32
	ContinuedResistanceChance    int32
}

// Difficulty ...
type Difficulty struct {
	Length                                 int32
	DifficultyLevelName                    [64]byte
	NumberOfCitizensBornContent            int32
	MaxGovernmentTransitionTime            int32
	NumberOfDefensiveLandUnitsAIStartsWith int32
	NumberOfOffensiveLandUnitsAIStartsWith int32
	ExtraStartUnit1                        int32
	ExtraStartUnit2                        int32
	AdditionalFreeSupport                  int32
	BonusForEachCity                       int32
	AttackBonusAgainstBarbarians           int32
	CostFactor                             int32
	PercentageOfOptimalCities              int32
	AIToAITradeRate                        int32
	CorruptionPct                          int32
	MilitaryLaw                            int32
}

// Era ...
type Era struct {
	Length                      int32
	EraName                     [64]byte
	CivilopediaEntry            [32]byte
	Researcher1                 [32]byte
	Researcher2                 [32]byte
	Researcher3                 [32]byte
	Researcher4                 [32]byte
	Researcher5                 [32]byte
	NumberOfUsedResearcherNames int32
	A                           int32
}

// Espn ...
type Espn struct {
	Length                   int32
	Description              [128]byte
	MissionName              [64]byte
	CivilopediaEntry         [32]byte
	MissionPerformedByBitMap int32
	BaseCost                 int32
}

// Expr ...
type Expr struct {
	Length              int32
	ExperienceLevelName [32]byte
	BaseHitPoints       int32
	RetreatBonus        int32
}

// Good are the resource types
type Good struct {
	LengthOfResourceData88   int32
	NaturalResourceName      [24]byte
	CivilopediaEntry         [32]byte
	Type                     int32
	AppearanceRatio          int32
	DisappearanceProbability int32
	Icon                     int32
	Prerequisite             int32
	FoodBonus                int32
	ShieldsBonus             int32
	CommerceBonus            int32
}

// Base is one of the basic section structures of the game data
type Base struct {
	Name    [4]byte
	Length  int32
	RawData []byte
}

func newBase(r io.ReadSeeker) (Base, error) {
	var base Base
	var err error
	err = binary.Read(r, binary.LittleEndian, &base.Name)
	if err != nil {
		return base, ReadError{err}
	}
	err = binary.Read(r, binary.LittleEndian, &base.Length)
	if err != nil {
		return base, ReadError{err}
	}
	base.RawData = make([]byte, base.Length)
	err = binary.Read(r, binary.LittleEndian, &base.RawData)
	if err != nil {
		return base, ReadError{err}
	}
	return base, nil
}

// List is one of the basic section structures of the game data
type List struct {
	Name  [4]byte
	Count int32
	List  [][]byte
}

func newList(r io.ReadSeeker) (List, error) {
	var list List
	var err error
	err = binary.Read(r, binary.LittleEndian, &list.Name)
	if err != nil {
		return list, ReadError{err}
	}
	err = binary.Read(r, binary.LittleEndian, &list.Count)
	if err != nil {
		return list, ReadError{err}
	}
	for i := int32(0); i < list.Count; i++ {
		var length int32
		err = binary.Read(r, binary.LittleEndian, &length)
		if err != nil {
			return list, ReadError{err}
		}

		temp := make([]byte, length)
		err = binary.Read(r, binary.LittleEndian, &temp)
		list.List = append(list.List, temp)

	}
	return list, nil
}

// Flav is one of the basic section structures of the game data
type Flav struct {
	Name  [4]byte
	Count int32
	List  [][]Flavor
}

func newFlav(r io.ReadSeeker) (Flav, error) {
	var flav Flav
	var err error
	err = binary.Read(r, binary.LittleEndian, &flav.Name)
	if err != nil {
		return flav, ReadError{err}
	}
	err = binary.Read(r, binary.LittleEndian, &flav.Count)
	if err != nil {
		return flav, ReadError{err}
	}
	for i := int32(0); i < flav.Count; i++ {
		var count int32
		err = binary.Read(r, binary.LittleEndian, &count)
		if err != nil {
			return flav, ReadError{err}
		}
		flavorGroups := make([]Flavor, count)
		flav.List = append(flav.List, flavorGroups)
		for j := int32(0); j < count; j++ {
			flav.List[i][j] = Flavor{}
			err = binary.Read(r, binary.LittleEndian, &flav.List[i][j])
			if err != nil {
				return flav, ReadError{err}
			}
		}
	}
	return flav, nil
}

// Flavor is the leaf element of FLAV
// Hard-coding FlavorRelations at 7. Hopefully that always works
type Flavor struct {
	A                  int32
	FlavorName         [0x100]byte
	NumFlavorRelations int32
	FlavorRelations    [7]int32
}

// Game is the first section after the BIC.
type Game struct {
	// First two fields count for "class base"
	Name                       [4]byte
	_                          int32
	_                          [3]int32
	RenderFlags                int32
	DifficultyLevel            int32
	_                          int32
	UnitsCount                 int32
	CitiesCount                int32
	_                          int32
	_                          int32
	GlobalWarmingLevel         int32
	_                          int32
	_                          int32
	_                          int32
	CurrentTurn                int32
	_                          int32
	Random                     int32
	_                          int32
	CivFlags2                  int32
	CivFlags1                  int32
	_                          int32
	_                          int32
	_                          int32
	_                          int32
	_                          int32
	_                          int32
	_                          [48]int32
	Value1                     int32
	_                          [72]int32
	GameLimitPoints            int32
	GameLimitTurns             int32
	_                          [50]int32
	_                          int32
	_                          int32
	GameLimitDestroyedCities   int32
	GameLimitCityCulture       int32
	GameLimitCivCulture        int32
	GameLimitPopulation        int32
	GameLimitTerritory         int32
	GameLimitWonders           int32
	GameLimitDestroyedWonders  int32
	GameLimitAdvances          int32
	GameLimitCapturedCities    int32
	GameLimitVictoryPointPrice int32
	GameLimitPrincessRansom    int32
	DefaultDate1               int32
}

// GameNext is what Antal1987's dumps suggest is next, but I don't think so
type GameNext struct {
	_                     [27]int32
	PLGI                  [10]int32
	Date2                 Date
	Date3                 Date
	GameAggression        int32
	_                     int32
	CityStatIntArray      int32
	ResearchedAdvances    int32
	Wonders               int32
	WonderFlags           int32
	ImprovementTypesData1 int32
	ImprovementTypesData2 int32
	UnitTypesData1        int32
	UnitTypesData2        int32
	_                     int32
	_                     int32
	DefaultGameSettings   DefaultGameSettings
}

// Date DATE section ... I don't think this is nearly right
type Date struct {
	Name         [4]byte
	Length       int32
	Text         [16]byte
	_            [12]int32
	BaseTimeUnit int32
	Month        int32
	Week         int32
	Year         int32
	_            int32
}

// DefaultGameSettings from Antal1987's dump
type DefaultGameSettings struct {
	TurnsLimit           int32
	PointsLimit          int32
	DestroyedCitiesLimit int32
	CityCultureLimit     int32
	CivCultureLimit      int32
	PopulationLimit      int32
	TerritoryLimit       int32
	WondersLimit         int32
	DestroyedUnitsLimit  int32
	AdvancesLimit        int32
	CapturedCitiesLimit  int32
	VictoryPointPrice    int32
	PrincessPrice        int32
	PrincessRansom       int32
}

// Wrld is the Conquests' 3 WRLD sections combined
type Wrld struct {
	Name                   [4]byte `json:"-"`
	Length                 int32   `json:"-"`
	NumContinents          int16
	Name2                  [4]byte `json:"-"`
	Length2                int32   `json:"-"`
	OceanContinentID       int32
	MapHeight              int32
	DistanceBetweencivs    int32
	MaybeCivCount          int32
	D                      int32
	E                      int32
	MapWidth               int32
	CivStartLocationTileID [32]int32
	WorldSeed              int32
	G                      int32
	Name3                  [4]byte `json:"-"`
	Length3                int32   `json:"-"`
	GenOptions             WorldFeatures
}

// WorldFeatures is the map generation settings
// Barbs: -1 is off, 0 sedendary...3 raging
type WorldFeatures struct {
	Climate            int32
	ClimateFinal       int32
	Barbarians         int32
	BarbariansFinal    int32
	Landmass           int32
	LandmassFinal      int32
	OceanCoverage      int32
	OceanCoverageFinal int32
	Temperature        int32
	TemperatureFinal   int32
	Age                int32
	AgeFinal           int32
	Size               int32
}

// Tile is Conquests' 4 TILE sections per world tile combined
type Tile struct {
	Name                        [4]byte `json:"-"`
	Length                      int32   `json:"-"`
	Rivers                      uint8
	MaybeTerritoryCivID         int8
	MaybeLandmarkTerrain        int16
	ResourceType                int32
	TileUnitID                  int32
	MaybeSquarePart             int16
	MaybeVictoryPoint           int16
	MaybePreConquestsTileInfo   int32
	BarbTribeID                 int16
	CityID                      int16
	MaybeColonyID               int16
	ContinentID                 int16
	A                           int32
	B                           int32
	Name2                       [4]byte `json:"-"`
	Length2                     int32   `json:"-"`
	Improvements                int32
	C                           int8
	Terrain                     uint8
	MaybeByCivBitMask           uint32
	TerrainFeatures             uint16
	Name3                       [4]byte `json:"-"`
	Length3                     int32   `json:"-"`
	D                           int32
	Name4                       [4]byte `json:"-"`
	Length4                     int32   `json:"-"`
	VisibleToCiv                uint32
	VisibleNowToCivUnits        uint32
	MaybeVisibleToColonies      uint32
	VisibleToCivCulture         uint32
	E                           int32
	CityID2                     int16
	TradeNetworkIDByCiv         [32]int16
	MaybeImprovementsKnownToCiv [32]uint8
	F                           uint16
	G                           int32
	H                           int32
}

// Continent is the CONT section
// Land is 1 for land, 0 for water
// Size is number of tiles
type Continent struct {
	Name   [4]byte
	Length int32
	Land   int32
	Size   int32
}

// Lead is the LEAD section in the game, not the BIC
type Lead struct {
	Name [4]byte
	A    int32
}
