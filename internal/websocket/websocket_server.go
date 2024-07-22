package websocket

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
)

type WebsocketServer struct {
	clients map[*websocket.Conn]bool
}

var websocketServer WebsocketServer

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func StartServer() {
	websocketServer = WebsocketServer{
		make(map[*websocket.Conn]bool),
	}

	slog.Info("Websocket Server Initialized")

	http.HandleFunc("/tpv-update", handleTpvUpdate)
	http.ListenAndServe(":4120", nil)
}

func handleTpvUpdate(w http.ResponseWriter, r *http.Request) {
	connection, _ := upgrader.Upgrade(w, r, nil)

	websocketServer.clients[connection] = true

	for {
		mt, message, err := connection.ReadMessage()

		if err != nil || mt == websocket.CloseMessage {
			break
		}

		go tpvUpdateMessageHandler(message)
	}

	delete(websocketServer.clients, connection)

	connection.Close()

	slog.Info("Closed Client Connection")
}

func tpvUpdateMessageHandler(message []byte) {
	slog.Info((string(message)))
}
