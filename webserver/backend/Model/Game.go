package Model

type Parameters struct {
	mapLength            int
	soldierSpeed         int
	soldierCreationSpeed int
	terrainChangeSpeed   int
	gameLength           int
}

func NewParameters(mapLength int, soldierSpeed int, soldierCreationSpeed int, terrainChangeSpeed int, gameLength int) *Parameters {
	return &Parameters{mapLength: mapLength, soldierSpeed: soldierSpeed, soldierCreationSpeed: soldierCreationSpeed, terrainChangeSpeed: terrainChangeSpeed, gameLength: gameLength}
}
func (p *Parameters) String() string {
	return "\tMap length : " + string(p.mapLength) + ",\n\t Soldier speed : " + string(p.soldierSpeed) + ",\n\t Soldier creation speed : " + string(p.soldierCreationSpeed) + ",\n\t Terrain change speed : " + string(p.terrainChangeSpeed) + ",\n\t Game length : " + string(p.gameLength)
}

type Game struct {
	Status     string
	Parameters Parameters
}

func NewGame(Status string, mapLength int, soldierSpeed int, soldierCreationSpeed int, terrainChangeSpeed int, gameLength int) *Game {
	return &Game{Status: Status, Parameters: Parameters{mapLength: mapLength, soldierSpeed: soldierSpeed, soldierCreationSpeed: soldierCreationSpeed, terrainChangeSpeed: terrainChangeSpeed, gameLength: gameLength}}
}

func (g *Game) String() string {
	return "Game status : " + g.Status + "\nParameters : \n" + g.Parameters.String()
}

func (g *Game) Start() {
	g.Status = "started"
	// TODO : Start the game
}

func (g *Game) Stop() {
	g.Status = "stopped"
	// TODO : Stop the game
}

func (g *Game) ChangeParameters(mapLength int, soldierSpeed int, soldierCreationSpeed int, terrainChangeSpeed int, gameLength int) {
	g.Parameters = Parameters{mapLength: mapLength, soldierSpeed: soldierSpeed, soldierCreationSpeed: soldierCreationSpeed, terrainChangeSpeed: terrainChangeSpeed, gameLength: gameLength}
}
func (g *Game) GetParameters() Parameters {
	return g.Parameters
}
