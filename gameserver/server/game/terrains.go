package game

import "github.com/google/uuid"

type TerrainState interface {
	tick()
	getFightingFunctions() (fightFunction, fightFunction)
}

type Terrain struct {
	id             uuid.UUID
	state          TerrainState
	incomingGroups []*SoldierGroup
	owner          *Player
	soldiers       *SoldierGroup
}

func NewTerrain(id uuid.UUID) *Terrain {
	return &Terrain{
		id:       id,
		state:    &EmptyTerrainState{},
		owner:    nil,
		soldiers: nil,
	}
}

func (t *Terrain) SetOwner(player *Player, soldierCount uint32) {
	t.owner = player
	t.soldiers = NewSoldierGroup(player, soldierCount)
}

func (t *Terrain) Tick() {
	// TODO : Gérer les ordres

	t.state.tick()
}

func (t *Terrain) PostTick() {
	var g1Function fightFunction
	var g2Function fightFunction
	g1Function, g2Function = t.state.getFightingFunctions()
	var lastGroupStanding *SoldierGroup = battleNGroup(t.soldiers, t.incomingGroups, g1Function, g2Function)
	t.soldiers = lastGroupStanding
	if lastGroupStanding.player != t.owner {
		// TODO : Tell the owner
		t.owner = lastGroupStanding.player
	}
	t.incomingGroups = []*SoldierGroup{}
}

func (t *Terrain) GetId() uuid.UUID {
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
