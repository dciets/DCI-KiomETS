package app

import "github.com/google/uuid"

type Player struct {
	id   uuid.UUID
	name string
}

func (p *Player) Id() string {
	return p.id.String()
}

func (p *Player) Name() string {
	return p.name
}
