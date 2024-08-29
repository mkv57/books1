package main

import (
	"bookServer/internal/api"
	"bookServer/internal/db"
	"bookServer/internal/domain"

	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

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

	/*var h = api.Server{
		Di: http.Server{
			//Addr:           ":8080",
			//Handler:        r,
			//ReadTimeout:    10 * time.Second,
			//WriteTimeout:   10 * time.Second,
			//MaxHeaderBytes: 1 << 20,
		},
	}
	*/
	ourServer := api.Server{
		Database: db.Repository{
			Store: make(map[int]domain.Book),
		},
	}
	r.HandleFunc("/book", ourServer.GetBook).Methods(http.MethodGet)
	r.HandleFunc("/book", ourServer.AddBook).Methods(http.MethodPost)
	r.HandleFunc("/book", ourServer.DeleteBook).Methods(http.MethodDelete)
	r.HandleFunc("/book", ourServer.UpdateBook).Methods(http.MethodPut)
	r.HandleFunc("/books", ourServer.AllBooks).Methods(http.MethodGet)

	logger.Warn("сервер запущен")
	fmt.Println("сервер запущен")
	err = http.ListenAndServe("127.0.0.1:8080", r)
	logger1.Warn("сервер отключён")
	if err != nil {
		logger.Error("сервер нe запустился")
		log.Fatal(err)
	}

}
