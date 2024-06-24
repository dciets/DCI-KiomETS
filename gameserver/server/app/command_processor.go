package app

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"server/game"
	"server/interfaces"
	"strconv"
	"strings"
	"sync"
)

type CommandProcessor struct {
	userRepository     *UserRepository
	gameRuntime        *GameRuntime
	clientListener     interfaces.Listener
	adminListener      interfaces.Listener
	superAdminListener interfaces.Listener
	mutex              sync.RWMutex
}

func NewCommandProcessor(userRepository *UserRepository, gameRuntime *GameRuntime, clientListener interfaces.Listener, adminListener interfaces.Listener, superAdminListener interfaces.Listener) *CommandProcessor {
	return &CommandProcessor{
		userRepository:     userRepository,
		gameRuntime:        gameRuntime,
		clientListener:     clientListener,
		adminListener:      adminListener,
		superAdminListener: superAdminListener,
		mutex:              sync.RWMutex{},
	}
}

func (c *CommandProcessor) status() string {
	if c.gameRuntime.Status() {
		return "1"
	}
	return "0"
}

func (c *CommandProcessor) start() string {
	if c.gameRuntime.Status() {
		return "0"
	}
	go c.gameRuntime.Start()
	return "1"
}

func (c *CommandProcessor) stop() string {
	if c.gameRuntime.Status() {
		c.gameRuntime.Stop()
		return "1"
	}
	return "0"
}

func (c *CommandProcessor) connect(id string) bool {
	var player *Player
	var err error
	player, err = c.userRepository.GetPlayer(id)
	if err == nil {
		if !c.gameRuntime.HasPlayer(player) {
			c.gameRuntime.AddPlayer(player)
		}
		return true
	} else {
		log.Printf(err.Error())
		return false
	}
}

func (c *CommandProcessor) play(id string, actions string) {
	log.Printf("actions from %s = %s", id, actions)
	if c.gameRuntime.Status() {
		var actionsObj []game.Action = make([]game.Action, 0)
		var err error

		err = json.Unmarshal([]byte(actions), &actionsObj)

		if err != nil {
			log.Printf("actions from %s got error %s", id, err.Error())
		} else {
			var player *Player
			player, _ = c.userRepository.GetPlayer(id)
			c.gameRuntime.Play(player.Name())
		}
	}
}

func (c *CommandProcessor) adminCreate(name string, id string) {
	log.Printf("Super admin create %s with id %s", name, id)
	c.userRepository.AddPlayer(NewPlayer(id, name))
}

func (c *CommandProcessor) Process(command string) {
	// TODO : Regarder user
	log.Print("Command " + command)
	var split []string = strings.Split(command, " ")
	switch split[0] {
	case "start":
		c.mutex.Lock()
		var ret = c.start()
		_ = c.adminListener.Write(ret)
		c.mutex.Unlock()
		break
	case "stop":
		c.mutex.Lock()
		var ret = c.stop()
		_ = c.adminListener.Write(ret)
		c.mutex.Unlock()
		break
	case "status":
		c.mutex.Lock()
		var ret = c.status()
		_ = c.adminListener.Write(ret)
		c.mutex.Unlock()
		break
	case "set-max-tick":
		c.mutex.Lock()
		if len(split) == 2 {
			var value uint64
			var err error
			value, err = strconv.ParseUint(split[1], 10, 32)
			if err == nil {
				c.gameRuntime.SetMaxTick(uint32(value))
			}
		}
		c.mutex.Unlock()
		break
	case "set-time-per-tick":
		c.mutex.Lock()
		if len(split) == 2 {
			var value uint64
			var err error
			value, err = strconv.ParseUint(split[1], 10, 32)
			if err == nil {
				c.gameRuntime.SetTimePerTick(uint32(value))
			}
		}
		c.mutex.Unlock()
		break
	case "stop-all":
		c.clientListener.Stop()
		c.superAdminListener.Stop()
		c.adminListener.Stop()
		break
	case "action":
		c.mutex.RLock()
		if len(split) == 3 {
			var decodedIdBytes []byte
			var decodedActionBytes []byte
			decodedIdBytes, _ = base64.StdEncoding.DecodeString(split[1])
			decodedActionBytes, _ = base64.StdEncoding.DecodeString(split[2])
			var decodedId string = string(decodedIdBytes)
			var decodedAction string = string(decodedActionBytes)
			if c.connect(decodedId) {
				c.play(decodedId, decodedAction)
			}
		}
		c.mutex.RUnlock()
	case "admin-create":
		c.mutex.Lock()
		if len(split) == 3 {
			var decodedNameBytes []byte
			var decodedIdBytes []byte
			decodedNameBytes, _ = base64.StdEncoding.DecodeString(split[1])
			decodedIdBytes, _ = base64.StdEncoding.DecodeString(split[2])
			var decodedName string = string(decodedNameBytes)
			var decodedId string = string(decodedIdBytes)
			c.adminCreate(decodedName, decodedId)
		}
		c.mutex.Unlock()
	}
}
