package main

import (
	"fmt"
	"log/slog"
	"positionrecieverservice/internal/db"
	"positionrecieverservice/internal/websocket"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelInfo)

	db.InitDB("electrack", "localhost:3306", "user", "pass")

	_, forceCloseErr := db.ForceCloseAllSessions()

	if forceCloseErr != nil {
		slog.Error(fmt.Sprintf("Failed to Force Open Sessions Closed: %s", forceCloseErr.Error()))
	}

	websocket.StartServer()
}
