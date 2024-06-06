package game

type Game struct {
	players  map[string]*Player
	terrains map[string]*Terrain
	pipes    []*Pipe
}

func NewGame() *Game {
	return &Game{
		players:  make(map[string]*Player),
		terrains: make(map[string]*Terrain),
		pipes:    make([]*Pipe, 0),
	}
}

func (g *Game) Stop() {

}
