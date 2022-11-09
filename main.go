package main

import (
	"log"
	"valyria-dc/game"
	"valyria-dc/services"
)

func main() {
	game.InitWatchdog()

	err := services.Start("0.0.0.0:8000")
	log.Fatalln(err.Error())
}
