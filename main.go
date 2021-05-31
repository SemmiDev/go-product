package main

import (
	"github.com/SemmiDev/go-product/internal/logger"
	"github.com/SemmiDev/go-product/internal/server"
)

func main() {
	err := server.Start()
	if err != nil {
		logger.Log().Fatal().Err(err).Msg("failed to run server")
	}
}