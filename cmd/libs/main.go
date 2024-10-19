package main

//dsn2: "postgres://mkv:book_server@localhost:5432/book_database?sslmode=disable"

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

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Config struct {
	DSN string `yaml:"dsn"`
	//DSN_r    string `yaml:"dsn2"`
	LogLevel int `yaml:"log_level"`
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

	var log2 = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
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

	/*rawSQLConn, err := sql.Open("postgres", systemconfig.DSN)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	m, err := migrate.New(
		"file://migrate",
		systemconfig.DSN_r)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		log.Fatal(err)
	}
	*/
	repo := db.NewRepository(gormDB)

	//r.Use(api.Logging(log))

	r.Use(api.Logging1(log2))

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

	log2.Warn("сервер запущен")
	//fmt.Println("сервер запущен")
	err = http.ListenAndServe("127.0.0.1:8080", r)
	log2.Warn("сервер отключён")
	if err != nil {
		log2.Debug("сервер нe запустился")
	}

}
