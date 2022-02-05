package main

import (
	"log"
	"time"

	"github.com/abibby/uptime/config"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	config.Init()

	uptimeTicker := time.NewTicker(config.CheckFrequency)
	rotateTicker := time.NewTicker(time.Hour)

	var t time.Time

	for {
		select {
		case t = <-uptimeTicker.C:
			err := uptime(t)
			if err != nil {
				log.Print(err)
				continue
			}

		case t = <-rotateTicker.C:
			err := rotate(t)
			if err != nil {
				log.Print(err)
			}
		}
	}
}
