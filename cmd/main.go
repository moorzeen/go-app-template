package main

import (
	"github.com/rs/zerolog/log"
	"template/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal().Err(err).Send()
	}
}
