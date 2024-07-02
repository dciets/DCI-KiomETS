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
	addPlayerMutex     sync.Mutex
	adminId            uint8
}

func NewCommandProcessor(userRepository *UserRepository, gameRuntime *GameRuntime, clientListener interfaces.Listener, adminListener interfaces.Listener, superAdminListener interfaces.Listener) *CommandProcessor {
	return &CommandProcessor{
		userRepository:     userRepository,
		gameRuntime:        gameRuntime,
		clientListener:     clientListener,
		adminListener:      adminListener,
		superAdminListener: superAdminListener,
		mutex:              sync.RWMutex{},
		adminId:            1,
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
			c.addPlayerMutex.Lock()
			c.gameRuntime.AddPlayer(player)
			c.addPlayerMutex.Unlock()
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
			c.gameRuntime.Play(player.Name(), actionsObj)
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
		if len(split) == 2 {
			var ret = c.start()
			_ = c.adminListener.Write(split[1] + " " + ret)
		}
		c.mutex.Unlock()
		break
	case "stop":
		c.mutex.Lock()
		if len(split) == 2 {
			var ret = c.stop()
			_ = c.adminListener.Write(split[1] + " " + ret)
		}
		c.mutex.Unlock()
		break
	case "status":
		c.mutex.Lock()
		if len(split) == 2 {
			var ret = c.status()
			_ = c.adminListener.Write(split[1] + " " + ret)
		}
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
	case "id-assign":
		c.mutex.Lock()
		if len(split) == 1 {
			var id string = strconv.Itoa(int(c.adminId))
			_ = c.adminListener.Write("01 " + id)
			c.adminId++
		}
		c.mutex.Unlock()
		break
	case "new-player":
		c.mutex.Lock()
		if len(split) == 3 {
			var messageId string = split[1]
			var decodedPlayerNameBytes []byte
			decodedPlayerNameBytes, _ = base64.StdEncoding.DecodeString(split[2])
			var decodedPlayerName string = string(decodedPlayerNameBytes)
			if c.userRepository.DoesPlayerExist(decodedPlayerName) {
				_ = c.adminListener.Write(messageId)
			} else {
				var id string = c.userRepository.CreatePlayer(decodedPlayerName)
				_ = c.adminListener.Write(messageId + " " + base64.StdEncoding.EncodeToString([]byte(id)))
			}
		}
		c.mutex.Unlock()
		break
	case "all-player":
		c.mutex.Lock()
		if len(split) == 2 {
			var messageId string = split[1]
			var players = c.userRepository.SerialisePlayers()

			_ = c.adminListener.Write(messageId + " " + base64.StdEncoding.EncodeToString([]byte(players)))
		}
		c.mutex.Unlock()
		break
	case "set-parameters":
		c.mutex.Lock()
		if len(split) == 3 {
			var parameters Parameters
			var decodedParametersBytes []byte
			decodedParametersBytes, _ = base64.StdEncoding.DecodeString(split[2])

			var err error
			err = json.Unmarshal(decodedParametersBytes, &parameters)
			if err == nil {
				c.gameRuntime.SetParameters(parameters)
			}
		}
		c.mutex.Unlock()
		break
	case "get-parameters":
		c.mutex.Lock()
		if len(split) == 2 {
			var messageId string = split[1]
			var parameters Parameters = c.gameRuntime.GetParameters()

			var encodedParametersBytes []byte
			encodedParametersBytes, _ = json.Marshal(parameters)
			_ = c.adminListener.Write(messageId + " " + base64.StdEncoding.EncodeToString(encodedParametersBytes))
		}
		c.mutex.Unlock()
		break
	}
}
