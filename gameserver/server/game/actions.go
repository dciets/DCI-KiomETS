package game

import "server/game/serialisation"

type ActionType int

const (
	Move ActionType = iota
	Build
)

type MoveAction struct {
	FromId   string `json:"fromId"`
	ToId     string `json:"toId"`
	Quantity uint   `json:"quantity"`
}

type BuildAction struct {
	TerrainId   string                    `json:"terrainId"`
	TerrainType serialisation.TerrainType `json:"terrainType"`
}

type Action struct {
	ActionType ActionType   `json:"actionType"`
	Move       *MoveAction  `json:"move"`
	Build      *BuildAction `json:"build"`
}

type ExecutableAction interface {
	ExecuteAction(terrain *Terrain)
}

type MoveExecutableAction struct {
	Pipe     *Pipe
	Quantity uint32
}

func (action *MoveExecutableAction) ExecuteAction(terrain *Terrain) {
	action.Pipe.AddToActionQueue(NewSoldierGroup(terrain.owner, action.Quantity), terrain)
	terrain.soldiers.count -= action.Quantity
}

type BuildExecutableAction struct {
	terrainType serialisation.TerrainType
}

func (action *BuildExecutableAction) ExecuteAction(terrain *Terrain) {
	switch action.terrainType {
	case serialisation.Barricade:
		if terrain.state.getTerrainType() == serialisation.Empty {
			var barricadeTS *BarricadeTerrainState = &BarricadeTerrainState{
				terrain,
			}
			var constructionTS *ConstructionTerrainState = &ConstructionTerrainState{
				terrain,
				0,
				barricadeTS,
			}
			terrain.SetTerrainState(constructionTS)
		}
		break
	case serialisation.Factory:
		if terrain.state.getTerrainType() == serialisation.Empty {
			var factoryTS *FactoryTerrainState = &FactoryTerrainState{
				terrain,
				0,
			}
			var constructionTS *ConstructionTerrainState = &ConstructionTerrainState{
				terrain,
				0,
				factoryTS,
			}
			terrain.SetTerrainState(constructionTS)
		}
		break
	case serialisation.Empty:
		if terrain.state.getTerrainType() == serialisation.Barricade || terrain.state.getTerrainType() == serialisation.Factory {
			var emptyTS *EmptyTerrainState = &EmptyTerrainState{
				terrain,
			}
			var constructionTS *ConstructionTerrainState = &ConstructionTerrainState{
				terrain,
				0,
				emptyTS,
			}
			terrain.SetTerrainState(constructionTS)
		}
		break
	default:
		break
	}
}
