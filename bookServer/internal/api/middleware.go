package api

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

var minimalLevel = slog.LevelInfo // другого метода пока не нашёл, чтобы работали логеры, которые на в "main"

var Logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{ // пришлось дублировать файлы из "main"
	Level: minimalLevel,
}))

func Logging(Logger *slog.Logger) mux.MiddlewareFunc {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			//b, _  := copy(r.Method)
			starter := time.Now()

			Logger.Info("Request", slog.String("uri", r.RequestURI), slog.String("remote_addr", r.RemoteAddr))

			switch r.Method {
			case "GET":
				Logger.Info("отвечаем на запрос")
			case "POST":
				Logger.Info("создаём")
			case "PUT":
				Logger.Info("обновляем")
			case "DELETE":
				Logger.Info("удаляем")
			}

			Logger.Info("Finished", slog.String("duration", time.Since(starter).String()))
			//r = b
			next.ServeHTTP(w, r)

		})
	}
}
