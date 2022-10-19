package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	grpc_client "github.com/AndrewMislyuk/go-shop-backend/internal/transport/grpc"
	"github.com/AndrewMislyuk/go-shop-backend/internal/transport/rest/config"
	"github.com/AndrewMislyuk/go-shop-backend/internal/transport/rest/handler"
	"github.com/AndrewMislyuk/go-shop-backend/internal/transport/rest/repository"
	"github.com/AndrewMislyuk/go-shop-backend/internal/transport/rest/service"
	"github.com/AndrewMislyuk/go-shop-backend/pkg/database"
	"github.com/AndrewMislyuk/go-shop-backend/pkg/server"
	"github.com/AndrewMislyuk/go-shop-backend/pkg/storage"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
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

// @title CRUD API Go Shop Backend
// @version 1.0
// @description API for frontend cliend

// @host 159.89.235.180:3000
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
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

	client, err := minio.New(cfg.FileStorageConfig.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.FileStorageConfig.AccessKey, cfg.FileStorageConfig.SecretKey, ""),
		Secure: true,
	})
	if err != nil {
		logrus.Fatal(err)
	}

	auditClient, err := grpc_client.NewClient(9000)
	if err != nil {
		logrus.Fatal(err)
	}

	provider := storage.NewFileStorage(client, cfg.FileStorageConfig.Bucket, cfg.FileStorageConfig.Endpoint)

	documentsRepo := repository.NewRepository(db)
	documentsService := service.NewService(documentsRepo, auditClient, provider)
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
		logrus.Errorf("error occurred on server shutting down: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		logrus.Errorf("error occurred on db connection close: %s", err.Error())
	}
}
