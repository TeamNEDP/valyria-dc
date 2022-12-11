package game

import (
	"encoding/json"
	"github.com/gofrs/uuid"
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
	"github.com/sasha-s/go-deadlock"
	"log"
	"net/http"
	"os"
	"time"
)

type simulatorSession struct {
	mu         deadlock.Mutex
	id         string
	authorized bool
	slots      uint
	running    uint
	conn       *websocket.Conn
}

var sessions = map[string]*simulatorSession{}
var sessionsMu = deadlock.Mutex{}

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
	conn.SetPingHandler(func(string) error {
		_ = conn.WriteMessage(websocket.PongMessage, []byte{})
		_ = conn.SetReadDeadline(time.Now().Add(time.Second * 60))
		return nil
	})
	conn.SetPongHandler(func(string) error {
		_ = conn.SetReadDeadline(time.Now().Add(time.Second * 60))
		return nil
	})

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
			gamesMu.Unlock()
			if ok {
				game.mu.Lock()
				game.lastUpdated = time.Now()
				if game.allocatedSession == session.id {
					game.Ticks = append(game.Ticks, data.Tick)
					livesMu.Lock()
					lives[data.ID].Send(data.Tick)
					livesMu.Unlock()
				}
				game.mu.Unlock()
			}
		} else if message.Event == "gameEnd" {
			data := GameEndData{}
			err := mapstructure.Decode(message.Data, &data)
			if err != nil {
				log.Printf("Invalid gameEnd data received from simulator: %v\n", err)
				return
			}
			session.mu.Lock()
			session.running--
			session.mu.Unlock()
			livesMu.Lock()
			lives[data.ID].Send(nil)
			lives[data.ID].Close()
			delete(lives, data.ID)
			livesMu.Unlock()
			gamesMu.Lock()
			game := games[data.ID]
			gamesMu.Unlock()
			game.mu.Lock()
			handleGameEnd(game, data.Result)
			game.mu.Unlock()
			delete(games, data.ID)
		} else {
			log.Printf("Invalid event type received from simulator: %s\n", message.Event)
			return
		}
	}

}

func OnGameEnd(handler func(process *GameProcess, result GameResult)) {
	handleGameEnd = handler
}
