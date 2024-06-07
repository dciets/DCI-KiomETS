package app

import "errors"

type UserRepository struct {
	usedIds []string
	players map[string]*Player
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) GetPlayer(id string) (*Player, error) {
	var player *Player
	var ok bool
	if player, ok = r.players[id]; ok {
		return player, nil
	}
	return nil, errors.New("player not found")
}
