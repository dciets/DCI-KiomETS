package app

import (
	"log"
	"server/communications"
	"server/game"
	"server/interfaces"
	"time"
)

type Parameters struct {
	MapSize              uint `json:"mapSize"`
	SoldierSpeed         uint `json:"soldierSpeed"`
	SoldierCreationSpeed uint `json:"soldierCreationSpeed"`
	TerrainChangeSpeed   uint `json:"terrainChangeSpeed"`
	GameLength           uint `json:"gameLength"`
}

type GameRuntime struct {
	maxTickPerGame       uint32
	currentGame          *game.Game
	ticking              bool
	timePerTick          uint32
	running              bool
	gameListener         interfaces.Listener
	mapSize              uint32
	soldierMovementSpeed uint32
	soldierCreationSpeed uint32
	terrainChangeSpeed   uint32
	userRepository       *UserRepository
}

func (runtime *GameRuntime) tick() {
	runtime.currentGame.Tick()

	_ = runtime.gameListener.Write(communications.NewMessage("action", runtime.currentGame.Serialize()).String())
}

func (runtime *GameRuntime) startNewGame() {
	runtime.currentGame = game.NewGame(runtime.mapSize, runtime.soldierCreationSpeed, runtime.soldierMovementSpeed, runtime.terrainChangeSpeed)

	var timePerTick uint32 = runtime.timePerTick

	var ticker *time.Ticker = time.NewTicker(time.Duration(timePerTick) * time.Millisecond)
	var maxTickPerGame uint32 = runtime.maxTickPerGame

	for i := 0; i < int(maxTickPerGame); i++ {
		var start time.Time = time.Now()
		runtime.tick()

		var duration time.Duration = time.Since(start)
		if uint32(duration.Milliseconds()) > timePerTick {
			log.Print("UPDATE TOO LONG")
		}
		_ = <-ticker.C
	}
	ticker.Stop()
	_ = runtime.gameListener.Write(communications.NewMessage("end", "").String())

	for _, player := range runtime.userRepository.players {
		var pt uint32 = runtime.currentGame.CalculatePlayerScore(player.name)
		player.points += pt
	}

	_ = runtime.gameListener.Write(communications.NewMessage("scoreboard", runtime.userRepository.SerialisePlayersScore()).String())

	runtime.currentGame = nil
}

func NewGameRuntime(gameListener interfaces.Listener, repository *UserRepository) *GameRuntime {
	var runtime *GameRuntime = &GameRuntime{
		maxTickPerGame:       12000,
		currentGame:          nil,
		ticking:              false,
		timePerTick:          333,
		running:              false,
		gameListener:         gameListener,
		mapSize:              1,
		soldierMovementSpeed: 1,
		soldierCreationSpeed: 1,
		terrainChangeSpeed:   1,
		userRepository:       repository,
	}

	return runtime
}

func (runtime *GameRuntime) Status() bool {
	return runtime.running
}

func (runtime *GameRuntime) Start() {
	runtime.ticking = true
	runtime.running = true

	for runtime.ticking {
		<-time.NewTimer(time.Duration(1) * time.Second).C
		runtime.startNewGame()
	}

	runtime.running = false
}

func (runtime *GameRuntime) Stop() {
	runtime.ticking = false
}

func (runtime *GameRuntime) SetMaxTick(maxTick uint32) {
	runtime.maxTickPerGame = maxTick
}

func (runtime *GameRuntime) SetTimePerTick(timePerTick uint32) {
	runtime.timePerTick = timePerTick
}

func (runtime *GameRuntime) HasPlayer(player *Player) bool {
	if !runtime.running {
		return true
	}

	return runtime.currentGame.HasPlayer(player.Name())
}

func (runtime *GameRuntime) AddPlayer(player *Player) {
	if runtime.running {
		runtime.currentGame.AddPlayer(player.Name())
	}
}

func (runtime *GameRuntime) Play(playerName string, actions []game.Action) {
	if runtime.running {
		runtime.currentGame.Play(playerName, actions)
	}
}

func (runtime *GameRuntime) SetParameters(parameters Parameters) {
	runtime.maxTickPerGame = uint32(parameters.GameLength)
	runtime.soldierMovementSpeed = uint32(parameters.SoldierSpeed)
	runtime.soldierCreationSpeed = uint32(parameters.SoldierCreationSpeed)
	runtime.terrainChangeSpeed = uint32(parameters.TerrainChangeSpeed)
	runtime.mapSize = uint32(parameters.MapSize)
}

func (runtime *GameRuntime) GetParameters() Parameters {
	return Parameters{
		MapSize:              uint(runtime.mapSize),
		SoldierSpeed:         uint(runtime.soldierMovementSpeed),
		SoldierCreationSpeed: uint(runtime.soldierCreationSpeed),
		TerrainChangeSpeed:   uint(runtime.terrainChangeSpeed),
		GameLength:           uint(runtime.maxTickPerGame),
	}
}
