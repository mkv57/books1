package main

import (
	"bookServer/struct_p"
	"fmt"
	"log"
	"log/slog"
	"os"

	"net/http"

	"github.com/gorilla/mux"
)

var minimalLevel = slog.LevelInfo

var file, err = os.OpenFile("app.log", os.O_APPEND, 0666)

var logger = slog.New(slog.NewTextHandler(file, &slog.HandlerOptions{
	Level: minimalLevel,
}))

func main() {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	r := mux.NewRouter()
	r.HandleFunc("/book", struct_p.GetBook).Methods(http.MethodGet)
	r.HandleFunc("/book", struct_p.AddBook).Methods(http.MethodPost)
	r.HandleFunc("/book", struct_p.DeleteBook).Methods(http.MethodDelete)
	r.HandleFunc("/book", struct_p.UpdateBook).Methods(http.MethodPut)
	r.HandleFunc("/books", struct_p.AllBooks).Methods(http.MethodGet)

	logger.Info("сервер запущен")
	fmt.Println("сервер запущен")
	err1 := http.ListenAndServe("127.0.0.1:8080", r)
	logger.Info("сервер отключён")
	if err1 != nil {
		logger.Error("сервер на запустился")
		log.Fatal(err)
	}

}
