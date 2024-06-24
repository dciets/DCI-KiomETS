package game

import "server/game/serialisation"

const TickToCreate uint8 = 2

type FactoryTerrainState struct {
	T           *Terrain
	currentTick uint8
}

func (f *FactoryTerrainState) tick() {
	f.currentTick++
	if f.currentTick == TickToCreate {
		f.T.soldiers.count++
	}
}

func (f *FactoryTerrainState) getFightingFunctions() (fightFunction, fightFunction) {
	return defaultTerrainF1Function, defaultTerrainF2Function
}

func (f *FactoryTerrainState) getTerrainType() serialisation.TerrainType {
	return serialisation.Factory
}
