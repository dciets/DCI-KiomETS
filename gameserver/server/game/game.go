package game

import (
	"encoding/json"
	"github.com/google/uuid"
	"math"
	"math/rand"
	"server/game/serialisation"
	"time"
)

type Game struct {
	players           map[string]*Player
	terrains          map[string]*Terrain
	terrainsNeighbors map[string][]Neighbor
	terrainsArray     []*Terrain
	pipes             []*Pipe
	soldierCreation   uint32
	soldierMove       uint32
	terrainChange     uint32
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

type Neighbor struct {
	Name string
	Pipe *Pipe
}

func addNeighbor(t1 *Terrain, t2 *Terrain, terrainNeighbors map[string][]Neighbor, pipe *Pipe) {
	n, ok := terrainNeighbors[t1.GetId()]
	if !ok {
		terrainNeighbors[t1.GetId()] = make([]Neighbor, 0)
	}
	terrainNeighbors[t1.GetId()] = append(n, Neighbor{t2.GetId(), pipe})
	n, ok = terrainNeighbors[t2.GetId()]
	if !ok {
		terrainNeighbors[t2.GetId()] = make([]Neighbor, 0)
	}
	terrainNeighbors[t2.GetId()] = append(n, Neighbor{t1.GetId(), pipe})
}

func generateNodesAndPipes(mapSize uint32) (map[string]*Terrain, []*Terrain, []*Pipe, map[string][]Neighbor) {
	var realSize uint32 = 1
	for i := uint32(0); i <= mapSize; i++ {
		realSize += i * 6
	}
	var terrains map[string]*Terrain = make(map[string]*Terrain)
	var terrainsArray []*Terrain = make([]*Terrain, realSize)
	var pipes []*Pipe = make([]*Pipe, 0)
	var ids []string = make([]string, realSize)
	var terrainNeighbors map[string][]Neighbor = make(map[string][]Neighbor)

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

			var rAdded uint32 = 0
			for r2 := uint32(0); r2 < mapSize; r2++ {
				rAdded += mapSize - r2
			}

			for i := uint32(0); i < 6; i++ {
				var baseX float32 = arr[i][0] * float32(rAdded)
				var baseY float32 = arr[i][1] * float32(rAdded)

				var endX float32 = arr[(i+1)%6][0] * float32(rAdded)
				var endY float32 = arr[(i+1)%6][1] * float32(rAdded)
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

	for i := uint32(1); i < 7; i++ {
		var pipe *Pipe = NewPipe(terrainsArray[0], terrainsArray[i], uint8(mapSize))
		pipes = append(pipes, pipe)
		addNeighbor(terrainsArray[0], terrainsArray[i], terrainNeighbors, pipe)
	}

	for r := uint32(1); r < mapSize; r++ {
		var minV uint32 = uint32(math.Pow(6, float64(r-1))) + 1
		if r == 1 {
			minV = 1
		}
		var maxV uint32 = minV + r*6 - 1

		for i := minV; i <= maxV; i++ {
			var indexM1 uint32 = (((i - minV) - 1) % ((r + 1) * 6)) + maxV + 1
			var index uint32 = ((i - minV) % ((r + 1) * 6)) + maxV + 1
			var indexP1 uint32 = (((i - minV) - 1) % ((r + 1) * 6)) + maxV + 1

			var p1 *Pipe = NewPipe(terrainsArray[i], terrainsArray[indexM1], uint8(mapSize-r))
			var p2 *Pipe = NewPipe(terrainsArray[i], terrainsArray[index], uint8(mapSize-r))
			var p3 *Pipe = NewPipe(terrainsArray[i], terrainsArray[indexP1], uint8(mapSize-r))
			pipes = append(pipes, p1)
			pipes = append(pipes, p2)
			pipes = append(pipes, p3)
			addNeighbor(terrainsArray[indexM1], terrainsArray[i], terrainNeighbors, p1)
			addNeighbor(terrainsArray[index], terrainsArray[i], terrainNeighbors, p2)
			addNeighbor(terrainsArray[indexP1], terrainsArray[i], terrainNeighbors, p3)
		}
	}

	return terrains, terrainsArray, pipes, terrainNeighbors
}

func NewGame(mapSize uint32, soldierCreationSpeed uint32, soldierMoveSpeed uint32, terrainChangeSpeed uint32) *Game {
	var terrains map[string]*Terrain
	var terrainsArray []*Terrain
	var pipes []*Pipe
	var terrainsNeighbors map[string][]Neighbor

	terrains, terrainsArray, pipes, terrainsNeighbors = generateNodesAndPipes(mapSize)

	var g *Game = &Game{
		players:           make(map[string]*Player),
		terrains:          terrains,
		terrainsArray:     terrainsArray,
		pipes:             pipes,
		soldierCreation:   soldierCreationSpeed,
		soldierMove:       soldierMoveSpeed,
		terrainChange:     terrainChangeSpeed,
		terrainsNeighbors: terrainsNeighbors,
	}

	return g
}

func (g *Game) Tick() {
	var playerIdToRemove []string = make([]string, 0)
	for _, player := range g.players {
		if len(player.possessedTerrains) == 0 {
			playerIdToRemove = append(playerIdToRemove, player.name)
		}
	}
	for _, name := range playerIdToRemove {
		delete(g.players, name)
	}

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
	for _, terrain := range g.terrainsArray {
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
	var unusedTerrain []string = make([]string, 0)
	for _, terrain := range g.terrainsArray {
		if terrain.owner == nil {
			unusedTerrain = append(unusedTerrain, terrain.id)
		}
	}

	if len(unusedTerrain) > 0 {
		rand.Seed(time.Now().Unix())
		var player *Player = NewPlayer(name, RandomColor())
		g.players[name] = player
		var terrain *Terrain = g.terrains[unusedTerrain[rand.Intn(len(unusedTerrain))]]
		terrain.SetOwner(player, 5)
		terrain.SetTerrainState(&FactoryTerrainState{terrain, 0})
	}
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

func (g *Game) Play(playerName string, actions []Action) {
	var player *Player
	player, _ = g.players[playerName]

	for _, action := range actions {
		var affectedTerrain *Terrain
		var ok bool
		switch action.ActionType {
		case Move:
			affectedTerrain, ok = g.terrains[action.Move.FromId]
			if ok && affectedTerrain.owner == player && affectedTerrain.soldiers.count >= uint32(action.Move.Quantity) {
				var toTerrain *Terrain
				toTerrain, ok = g.terrains[action.Move.ToId]
				if ok {
					var neighbors []Neighbor = g.terrainsNeighbors[affectedTerrain.GetId()]
					var pipe *Pipe = nil
					for _, neighbor := range neighbors {
						if neighbor.Name == toTerrain.GetId() {
							pipe = neighbor.Pipe
							break
						}
					}
					if pipe != nil {
						var moveAction *MoveExecutableAction = &MoveExecutableAction{Pipe: pipe, Quantity: uint32(action.Move.Quantity)}
						affectedTerrain.currentAction = moveAction
					}
				}
			}
		case Build:
			affectedTerrain, ok = g.terrains[action.Build.TerrainId]
			if ok && affectedTerrain.owner == player {
				if action.Build.TerrainType == serialisation.Barricade ||
					action.Build.TerrainType == serialisation.Factory ||
					action.Build.TerrainType == serialisation.Empty {
					var buildAction *BuildExecutableAction = &BuildExecutableAction{
						action.Build.TerrainType,
					}
					affectedTerrain.currentAction = buildAction
				}
			}
		}
	}
}
