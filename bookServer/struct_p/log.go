package struct_p

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

var minimalLevel = slog.LevelInfo // другого метода пока не нашёл, чтобы работали логеры, которые на в "main"

var logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{ // пришлось дублировать файлы из "main"
	Level: minimalLevel,
}))

func Logging(logger *slog.Logger) mux.MiddlewareFunc {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			starter := time.Now()
			logger.Info("Request", slog.String("uri", r.RequestURI), slog.String("remote_addr", r.RemoteAddr))
			next.ServeHTTP(w, r)
			n := r.ContentLength
			f := r.Method
			Body(n, f) // чтобы не делать ф-цию большой, создал и вызываю новую функцию
			logger.Info("Finished", slog.String("duration", time.Since(starter).String()))

		})
	}
}

func Body(c int64, r string) {

	switch r {
	case "GET":
		logger.Info("отвечаем на запрос")
	case "POST":
		logger.Info("создаём")
	case "PUT":
		logger.Info("обновляем")
	case "DELETE":
		logger.Info("удаляем")
	}

	if c == 0 {
		logger.Info("тело запроса отсутствует")
	}
	if c != 0 {
		logger.Info("тело запроса есть")
	}
}
