package app

import (
	"errors"
	"github.com/google/uuid"
)

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

func (r *UserRepository) DoesPlayerExist(playerName string) bool {
	for _, player := range r.players {
		if player.name == playerName {
			return true
		}
	}
	return false
}

func (r *UserRepository) AddPlayer(player *Player) {
	r.players[player.id] = player
	r.usedIds = append(r.usedIds, player.id)
}

func doesIdExist(id string, ids []string) bool {
	for i := uint32(0); i < uint32(len(ids)); i++ {
		if ids[i] == id {
			return true
		}
	}
	return false
}

func generateId(ids []string) string {
	var id string
	for {
		var uid uuid.UUID
		uid, _ = uuid.NewUUID()
		id = uid.String()
		if !doesIdExist(id, ids) {
			return id
		}
	}
}

func (r *UserRepository) CreatePlayer(name string) string {
	var id string = generateId(r.usedIds)
	r.AddPlayer(NewPlayer(id, name))
	return id
}
