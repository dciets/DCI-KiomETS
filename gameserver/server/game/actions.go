package game

import "server/game/serialisation"

type ActionType int

const (
	Move ActionType = iota
	Build
)

type MoveAction struct {
	FromId   string `json:"fromId"`
	ToId     string `json:"toId"`
	Quantity uint   `json:"quantity"`
}

type BuildAction struct {
	TerrainId   string                    `json:"terrainId"`
	TerrainType serialisation.TerrainType `json:"terrainType"`
}

type Action struct {
	ActionType ActionType   `json:"actionType"`
	Move       *MoveAction  `json:"move"`
	Build      *BuildAction `json:"build"`
}
