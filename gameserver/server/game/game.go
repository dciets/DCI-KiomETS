package game

import (
	"encoding/json"
	"github.com/google/uuid"
	"math"
	"server/game/serialisation"
)

type Game struct {
	players         map[string]*Player
	terrains        map[string]*Terrain
	pipes           []*Pipe
	soldierCreation uint32
	soldierMove     uint32
	terrainChange   uint32
}

func doesIdExist(id string, ids []string, maxIndex uint32) bool {
	for i := uint32(0); i < maxIndex; i++ {
		if ids[i] == id {
			return true
		}
	}
	return false
}

func generateId(ids []string, currentIndex uint32) string {
	var id string
	for {
		var uid uuid.UUID
		uid, _ = uuid.NewUUID()
		id = uid.String()
		if !doesIdExist(id, ids, currentIndex) {
			ids[currentIndex] = id
			return id
		}
	}
}

func generateNodesAndPipes(mapSize uint32) (map[string]*Terrain, []*Pipe) {
	var realSize uint32 = uint32(math.Pow(6, float64(mapSize))) + 1
	var terrains map[string]*Terrain = make(map[string]*Terrain)
	var terrainsArray []*Terrain = make([]*Terrain, realSize)
	var pipes []*Pipe = make([]*Pipe, 0)
	var ids []string = make([]string, realSize)

	var side1 [2]float32 = [2]float32{1, 0}
	var up1 [2]float32 = [2]float32{0.5, 0.866}
	var up2 [2]float32 = [2]float32{-0.5, 0.866}
	var side2 [2]float32 = [2]float32{-1, 0}
	var down1 [2]float32 = [2]float32{-0.5, -0.866}
	var down2 [2]float32 = [2]float32{0.5, -0.866}

	var arr [6][2]float32 = [6][2]float32{side1, up1, up2, side2, down1, down2}

	for r := uint32(0); r <= mapSize; r++ {
		if r == 0 {
			terrainsArray[0] = NewTerrain(generateId(ids, 0))
			terrains[terrainsArray[0].GetId()] = terrainsArray[0]
			terrainsArray[0].SetPosition([2]float32{0, 0})
		} else {
			var minV uint32 = uint32(math.Pow(6, float64(r-1))) + 1
			if r == 1 {
				minV = 1
			}

			for i := uint32(0); i < 6; i++ {
				var baseX float32 = arr[i][0]
				var baseY float32 = arr[i][1]

				var endX float32 = arr[(i+1)%6][0]
				var endY float32 = arr[(i+1)%6][1]
				var dx float32 = (endX - baseX) / float32(r)
				var dy float32 = (endY - baseY) / float32(r)

				for j := uint32(0); j < r; j++ {
					var index uint32 = minV + i*r + j
					terrainsArray[index] = NewTerrain(generateId(ids, index))
					terrains[terrainsArray[index].GetId()] = terrainsArray[index]
					terrainsArray[index].SetPosition([2]float32{baseX + float32(j)*dx, baseY + float32(j)*dy})
				}
			}
		}
	}

	for r := uint32(1); r < mapSize-1; r++ {
		var minV uint32 = uint32(math.Pow(6, float64(r-1))) + 1
		if r == 1 {
			minV = 1
		}
		var maxV uint32 = minV + r*6 - 1

		for i := minV; i <= maxV; i++ {
			var indexM1 uint32 = (((i - minV) - 1) % ((r + 1) * 6)) + maxV + 1
			var index uint32 = ((i - minV) % ((r + 1) * 6)) + maxV + 1
			var indexP1 uint32 = (((i - minV) - 1) % ((r + 1) * 6)) + maxV + 1

			pipes = append(pipes, NewPipe(terrainsArray[i], terrainsArray[indexM1], uint8(mapSize-r)))
			pipes = append(pipes, NewPipe(terrainsArray[i], terrainsArray[index], uint8(mapSize-r)))
			pipes = append(pipes, NewPipe(terrainsArray[i], terrainsArray[indexP1], uint8(mapSize-r)))
		}
	}

	return terrains, pipes
}

func NewGame(mapSize uint32, soldierCreationSpeed uint32, soldierMoveSpeed uint32, terrainChangeSpeed uint32) *Game {
	var terrains map[string]*Terrain
	var pipes []*Pipe

	terrains, pipes = generateNodesAndPipes(mapSize)

	var g *Game = &Game{
		players:         make(map[string]*Player),
		terrains:        terrains,
		pipes:           pipes,
		soldierCreation: soldierCreationSpeed,
		soldierMove:     soldierMoveSpeed,
		terrainChange:   terrainChangeSpeed,
	}

	return g
}

func (g *Game) Tick() {
	for _, terrain := range g.terrains {
		terrain.Tick()
	}
	for i := uint32(0); i < g.soldierMove; i++ {
		for _, pipe := range g.pipes {
			pipe.Tick()
		}
	}

	for _, terrain := range g.terrains {
		terrain.PostTick()
	}
	for _, pipe := range g.pipes {
		pipe.PostTick()
	}
}

func (g *Game) Serialize() string {
	var serialisedGame *serialisation.GameSerialisation = serialisation.NewGameSerialisation(len(g.players), len(g.terrains), len(g.pipes))
	var serialisationBytes []byte

	var playerNameIndexMap = make(map[string]int)
	var terrainIdIndexMap = make(map[string]int)

	var i int = 0
	for _, player := range g.players {
		playerNameIndexMap[player.name] = i
		serialisedGame.Players[i] = player.Serialize()
		i++
	}

	i = 0
	for _, terrain := range g.terrains {
		terrainIdIndexMap[terrain.id] = i
		serialisedGame.Terrains[i] = terrain.Serialize(playerNameIndexMap)
		i++
	}

	i = 0
	for _, pipe := range g.pipes {
		serialisedGame.Pipes[i] = pipe.Serialize(playerNameIndexMap, terrainIdIndexMap)
		i++
	}

	serialisationBytes, _ = json.Marshal(serialisedGame)
	return string(serialisationBytes)
}

func (g *Game) HasPlayer(name string) bool {
	var ok bool
	_, ok = g.players[name]
	return ok
}

func (g *Game) AddPlayer(name string) {
	g.players[name] = NewPlayer(name, RandomColor())
}

func (g *Game) CalculatePlayerScore(name string) uint32 {
	var player *Player
	var ok bool
	var points uint32 = 0
	player, ok = g.players[name]
	if !ok {
		return 0
	}

	points += uint32(len(player.possessedTerrains))
	points += player.numberOfKill * 5

	return points
}
