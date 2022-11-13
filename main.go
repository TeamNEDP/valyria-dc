package main

import (
	"log"
	"math/rand"
	"time"
	"valyria-dc/game"
	"valyria-dc/services"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	game.InitWatchdog()

	err := services.Start("0.0.0.0:8000")
	log.Fatalln(err.Error())
}
