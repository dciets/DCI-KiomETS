package game

const TickToChangeBuilding uint32 = 20

type ConstructionTerrainState struct {
	T                *Terrain
	currentTick      uint32
	nextTerrainState TerrainState
}

func (c *ConstructionTerrainState) tick() {
	c.currentTick++
	if c.currentTick == TickToChangeBuilding {
		c.T.state = c.nextTerrainState
	}
}

func (c *ConstructionTerrainState) getFightingFunctions() (fightFunction, fightFunction) {
	return defaultTerrainF1Function, defaultTerrainF2Function
}
