package tests

import (
	"github.com/google/uuid"
	"server/game"
	"testing"
)

func TestBattleOnNonBarricadeTerrainWithOneAttackerWithEqualForce(t *testing.T) {
	var p1 *game.Player = &game.Player{}
	var p2 *game.Player = &game.Player{}

	var g1 *game.SoldierGroup = game.NewSoldierGroup(p2, 2)

	var terrain *game.Terrain = game.NewTerrain(uuid.Max.String())

	terrain.SetOwner(p1, 1)
	terrain.AddToActionQueue(g1, nil)

	terrain.PostTick()

	if terrain.GetOwner() != p1 {
		t.Fatal("Terrain owner shouldn't have change")
	}
	if terrain.GetNumberOfSoldiers() != 0 {
		t.Fatal("All soldier should be dead")
	}
}

func TestBattleOnNonBarricadeTerrainWithOneAttackerWithWinningAttacker(t *testing.T) {
	var p1 *game.Player = &game.Player{}
	var p2 *game.Player = &game.Player{}

	var g1 *game.SoldierGroup = game.NewSoldierGroup(p2, 3)

	var terrain *game.Terrain = game.NewTerrain(uuid.Max.String())

	terrain.SetOwner(p1, 1)
	terrain.AddToActionQueue(g1, nil)

	terrain.PostTick()

	if terrain.GetOwner() != p2 {
		t.Fatal("Terrain owner should be P2")
	}

	if terrain.GetNumberOfSoldiers() != 1 {
		t.Fatal("One soldier of the attacker should remain")
	}
}

func TestBattleOnNonBarricadeTerrainWithTwoAttackersWithEqualForce(t *testing.T) {
	var p1 *game.Player = &game.Player{}
	var p2 *game.Player = &game.Player{}
	var p3 *game.Player = &game.Player{}

	var g1 *game.SoldierGroup = game.NewSoldierGroup(p2, 1)
	var g2 *game.SoldierGroup = game.NewSoldierGroup(p3, 1)

	var terrain *game.Terrain = game.NewTerrain(uuid.Max.String())

	terrain.SetOwner(p1, 1)
	terrain.AddToActionQueue(g1, nil)
	terrain.AddToActionQueue(g2, nil)

	terrain.PostTick()

	if terrain.GetOwner() != p1 {
		t.Fatal("Terrain owner shouldn't have change")
	}
	if terrain.GetNumberOfSoldiers() != 0 {
		t.Fatal("All soldier should be dead")
	}
}

func TestBattleOnNonBarricadeTerrainWithTwoAttackerWithWinningAttackerAndNotEnoughSoldier(t *testing.T) {
	var p1 *game.Player = &game.Player{}
	var p2 *game.Player = &game.Player{}
	var p3 *game.Player = &game.Player{}

	var g1 *game.SoldierGroup = game.NewSoldierGroup(p2, 2)
	var g2 *game.SoldierGroup = game.NewSoldierGroup(p3, 2)

	var terrain *game.Terrain = game.NewTerrain(uuid.Max.String())

	terrain.SetOwner(p1, 1)
	terrain.AddToActionQueue(g1, nil)
	terrain.AddToActionQueue(g2, nil)

	terrain.PostTick()

	if terrain.GetOwner() != p1 {
		t.Fatal("Terrain owner shouldn't have change")
	}
	if terrain.GetNumberOfSoldiers() != 0 {
		t.Fatal("All soldier should be dead")
	}
}

func TestBattleOnNonBarricadeTerrainWithTwoAttackersWithWinningAttackerAndEnoughSoldier(t *testing.T) {
	var p1 *game.Player = &game.Player{}
	var p2 *game.Player = &game.Player{}
	var p3 *game.Player = &game.Player{}

	var g1 *game.SoldierGroup = game.NewSoldierGroup(p2, 3)
	var g2 *game.SoldierGroup = game.NewSoldierGroup(p3, 2)

	var terrain *game.Terrain = game.NewTerrain(uuid.Max.String())

	terrain.SetOwner(p1, 1)
	terrain.AddToActionQueue(g1, nil)
	terrain.AddToActionQueue(g2, nil)

	terrain.PostTick()

	if terrain.GetOwner() != p2 {
		t.Fatal("Terrain owner should be P2")
	}
	if terrain.GetNumberOfSoldiers() != 1 {
		t.Fatal("One soldier of the attacker should remain")
	}
}

func TestBattleOnBarricadeTerrainWithOneAttackerWithEqualForce(t *testing.T) {
	var p1 *game.Player = &game.Player{}
	var p2 *game.Player = &game.Player{}

	var g1 *game.SoldierGroup = game.NewSoldierGroup(p2, 7)

	var terrain *game.Terrain = game.NewTerrain(uuid.Max.String())

	var state game.TerrainState = &game.BarricadeTerrainState{
		T: terrain,
	}
	terrain.SetTerrainState(state)

	terrain.SetOwner(p1, 1)
	terrain.AddToActionQueue(g1, nil)

	terrain.PostTick()

	if terrain.GetOwner() != p1 {
		t.Fatal("Terrain owner shouldn't have change")
	}
	if terrain.GetNumberOfSoldiers() != 0 {
		t.Fatal("All soldier should be dead")
	}
}

func TestBattleOnBarricadeTerrainWithOneAttackerWithWinningAttacker(t *testing.T) {
	var p1 *game.Player = &game.Player{}
	var p2 *game.Player = &game.Player{}

	var g1 *game.SoldierGroup = game.NewSoldierGroup(p2, 8)

	var terrain *game.Terrain = game.NewTerrain(uuid.Max.String())

	var state game.TerrainState = &game.BarricadeTerrainState{
		T: terrain,
	}
	terrain.SetTerrainState(state)

	terrain.SetOwner(p1, 1)
	terrain.AddToActionQueue(g1, nil)

	terrain.PostTick()

	if terrain.GetOwner() != p2 {
		t.Fatal("Terrain should be P2")
	}
	if terrain.GetNumberOfSoldiers() != 1 {
		t.Fatal("One soldier of the attacker should remain")
	}
}

func TestBattleOnNonBarricadeTerrainWithZeroDefender(t *testing.T) {
	var p1 *game.Player = &game.Player{}
	var p2 *game.Player = &game.Player{}

	var g1 *game.SoldierGroup = game.NewSoldierGroup(p2, 1)

	var terrain *game.Terrain = game.NewTerrain(uuid.Max.String())

	terrain.SetOwner(p1, 0)
	terrain.AddToActionQueue(g1, nil)

	terrain.PostTick()

	if terrain.GetOwner() != p2 {
		t.Fatal("Terrain owner should have changed")
	}
	if terrain.GetNumberOfSoldiers() != 1 {
		t.Fatal("No soldier should be dead")
	}
}

func TestBattleOnBarrierTerrainWithZeroDefender(t *testing.T) {
	var p1 *game.Player = &game.Player{}
	var p2 *game.Player = &game.Player{}

	var g1 *game.SoldierGroup = game.NewSoldierGroup(p2, 1)

	var terrain *game.Terrain = game.NewTerrain(uuid.Max.String())

	var state game.TerrainState = &game.BarricadeTerrainState{
		T: terrain,
	}
	terrain.SetTerrainState(state)

	terrain.SetOwner(p1, 0)
	terrain.AddToActionQueue(g1, nil)

	terrain.PostTick()

	if terrain.GetOwner() != p2 {
		t.Fatal("Terrain owner should have changed")
	}
	if terrain.GetNumberOfSoldiers() != 1 {
		t.Fatal("No soldier should be dead")
	}
}
