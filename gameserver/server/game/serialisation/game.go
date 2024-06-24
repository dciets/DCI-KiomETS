package serialisation

type GameSerialisation struct {
	Players  []PlayerSerialisation  `json:"players"`
	Terrains []TerrainSerialisation `json:"terrains"`
	Pipes    []PipeSerialisation    `json:"pipes"`
}

func NewGameSerialisation(numberOfPlayer int, numberOfTerrain int, numberOfPipe int) *GameSerialisation {
	return &GameSerialisation{
		Players:  make([]PlayerSerialisation, numberOfPlayer),
		Terrains: make([]TerrainSerialisation, numberOfTerrain),
		Pipes:    make([]PipeSerialisation, numberOfPipe),
	}
}
