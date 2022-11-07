package main

import (
	"log"
	"valyria-dc/services"
)

func main() {
	err := services.Start("0.0.0.0:8000")
	log.Fatalln(err.Error())
}
