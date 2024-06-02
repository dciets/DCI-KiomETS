package tests

import (
	"github.com/google/uuid"
	"server/game"
	"testing"
)

type mockNode struct {
	Pipe            *game.Pipe
	SendSoldiers    []*game.SoldierGroup
	ReceiveSoldiers []*game.SoldierGroup
}

func (n *mockNode) Tick() {
	for _, soldier := range n.SendSoldiers {
		n.Pipe.AddToActionQueue(soldier, n)
	}
	n.SendSoldiers = []*game.SoldierGroup{}
}

func (n *mockNode) PostTick() {

}

func (n *mockNode) GetId() uuid.UUID {
	return uuid.Max
}

func (n *mockNode) AddToActionQueue(soldierGroup *game.SoldierGroup, from game.Location) {
	n.ReceiveSoldiers = append(n.ReceiveSoldiers, soldierGroup)
}

func TestPipeMovingSoldier(t *testing.T) {
	var n1 *mockNode = &mockNode{}
	var n2 *mockNode = &mockNode{}

	var pipe *game.Pipe = game.NewPipe(n1, n2, 1)

	n1.Pipe = pipe
	n2.Pipe = pipe

	var player *game.Player = game.NewPlayer("", game.Color{})
	var g1 *game.SoldierGroup = game.NewSoldierGroup(player, 1)
	var g2 *game.SoldierGroup = game.NewSoldierGroup(player, 2)
	n1.SendSoldiers = append(n1.SendSoldiers, g1)

	if pipe.GetGroupCount() != 0 {
		t.Fatal("pipe group count should be 0")
	}

	n1.Tick()
	n2.Tick()
	pipe.Tick()

	if pipe.GetGroupCount() != 1 {
		t.Fatal("pipe group count should be 1")
	}

	n1.SendSoldiers = append(n1.SendSoldiers, g2)

	n1.Tick()
	n2.Tick()
	pipe.Tick()

	if pipe.GetGroupCount() != 1 {
		t.Fatal("pipe group count should be 1")
	}

	if len(n2.ReceiveSoldiers) != 1 {
		t.Fatal("Node 2 received count should be 1")
	}

	if n2.ReceiveSoldiers[0] != g1 {
		t.Fatal("Node 2 didn't receive g1")
	}

	n2.ReceiveSoldiers = []*game.SoldierGroup{}

	n1.Tick()
	n2.Tick()
	pipe.Tick()

	if pipe.GetGroupCount() != 0 {
		t.Fatal("pipe group count should be 0")
	}

	if len(n2.ReceiveSoldiers) != 1 {
		t.Fatal("Node 2 received count should be 1")
	}

	if n2.ReceiveSoldiers[0] != g2 {
		t.Fatal("Node 2 didn't receive g2")
	}

}

func TestPipeBattle(t *testing.T) {
	var n1 *mockNode = &mockNode{}
	var n2 *mockNode = &mockNode{}

	var pipe *game.Pipe = game.NewPipe(n1, n2, 1)

	n1.Pipe = pipe
	n2.Pipe = pipe

	var player1 *game.Player = game.NewPlayer("", game.Color{})
	var player2 *game.Player = game.NewPlayer("", game.Color{})
	var g1 *game.SoldierGroup = game.NewSoldierGroup(player1, 1)
	var g2 *game.SoldierGroup = game.NewSoldierGroup(player2, 2)

	n1.SendSoldiers = append(n1.SendSoldiers, g1)
	n2.SendSoldiers = append(n2.SendSoldiers, g2)

	n1.Tick()
	n2.Tick()
	pipe.Tick()

	if pipe.GetGroupCount() != 2 {
		t.Fatal("pipe group count should be 2")
	}

	n1.PostTick()
	n2.PostTick()
	pipe.PostTick()

	if pipe.GetGroupCount() != 1 {
		t.Fatal("pipe group count should be 1")
	}

	if g1.Count() != 0 {
		t.Fatal("G1 should be dead")
	}

	if g2.Count() != 1 {
		t.Fatal("G2 should have 1 remaining soldier")
	}
}
