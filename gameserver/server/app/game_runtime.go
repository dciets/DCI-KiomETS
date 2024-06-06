package app

type GameRuntime struct {
	maxTickPerGame uint32
}

func NewGameRuntime() *GameRuntime {
	var runtime *GameRuntime = &GameRuntime{
		maxTickPerGame: 12000,
	}

	return runtime
}

func (runtime *GameRuntime) Start() {

}
