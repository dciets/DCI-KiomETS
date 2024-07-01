package serialisation

type TerrainType uint8

const (
	Barricade TerrainType = iota
	Factory
	Empty
	Construction
)

type TerrainSerialisation struct {
	TerrainType     TerrainType `json:"terrainType"`
	OwnerIndex      int         `json:"ownerIndex"`
	NumberOfSoldier uint        `json:"numberOfSoldier"`
	Position        [2]float32  `json:"position"`
}

func NewTerrainSerialisation(terrainType TerrainType, ownerIndex int, numberOfSoldier uint, position [2]float32) *TerrainSerialisation {
	return &TerrainSerialisation{TerrainType: terrainType, OwnerIndex: ownerIndex, NumberOfSoldier: numberOfSoldier, Position: position}
}
