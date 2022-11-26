package game

import (
	"github.com/gorilla/websocket"
	"github.com/grafov/bcast"
	"github.com/sasha-s/go-deadlock"
	"net/http"
	"time"
)

type GameIntro struct {
	Map   GameMap    `json:"map"`
	Ticks []GameTick `json:"ticks"`
}

var livesMu = deadlock.Mutex{}
var lives = map[string]*bcast.Group{}

func ServeLive(id string, w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	// write game intro
	gamesMu.Lock()
	game, ok := games[id]
	gamesMu.Unlock()

	if !ok {
		go conn.Close()
		return
	}

	game.mu.Lock()
	intro := GameIntro{
		Map:   game.Setting.Map,
		Ticks: game.Ticks,
	}
	game.mu.Unlock()

	err = conn.WriteJSON(Message{
		Event: "intro",
		Data:  intro,
	})
	if err != nil {
		conn.Close()
		return
	}

	// join group
	livesMu.Lock()
	group, ok := lives[id]
	if !ok {
		go conn.Close()
		return
	}
	livesMu.Unlock()

	if !ok {
		conn.Close()
		return
	}

	go func() {
		self := group.Join()
		ticker := time.NewTicker(time.Second * 10)
		defer func() {
			_ = conn.Close()
			self.Close()
			ticker.Stop()
		}()
		for {
			select {
			case <-ticker.C:
				err := conn.WriteMessage(websocket.PingMessage, []byte{})
				if err != nil {
					return
				}
			case msg := <-self.Read:
				tick, ok := msg.(GameTick)
				if ok {
					err = conn.WriteJSON(Message{
						Event: "update",
						Data:  tick,
					})
				} else {
					err = conn.WriteJSON(Message{
						Event: "end",
						Data:  nil,
					})
				}
				if err != nil {
					return
				}
			}
		}
	}()

	defer conn.Close()

	_ = conn.SetReadDeadline(time.Now().Add(time.Second * 60))
	conn.SetPongHandler(func(string) error {
		_ = conn.SetReadDeadline(time.Now().Add(time.Second * 60))
		return nil
	})

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}
