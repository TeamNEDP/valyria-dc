package game

import (
	"github.com/grafov/bcast"
	"log"
	"sync"
	"time"
)

type GameProcess struct {
	ID               string
	Setting          GameSetting
	Ticks            []GameTick
	lastUpdated      time.Time
	allocatedSession string
	mu               sync.Mutex
}

var gamesMu = sync.Mutex{}
var games = map[string]*GameProcess{}

func StartGame(id string, setting GameSetting) {
	process := GameProcess{
		ID:               id,
		Setting:          setting,
		Ticks:            nil,
		lastUpdated:      time.Now(),
		allocatedSession: "",
	}

	gamesMu.Lock()
	games[id] = &process
	gamesMu.Unlock()

	livesMu.Lock()
	oldLive, ok := lives[id]
	if ok {
		oldLive.Close()
	}
	lives[id] = bcast.NewGroup()
	go lives[id].Broadcast(0)
	livesMu.Unlock()

	go allocateGame(&process)
}

func IsRunning(id string) bool {
	gamesMu.Lock()
	game, ok := games[id]
	gamesMu.Unlock()
	if !ok {
		return false
	}
	game.mu.Lock()
	defer func() {
		game.mu.Unlock()
	}()
	return games[id].allocatedSession != ""
}

func allocateGame(process *GameProcess) {
	process.mu.Lock()
	defer process.mu.Unlock()
	sessionsMu.Lock()
	defer sessionsMu.Unlock()

	process.Ticks = nil
	process.lastUpdated = time.Now()

	for _, session := range sessions {
		session.mu.Lock()
		if !session.authorized {
			continue
		}
		if session.slots < session.running {
			continue
		}
		session.running++
		process.allocatedSession = session.id
		go func() {
			session.mu.Lock()
			_ = session.conn.WriteJSON(Message{
				Event: "gameStart",
				Data: GameStartData{
					ID:      process.ID,
					Setting: process.Setting,
				},
			})
			session.mu.Unlock()
			//log.Printf("Sent game start message\n")
		}()
		session.mu.Unlock()
		return
	}

	log.Printf("no available simulator for game %s\n", process.ID)
}

func InitWatchdog() {
	ticker := time.NewTicker(time.Second * 5)
	go func() {
		defer ticker.Stop()
		for {
			<-ticker.C
			gamesMu.Lock()
			for _, game := range games {
				game := game
				go func() {
					game.mu.Lock()
					if game.lastUpdated.Add(time.Second * 5).Before(time.Now()) {
						game.lastUpdated = time.Now()
						log.Printf("Reallocating game %s\n", game.ID)
						go allocateGame(game)
					}
					game.mu.Unlock()
				}()
			}
			gamesMu.Unlock()
		}
	}()
}
