package app

type PlayerSerialisation struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type PlayerScoreSerialisation struct {
	Name  string `json:"name"`
	Score uint   `json:"score"`
}

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

func (p *Player) Serialise() PlayerSerialisation {
	return PlayerSerialisation{
		Id:   p.id,
		Name: p.name,
	}
}

func (p *Player) SerialiseScore() PlayerScoreSerialisation {
	return PlayerScoreSerialisation{
		Name:  p.name,
		Score: uint(p.points),
	}
}
