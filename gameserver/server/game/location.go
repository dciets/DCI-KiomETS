package game

import "github.com/google/uuid"

type Location interface {
	Tick()
	PostTick()
	AddToActionQueue(*SoldierGroup, Location)
}

type Node interface {
	Location
	GetId() uuid.UUID
}
