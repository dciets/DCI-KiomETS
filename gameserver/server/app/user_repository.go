package app

import "errors"

type UserRepository struct {
	usedIds []string
	players map[string]*Player
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		usedIds: make([]string, 0),
		players: make(map[string]*Player),
	}
}

func (r *UserRepository) GetPlayer(id string) (*Player, error) {
	var player *Player
	var ok bool
	if player, ok = r.players[id]; ok {
		return player, nil
	}
	return nil, errors.New("Player " + id + " not found")
}

func (r *UserRepository) AddPlayer(player *Player) {
	r.players[player.id] = player
}
