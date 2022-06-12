package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/AndrewMislyuk/go-shop-backend/internal/config"
	"github.com/AndrewMislyuk/go-shop-backend/internal/handler"
	"github.com/AndrewMislyuk/go-shop-backend/internal/repository"
	"github.com/AndrewMislyuk/go-shop-backend/internal/service"
	"github.com/AndrewMislyuk/go-shop-backend/pkg/database"
	"github.com/AndrewMislyuk/go-shop-backend/pkg/server"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

const (
	CONFIG_DIR  = "configs"
	CONFIG_FILE = "main"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		logrus.Fatal(err)
	}

	cfg, err := config.New(CONFIG_DIR, CONFIG_FILE)
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Printf("%+v\n", cfg)

	db, err := database.NewPostgresConnection(database.ConnectionInfo{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		Username: cfg.Postgres.User,
		DBName:   cfg.Postgres.Db,
		SSLMode:  cfg.DB.SSLMode,
		Password: cfg.Postgres.Password,
	})
	if err != nil {
		logrus.Fatal(err)
	}

	documentsRepo := repository.NewRepository(db)
	documentsService := service.NewService(documentsRepo)
	handler := handler.NewHandler(documentsService)

	srv := new(server.Server)

	go func() {
		if err := srv.Run(fmt.Sprintf("%d", cfg.Server.Port), handler.InitRouter()); err != nil {
			logrus.Fatal(err)
		}
	}()

	logrus.Infoln("Server has been running...")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Infoln("Server was stopped")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		logrus.Errorf("error occured on db connection close: %s", err.Error())
	}
}
