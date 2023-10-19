package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	error_message "github.com/MarselBissengaliyev/advertising/pkg/error_message"
	"github.com/MarselBissengaliyev/advertising/pkg/handler"
	"github.com/MarselBissengaliyev/advertising/pkg/model"
	"github.com/MarselBissengaliyev/advertising/pkg/repository"
	"github.com/MarselBissengaliyev/advertising/pkg/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// @title Advertising API
// @version 1.0
// @description Api Server for Advertising Application

// @contact.name Marsel Bissengaliyev
// @contact.url https://t.me/marsel_bisengaliev
// @contact.email marselbisengaliev1@gmail.com

// @host localhost:8000
// @BasePath /
func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil {
		logrus.Fatalf("%s: %s", error_message.ErrorInitConfig, err.Error())
	}

	if err := godotenv.Load("configs/.env"); err != nil {
		logrus.Fatalf("%s: %s", error_message.ErrorLoadEnvVars, err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		UserName: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})

	if err != nil {
		logrus.Fatalf("%s: %s", error_message.FailedToInitDB, err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(model.Server)
	go func() {
		port := viper.GetString("port")
		if err := srv.Run(port, handlers.InitRoutes()); err != nil {
			logrus.Fatalf("%s: %s", error_message.RunHttpServerError, err.Error())
		}
	}()

	logrus.Print("Advertising App Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("Advertising App Shutting Down")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("%s: %s", error_message.ShutDownServerError, err.Error())
	}

	if err := db.Close(); err != nil {
		logrus.Errorf("%s: %s", error_message.CloseDBConnectionError, err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigType("yml")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
