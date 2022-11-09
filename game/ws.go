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
	"valyria-dc/services"
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
				if ok && game.allocatedSession == session.id {
					game.Ticks = append(game.Ticks, data.Tick)
					// TODO: trigger live update
				}
				gamesMu.Unlock()
			} else if message.Event == "gameEnd" {
				data := GameEndData{}
				err := mapstructure.Decode(message.Data, &data)
				if err != nil {
					log.Printf("Invalid gameEnd data received from simulator: %v\n", err)
					return
				}
				gamesMu.Lock()
				services.HandleGameEnd(games[data.ID], data.Result)
				delete(games, data.ID)
				gamesMu.Unlock()
			}
		}
	}()
}
