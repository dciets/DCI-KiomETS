package app

import (
	"log"
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
	}
}
