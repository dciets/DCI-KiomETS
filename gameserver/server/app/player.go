package app

type Player struct {
	id     string
	name   string
	points uint32
}

func (p *Player) Id() string {
	return p.id
}

func (p *Player) Name() string {
	return p.name
}
