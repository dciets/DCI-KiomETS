package app

type Player struct {
	id     string
	name   string
	points uint32
}

func NewPlayer(id string, name string) *Player {
	return &Player{id: id, name: name, points: 0}
}

func (p *Player) Id() string {
	return p.id
}

func (p *Player) Name() string {
	return p.name
}
