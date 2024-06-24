package game

import "server/game/serialisation"

type BarricadeTerrainState struct {
	T *Terrain
}

func (b *BarricadeTerrainState) tick() {

}

func (b *BarricadeTerrainState) getFightingFunctions() (fightFunction, fightFunction) {
	return barrierTerrainF1Function, barrierTerrainF2Function
}

func (b *BarricadeTerrainState) getTerrainType() serialisation.TerrainType {
	return serialisation.Barricade
}
