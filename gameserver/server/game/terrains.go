package game

import (
	"server/game/serialisation"
)

type TerrainState interface {
	tick()
	getFightingFunctions() (fightFunction, fightFunction)
	getTerrainType() serialisation.TerrainType
}

type Terrain struct {
	id             string
	state          TerrainState
	incomingGroups []*SoldierGroup
	owner          *Player
	soldiers       *SoldierGroup
	position       [2]float32
	currentAction  ExecutableAction
}

func NewTerrain(id string) *Terrain {
	return &Terrain{
		id:            id,
		state:         &EmptyTerrainState{},
		owner:         nil,
		soldiers:      nil,
		position:      [2]float32{0, 0},
		currentAction: nil,
	}
}

func (t *Terrain) SetOwner(player *Player, soldierCount uint32) {
	t.owner = player
	player.AddTerrain(t.id)
	t.soldiers = NewSoldierGroup(player, soldierCount)
}

func (t *Terrain) SetPosition(position [2]float32) {
	t.position = position
}

func (t *Terrain) Tick() {
	if t.currentAction != nil {
		t.currentAction.ExecuteAction(t)
		t.currentAction = nil
	}

	t.state.tick()
}

func (t *Terrain) PostTick() {
	if t != nil && len(t.incomingGroups) > 0 {
		var g1Function fightFunction
		var g2Function fightFunction
		g1Function, g2Function = t.state.getFightingFunctions()
		var lastGroupStanding *SoldierGroup = battleNGroup(t.soldiers, t.incomingGroups, g1Function, g2Function)
		t.soldiers = lastGroupStanding
		if lastGroupStanding != nil && lastGroupStanding.player != t.owner && lastGroupStanding.player != nil {
			if t.owner != nil {
				t.owner.RemoveTerrain(t.id)
			}
			t.owner = lastGroupStanding.player
			lastGroupStanding.player.AddTerrain(t.id)
		}
		t.incomingGroups = []*SoldierGroup{}
	}
}

func (t *Terrain) GetId() string {
	return t.id
}

func (t *Terrain) AddToActionQueue(soldierGroup *SoldierGroup, from Location) {
	if soldierGroup.player == t.owner {
		t.soldiers.count += soldierGroup.count
	} else {
		t.incomingGroups = append(t.incomingGroups, soldierGroup)
	}
}

func (t *Terrain) GetOwner() *Player {
	return t.owner
}

func (t *Terrain) GetNumberOfSoldiers() uint32 {
	return t.soldiers.count
}

func (t *Terrain) SetTerrainState(state TerrainState) {
	t.state = state
}

func (t *Terrain) Serialize(playerNameIndexMap map[string]int) serialisation.TerrainSerialisation {
	var index int = -1
	var numberOfSoldier uint = 0
	if t.owner != nil {
		index = playerNameIndexMap[t.owner.name]
		numberOfSoldier = uint(t.soldiers.count)
	}

	return *serialisation.NewTerrainSerialisation(t.state.getTerrainType(), index, numberOfSoldier, t.position, t.id)
}
