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

	/*
		// создаём структуру которая слушает порты для получениязапроса по grpc( и только)
		ln, err := net.Listen("tcp", "localhost:8080")
				if err != nil { log.Fatal("net.Listen(): %v", err)}

			server := grpc.NewServer(
			grpc.Creds(insecure.NewCredentials()),
			)
		// здесь регестрируем хендеры для grpc
			pb.RegisterBookAPIServer(server, &ourServer)
		// здесь запускаем только grpc сервер
		go func() {
			if err := server.Serve(ln); err != nil {
				log.Fatal("server.Serve(): %v", err)
		}
		}()
									// здесь создаём конструкцию, которая умеет вызывать (grpc) методы
		conn, err := grpc.NewClient("localhost:8080",
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			)
			if err != nil {
				log.Fatal("grpc.NewClient(): %v", err)
			}
		defer conn.Close()
				// TODO: fix me
			// HTTP  сервер не умеет обрабатывать grpc запросы, поэтому используем grpc-gateway
		gw := grpc_gateway_runtime.NewServeMux() // нужно сделать тег grpc_gateway_runtime сдесь в main.go
				// grpc- gateway принимает клиент до grpc сервера и преобразует HTTP (REST) в grpc
				// и перенаправляет через client
		err = pb.RegisterBookAPIHandler(context.TODO(), gw, conn)
																	if err != nil { log.Fatal("RegisterServiceAPIHandler(): %v", err) }
			// web стучится по REST в gateway, который конвертирует из HTTP в grpc
			// тот получил и отправил grpc-client и тот в grpc-server
		gwServer := &http.Server{
			Addr: "0.0.0.0:8080",
			Handler: gw,
		}

	*/

	r.HandleFunc("/book", ourServer.GetBook).Methods(http.MethodGet)
	r.HandleFunc("/book", ourServer.AddBook).Methods(http.MethodPost)
	r.HandleFunc("/book", ourServer.DeleteBook).Methods(http.MethodDelete)
	r.HandleFunc("/book", ourServer.UpdateBook).Methods(http.MethodPut)
	r.HandleFunc("/books", ourServer.AllBooks).Methods(http.MethodGet)

	log2.Warn("сервер запущен")
	// err1 := gwServer.ListenAndServe()
	err = http.ListenAndServe("localhost:8080", r)
	log2.Warn("сервер отключён")
	if err != nil {
		fmt.Println(555)
		log2.Debug("сервер нe запустился")
	}

}
