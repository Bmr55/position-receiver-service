package websocket

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"positionrecieverservice/internal/common"
	"positionrecieverservice/internal/db"
	"positionrecieverservice/internal/geocoder"
)

const TPVUpdateMessageTypeId = 0
const NominatimGeocoderBaseUrl = "http://10.0.0.129/nominatim/"

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
	sessionId := client.SessionId.String()

	db.AddSession(sessionId)

	slog.Info(fmt.Sprintf("New Connection with SessionId: %s", sessionId))

	for {
		mt, message, err := connection.ReadMessage()

		if err != nil || mt == websocket.CloseMessage {
			_, closeErr := db.CloseSession(sessionId)

			if closeErr != nil {
				slog.Error(fmt.Sprintf("Failed to mark session '%s' closed in database: %s", sessionId, closeErr.Error()))
			}

			break
		}

		client.IncrementMessageCount()

		_, incrementErr := db.IncrementSessionMessageCount(sessionId)

		if incrementErr != nil {
			slog.Error(fmt.Sprintf("Failed to increment session '%s' message count in database: %s", sessionId, incrementErr.Error()))
		}

		slog.Info(fmt.Sprintf("Received TPVUpdate message for session: '%s'", sessionId))

		go tpvUpdateMessageHandler(client, message)
	}

	messageCountStr := strconv.FormatInt(client.MessageCount, 10)

	slog.Info(fmt.Sprintf("Closing client connection with sessionId '%s' after %s messages", sessionId, messageCountStr))

	delete(websocketServer.clients, connection)

	connection.Close()
}

func (c *Client) IncrementMessageCount() {
	c.MessageCount = c.MessageCount + 1
}

func tpvUpdateMessageHandler(client Client, message []byte) {
	var tpvUpdateMessage common.TPVUpdateMessage

	unmarshalErr := json.Unmarshal(message, &tpvUpdateMessage)

	if unmarshalErr != nil {
		slog.Error(fmt.Sprintf("Failed to Unmarshal TPVUpdateMessage: %s", unmarshalErr.Error()))

		return
	}

	device := db.Device{DeviceId: tpvUpdateMessage.Device.Id, Name: tpvUpdateMessage.Device.Name}

	_, deviceErr := db.UpsertDevice(device)

	if deviceErr != nil {
		slog.Error(fmt.Sprintf("Device '%s' Upsert Failed: %s", device.DeviceId, deviceErr.Error()))

		return
	}

	messageSequenceNumber, seqErr := tpvUpdateMessage.MessageSequenceNumber.Int64()

	if seqErr != nil {
		slog.Warn(fmt.Sprintf("Failed to convert MessageSequenceNumber to Int64: %s", seqErr.Error()))

		messageSequenceNumber = -1
	}

	incomingMessage := db.IncomingMessage{
		MessageId:             uuid.New().String(),
		MessageTypeId:         TPVUpdateMessageTypeId,
		MessageSequenceNumber: messageSequenceNumber,
		SessionId:             client.SessionId.String(),
		DeviceId:              tpvUpdateMessage.Device.Id,
	}

	_, incomingMessageErr := db.AddIncomingMessage(incomingMessage)

	if incomingMessageErr != nil {
		slog.Error(fmt.Sprintf("Incoming Message Insert Failed: %s", incomingMessageErr.Error()))
	}

	geocoder := geocoder.NominatimGeocoder{BaseUrl: NominatimGeocoderBaseUrl}
	tpv := common.CreateFlatTPV(tpvUpdateMessage.TPV)
	resp, geocodeErr := geocoder.ReverseGeocode(tpv.Latitude, tpv.Longitude)

	var deviceLocation db.DeviceLocation

	if geocodeErr != nil {
		slog.Error(fmt.Sprintf("Failed to reverse geocode (lat=%s, lon=%s): %s",
			tpvUpdateMessage.TPV.Position.Latitude.String(), tpvUpdateMessage.TPV.Position.Longitude.String(), geocodeErr.Error()))

		deviceLocation = db.CreateDeviceLocationModelWithUnknownAddress(device.DeviceId, tpv)
	} else {
		deviceLocation = db.CreateDeviceLocationModel(device.DeviceId, tpv, resp.Address)
	}

	_, deviceLocationErr := db.AddDeviceLocation(deviceLocation)

	if deviceLocationErr != nil {
		slog.Error(fmt.Sprintf("Device Location Insert Failed: %s", deviceLocationErr.Error()))
	}

	latStr := strconv.FormatFloat(tpv.Latitude, 'f', -1, 64)
	lonStr := strconv.FormatFloat(tpv.Longitude, 'f', -1, 64)

	slog.Info(fmt.Sprintf("Persisted Location (%s, %s) of Device with Id %s", latStr, lonStr, device.DeviceId))
}
