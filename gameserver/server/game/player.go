package game

type Color struct {
	r, g, b uint8
}

type Player struct {
	name  string
	color Color
}

func NewPlayer(name string, color Color) *Player {
	return &Player{name: name, color: color}
}

func (p *Player) Name() string {
	return p.name
}

func (p *Player) Color() Color {
	return p.color
}
