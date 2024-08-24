package main

import (
	"bookServer/struct_p"
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

	minimalLevel := slog.LevelInfo
	// Создаем новый логгер, который будет писать в стандартный вывод
	options := &slog.HandlerOptions{
		Level: minimalLevel,
	}
	handler := slog.NewTextHandler(os.Stdout, options)
	logger := slog.New(handler) // логи отправляю в поток

	handler1 := slog.NewTextHandler(file, options)
	logger1 := slog.New(handler1) // логи отправляю в файл

	r.Use(Logging(logger))

	r.HandleFunc("/book", struct_p.GetBook).Methods(http.MethodGet)
	r.HandleFunc("/book", struct_p.AddBook).Methods(http.MethodPost)
	r.HandleFunc("/book", struct_p.DeleteBook).Methods(http.MethodDelete)
	r.HandleFunc("/book", struct_p.UpdateBook).Methods(http.MethodPut)
	r.HandleFunc("/books", struct_p.AllBooks).Methods(http.MethodGet)

	logger1.Warn("сервер запущен")
	fmt.Println("сервер запущен")
	err = http.ListenAndServe("127.0.0.1:8080", r)
	logger.Warn("сервер отключён")
	if err != nil {
		logger.Error("сервер нe запустился")
		log.Fatal(err)
	}

}

func Logging(logger *slog.Logger) mux.MiddlewareFunc {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info(r.Method)
			starter := time.Now()
			logger.Info("Request", slog.String("uri", r.RequestURI), slog.String("remote_addr", r.RemoteAddr))
			next.ServeHTTP(w, r)
			logger.Info("Finished", slog.String("duration", time.Since(starter).String()))
		})
	}
}
