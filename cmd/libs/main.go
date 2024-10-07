package main

import (
	"log"

	"gopkg.in/yaml.v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"books/internal/api"
	"books/internal/db"
	"books/internal/domain"

	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Config struct {
	DSN      string `yaml:"dsn"`
	LogLevel int    `yaml:"log_level"`
}

func main() {

	yamlContent, err := os.ReadFile("C:/Users/Konstantin/Desktop/books/config.yml")
	if err != nil {
		log.Fatal(err)
	}
	var systemconfig Config
	err = yaml.Unmarshal(yamlContent, &systemconfig)
	if err != nil {
		log.Fatal(err)
	}

	var log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.Level(systemconfig.LogLevel),
	}))

	file, err := os.OpenFile("../../app.log", os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	config := postgres.Open(systemconfig.DSN)
	gormDB, err := gorm.Open(config, &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = gormDB.AutoMigrate(&domain.Book{})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	r := mux.NewRouter()

	repo := db.NewRepository(gormDB)

	//r.Use(api.Logging(log))

	r.Use(api.Logging1(log))

	ourServer := api.Server{
		//Database: db.Repository{
		//	Store: make(map[int]domain.Book),
		//},
		Database: repo,
	}

	r.HandleFunc("/book", ourServer.GetBook).Methods(http.MethodGet)
	r.HandleFunc("/book", ourServer.AddBook).Methods(http.MethodPost)
	r.HandleFunc("/book", ourServer.DeleteBook).Methods(http.MethodDelete)
	r.HandleFunc("/book", ourServer.UpdateBook).Methods(http.MethodPut)
	r.HandleFunc("/books", ourServer.AllBooks).Methods(http.MethodGet)

	log.Warn("сервер запущен")
	//fmt.Println("сервер запущен")
	err = http.ListenAndServe("127.0.0.1:8080", r)
	log.Warn("сервер отключён")
	if err != nil {
		log.Debug("сервер нe запустился")
		log.Error("сервер не запустился")
	}

}
