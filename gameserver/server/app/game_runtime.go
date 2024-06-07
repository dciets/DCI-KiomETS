package app

import (
	"server/game"
	"server/interfaces"
	"time"
)

type GameRuntime struct {
	maxTickPerGame uint32
	currentGame    *game.Game
	ticking        bool
	timePerTick    uint32
	running        bool
	gameListener   interfaces.Listener
}

func (runtime *GameRuntime) tick() {
	runtime.currentGame.Tick()
	_ = runtime.gameListener.Write(runtime.currentGame.Serialize())
}

func (runtime *GameRuntime) startNewGame() {
	runtime.currentGame = game.NewGame()

	var timePerTick uint32 = runtime.timePerTick

	var ticker *time.Ticker = time.NewTicker(time.Duration(timePerTick) * time.Millisecond)
	var maxTickPerGame uint32 = runtime.maxTickPerGame

	for i := 0; i < int(maxTickPerGame); i++ {
		var start time.Time = time.Now()
		runtime.tick()

		var duration time.Duration = time.Since(start)
		if uint32(duration.Milliseconds()) > timePerTick {
			// TODO : Log time too short
		}
		_ = <-ticker.C
	}
	ticker.Stop()
	// TODO : Collect scores
}

func NewGameRuntime(gameListener interfaces.Listener) *GameRuntime {
	var runtime *GameRuntime = &GameRuntime{
		maxTickPerGame: 12000,
		currentGame:    nil,
		ticking:        false,
		timePerTick:    333,
		running:        false,
		gameListener:   gameListener,
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
