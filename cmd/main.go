package main

import (
	"TestTask"
	"TestTask/pkg/handler"
	"TestTask/pkg/repository"
	"TestTask/pkg/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if err := initConfig(); err != nil {
		log.Println("error initializing configs:", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		log.Println("error loading env variables:", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})

	if err != nil {
		log.Println("failed initialize db:", err)
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(TestTask.Server)
	go func() {
		if err = srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			log.Println("error occurred when trying to run http server:", err)
		}
	}()
	log.Println("App started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("App shutting down")

	if err = db.Close(); err != nil {
		log.Println("error occurred on close db:", err)
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
