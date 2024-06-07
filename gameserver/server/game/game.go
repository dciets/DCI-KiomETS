package game

type Game struct {
	players  map[string]*Player
	terrains map[string]*Terrain
	pipes    []*Pipe
}

func NewGame() *Game {
	var g *Game = &Game{
		players:  make(map[string]*Player),
		terrains: make(map[string]*Terrain),
		pipes:    make([]*Pipe, 0),
	}

	return g
}

func (g *Game) Tick() {
	for _, terrain := range g.terrains {
		terrain.Tick()
	}
	for _, pipe := range g.pipes {
		pipe.Tick()
	}

	for _, terrain := range g.terrains {
		terrain.PostTick()
	}
	for _, pipe := range g.pipes {
		pipe.PostTick()
	}
}

func (g *Game) Serialize() string {
	return "[]"
}
