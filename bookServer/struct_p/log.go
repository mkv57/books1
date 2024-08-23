package struct_p

import (
	"log/slog"
	"os"
)

var minimalLevel = slog.LevelInfo

var file, _ = os.OpenFile("app.log", os.O_APPEND, 0666)

var logger = slog.New(slog.NewTextHandler(file, &slog.HandlerOptions{
	Level: minimalLevel,
}))
