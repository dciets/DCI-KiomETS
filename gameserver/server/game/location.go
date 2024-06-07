package game

type Location interface {
	Tick()
	PostTick()
	AddToActionQueue(*SoldierGroup, Location)
}

type Node interface {
	Location
	GetId() string
}
