package game

type EmptyTerrainState struct {
	T *Terrain
}

func (e *EmptyTerrainState) tick() {

}

func (e *EmptyTerrainState) getFightingFunctions() (fightFunction, fightFunction) {
	return defaultTerrainF1Function, defaultTerrainF2Function
}
