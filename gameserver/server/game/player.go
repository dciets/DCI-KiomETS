package game

import (
	"fmt"
	"math/rand"
	"server/game/serialisation"
)

type Color struct {
	r, g, b uint8
}

func RandomColor() Color {
	var r, g, b uint8
	r = uint8(rand.Uint32())
	g = uint8(rand.Uint32())
	b = uint8(rand.Uint32())
	return Color{r, g, b}
}

func (c Color) String() string {
	return fmt.Sprintf("#%02x%02x%02x", c.r, c.g, c.b)
}

type Player struct {
	name              string
	color             Color
	numberOfKill      uint32
	possessedTerrains []string
}

func NewPlayer(name string, color Color) *Player {
	return &Player{name: name, color: color, numberOfKill: 0, possessedTerrains: []string{}}
}

func (p *Player) Serialize() serialisation.PlayerSerialisation {
	return *serialisation.NewPlayerSerialisation(p.name, p.color.String(), int(p.numberOfKill), len(p.possessedTerrains))
}

func (p *Player) AddTerrain(terrainId string) {
	p.possessedTerrains = append(p.possessedTerrains, terrainId)
}

func (p *Player) RemoveTerrain(terrainId string) {
	for i, t := range p.possessedTerrains {
		if t == terrainId {
			p.possessedTerrains = append(p.possessedTerrains[:i], p.possessedTerrains[i+1:]...)
			return
		}
	}
}
