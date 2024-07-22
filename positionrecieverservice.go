package main

import (
	"log/slog"

	"positionrecieverservice/internal/websocket"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelInfo)

	websocket.StartServer()
}
