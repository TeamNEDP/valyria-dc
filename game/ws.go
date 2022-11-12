package game

import (
	"encoding/json"
	"github.com/gofrs/uuid"
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type simulatorSession struct {
	mu         sync.Mutex
	id         string
	authorized bool
	slots      uint
	running    uint
	conn       *websocket.Conn
}

var sessions = map[string]*simulatorSession{}
var sessionsMu = sync.Mutex{}

var handleGameEnd = func(process *GameProcess, result GameResult) {}

func ServeWs(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	sessionUUID, err := uuid.NewV4()
	if err != nil {
		log.Printf("Failed to generate UUID: %v\n", err)
		_ = conn.Close()
		return
	}
	sessionId := sessionUUID.String()
	session := simulatorSession{
		id:         sessionId,
		authorized: false,
		running:    0,
		slots:      0,
		conn:       conn,
	}

	sessionsMu.Lock()
	sessions[sessionId] = &session
	sessionsMu.Unlock()

	defer func() {
		conn.Close()
		sessionsMu.Lock()
		defer sessionsMu.Unlock()
		delete(sessions, sessionId)
	}()

	_ = conn.SetReadDeadline(time.Now().Add(time.Second * 60))
	conn.SetPongHandler(func(string) error {
		_ = conn.SetReadDeadline(time.Now().Add(time.Second * 60))
		return nil
	})

	go func() {
		defer func() {
			conn.Close()
			sessionsMu.Lock()
			defer sessionsMu.Unlock()
			delete(sessions, sessionId)
		}()
		for {
			typ, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}
			if typ == websocket.PingMessage {
				_ = conn.SetReadDeadline(time.Now().Add(time.Second * 60))
				_ = conn.WriteMessage(websocket.PongMessage, []byte{})
				continue
			}

			message := Message{}
			if err := json.Unmarshal(msg, &message); err != nil {
				log.Printf("Invalid message received from simulator: %v\n", err)
				return
			}

			if message.Event == "auth" {
				data := AuthData{}
				err := mapstructure.Decode(message.Data, &data)
				if err != nil {
					log.Printf("Invalid auth data received from simulator: %v\n", err)
					return
				}
				if os.Getenv("SIMULATOR_RPC_SECRET") != data.Token {
					log.Printf("Invalid token received from simulator\n")
					return
				}
				session.mu.Lock()
				session.authorized = true
				session.slots = data.Slots
				session.mu.Unlock()
				continue
			} else if message.Event == "gameUpdate" {
				data := GameUpdateData{}
				err := mapstructure.Decode(message.Data, &data)
				if err != nil {
					log.Printf("Invalid gameUpdate data received from simulator: %v\n", err)
					return
				}
				gamesMu.Lock()
				game, ok := games[data.ID]
				if ok {
					game.mu.Lock()
					if game.allocatedSession == session.id {
						game.Ticks = append(game.Ticks, data.Tick)
						livesMu.Lock()
						lives[data.ID].Send(data.Tick)
						livesMu.Unlock()
					}
					game.mu.Unlock()
				}
				gamesMu.Unlock()
			} else if message.Event == "gameEnd" {
				data := GameEndData{}
				err := mapstructure.Decode(message.Data, &data)
				if err != nil {
					log.Printf("Invalid gameEnd data received from simulator: %v\n", err)
					return
				}
				livesMu.Lock()
				lives[data.ID].Send(nil)
				lives[data.ID].Close()
				delete(lives, data.ID)
				livesMu.Unlock()
				gamesMu.Lock()
				game := games[data.ID]
				game.mu.Lock()
				handleGameEnd(game, data.Result)
				game.mu.Unlock()
				delete(games, data.ID)
				gamesMu.Unlock()
			} else {
				log.Printf("Invalid event type received from simulator: %s\n", message.Event)
				return
			}
		}
	}()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func OnGameEnd(handler func(process *GameProcess, result GameResult)) {
	handleGameEnd = handler
}
