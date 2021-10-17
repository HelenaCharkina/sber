package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"sber/pkg/handler"
	"sber/pkg/repository"
	"sber/pkg/service"
	"sber/pkg/settings"
	"syscall"
)

// @title Test app for Sber
// @version 1.0
// @description API Server for Employees Tree

// @host localhost:9000
// @BasePath /

func main() {
	ctx := context.TODO()
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("Env variables loading error: %s", err)
	}

	if err := settings.InitConfig(); err != nil {
		logrus.Fatalf("Config initialization error: %s", err)
	}

	client, err := repository.NewMongoDB(ctx)
	if err != nil {
		logrus.Fatalf("DB initialization error: %s", err)
	}

	repo := repository.NewRepository(client)
	services := service.NewService(repo)
	handlers := handler.NewHandler(services)

	srv := new(Server)
	go func() {
		if err := srv.Run(settings.Config.Port, handlers.InitRoutes()); err != nil {
			logrus.Fatalf("Server running error: %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	if err = srv.Shutdown(ctx); err != nil {
		logrus.Fatalf("Server shutting down error: %s", err)
	}
	if err = client.Disconnect(ctx); err != nil {
		logrus.Fatalf("DB connection close error: %s", err)
	}
}
