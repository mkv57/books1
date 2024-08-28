package main

import (
	"bookServer/internal/api"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

func main() {

	file, err := os.OpenFile("app.log", os.O_APPEND, 0666)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	r := mux.NewRouter()

	//minimalLevel := slog.LevelInfo
	// Создаем новый логгер, который будет писать в стандартный вывод
	options := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	handler := slog.NewTextHandler(os.Stdout, options)
	logger := slog.New(handler) // логи отправляю в поток

	handler1 := slog.NewTextHandler(file, options)
	logger1 := slog.New(handler1) // логи отправляю в файл

	r.Use(api.Logging(logger))

	var h = api.API{
		Di: &http.Server{
			//Addr:           ":8080",
			//Handler:        r,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
	}

	r.HandleFunc("/book", h.GetBook).Methods(http.MethodGet)
	r.HandleFunc("/book", h.AddBook).Methods(http.MethodPost)
	r.HandleFunc("/book", h.DeleteBook).Methods(http.MethodDelete)
	r.HandleFunc("/book", h.UpdateBook).Methods(http.MethodPut)
	r.HandleFunc("/books", h.AllBooks).Methods(http.MethodGet)

	logger.Warn("сервер запущен")
	fmt.Println("сервер запущен")
	err = http.ListenAndServe("127.0.0.1:8080", r)
	logger1.Warn("сервер отключён")
	if err != nil {
		logger.Error("сервер нe запустился")
		log.Fatal(err)
	}

}
