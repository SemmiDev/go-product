package server

import (
	"context"
	"fmt"
	"github.com/SemmiDev/go-product/internal/config"
	"github.com/SemmiDev/go-product/internal/db/mysql"
	"github.com/SemmiDev/go-product/internal/db/redis"
	"github.com/SemmiDev/go-product/internal/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"

)

func Start() error {
	mysqlClient, err := mysql.NewClient()
	if err != nil {
		return err
	}
	defer mysqlClient.Close()

	redisClient, err := redis.NewClient()
	if err != nil {
		return err
	}
	defer redisClient.Close()

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Cfg().AppPort),
		Handler: NewRouter(mysqlClient, redisClient),
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		defer close(idleConnsClosed)

		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		signal.Notify(sigint, syscall.SIGTERM)

		<-sigint

		err := httpServer.Shutdown(context.Background())
		if err != nil {
			logger.Log().Err(err).Msg("failed to shutdown server")
		}
	}()

	logger.Log().Info().Msgf("starting server on %s", httpServer.Addr)
	err = httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	<-idleConnsClosed

	logger.Log().Info().Msg("stopped server gracefully")
	return nil
}
