package game

type Game struct {
	players   map[string]*Player
	locations map[string]*Location
}
