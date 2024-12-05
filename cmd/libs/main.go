package main

//dsт: "host=localhost user=mkv password=book_server dbname=book_database port=5432 sslmode=disable"
//dsn2: "postgres://mkv:book_server@localhost:5432/book_database?sslmode=disable"

import (
	"database/sql"
	"log"

	"gopkg.in/yaml.v3"

	"books1/internal/api"
	"books1/internal/db"

	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Config struct {
	DSN      string `yaml:"dsn"`
	LogLevel int    `yaml:"log_level"`
}

func main() {

	yamlContent, err := os.ReadFile("C:/Users/Konstantin/Desktop/books1/config.yml")
	if err != nil {
		log.Fatal(err)
	}
	var systemconfig Config
	err = yaml.Unmarshal(yamlContent, &systemconfig)
	if err != nil {
		log.Fatal(err)
	}

	var log2 = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.Level(systemconfig.LogLevel),
	}))

	file, err := os.OpenFile("../../app.log", os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	rawSQLConn, err := sql.Open("postgres", systemconfig.DSN)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	r := mux.NewRouter()

	m, err := migrate.New(
		"file://../../migrate",
		systemconfig.DSN)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	repo := db.NewRepository(rawSQLConn)

	r.Use(api.Log(log2))

	ourServer := api.Server{

		Database: repo,
	}

	r.HandleFunc("/book", ourServer.GetBook).Methods(http.MethodGet)
	r.HandleFunc("/book", ourServer.AddBook).Methods(http.MethodPost)
	r.HandleFunc("/book", ourServer.DeleteBook).Methods(http.MethodDelete)
	r.HandleFunc("/book", ourServer.UpdateBook).Methods(http.MethodPut)
	r.HandleFunc("/books", ourServer.AllBooks).Methods(http.MethodGet)

	log2.Warn("сервер запущен")

	err = http.ListenAndServe("localhost:8080", r)
	log2.Warn("сервер отключён")
	if err != nil {
		log2.Debug("сервер нe запустился")
	}

}
