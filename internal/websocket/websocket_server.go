package websocket

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type WebsocketServer struct {
	clients map[*websocket.Conn]Client
}

type Client struct {
	SessionId    uuid.UUID
	MessageCount int64
}

var websocketServer WebsocketServer

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func StartServer() {
	websocketServer = WebsocketServer{
		make(map[*websocket.Conn]Client),
	}

	slog.Info("Websocket Server Initialized")

	http.HandleFunc("/tpv-update", handleTpvUpdate)
	http.ListenAndServe(":4120", nil)
}

func handleTpvUpdate(w http.ResponseWriter, r *http.Request) {
	connection, _ := upgrader.Upgrade(w, r, nil)

	websocketServer.clients[connection] = Client{SessionId: uuid.New(), MessageCount: 0}

	client := websocketServer.clients[connection]

	for {
		mt, message, err := connection.ReadMessage()

		if err != nil || mt == websocket.CloseMessage {
			break
		}

		client.IncrementMessageCount()

		go tpvUpdateMessageHandler(message)
	}

	sessionId := client.SessionId.String()
	messageCountStr := strconv.FormatInt(client.MessageCount, 10)

	slog.Info("Closing client connection with sessionId " + sessionId + " after " + messageCountStr + " messages")

	delete(websocketServer.clients, connection)

	connection.Close()
}

func tpvUpdateMessageHandler(message []byte) {
	slog.Info((string(message)))
}

func (c *Client) IncrementMessageCount() {
	c.MessageCount = c.MessageCount + 1
}
