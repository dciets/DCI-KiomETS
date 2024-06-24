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
}

func NewTerrainSerialisation(terrainType TerrainType, ownerIndex int, numberOfSoldier uint) *TerrainSerialisation {
	return &TerrainSerialisation{TerrainType: terrainType, OwnerIndex: ownerIndex, NumberOfSoldier: numberOfSoldier}
}
