package db

import (
	"database/sql"
)

var DB *sql.DB

func AddSession(sessionId string) (sql.Result, error) {
	result, err := DB.Exec("INSERT INTO SESSIONS (SESSION_ID) VALUES (?)", sessionId)

	return result, err
}

func IncrementSessionMessageCount(sessionId string) (sql.Result, error) {
	statement := `
		UPDATE SESSIONS S
		SET S.MESSAGE_COUNT = S.MESSAGE_COUNT + 1,
		S.LAST_MESSAGE_TIME = CURRENT_TIMESTAMP()
		WHERE S.SESSION_ID = ?;`

	result, err := DB.Exec(statement, sessionId)

	return result, err
}

func CloseSession(sessionId string) (sql.Result, error) {
	statement := `
		UPDATE SESSIONS S
		SET S.OPEN = FALSE,
		S.CLOSED_TIME = CURRENT_TIMESTAMP()
		WHERE S.SESSION_ID = ?;`

	result, err := DB.Exec(statement, sessionId)

	return result, err
}

func ForceCloseAllSessions() (sql.Result, error) {
	statement := `
		UPDATE SESSIONS S
		SET S.OPEN = FALSE,
		S.CLOSED_TIME = CURRENT_TIMESTAMP(),
		S.FORCE_CLOSED = TRUE
		WHERE S.OPEN = TRUE;`

	result, err := DB.Exec(statement)

	return result, err
}

func AddIncomingMessage(incomingMessage IncomingMessage) (sql.Result, error) {
	result, err := DB.Exec(
		"INSERT INTO INCOMING_MESSAGES (MESSAGE_ID, MESSAGE_TYPE_ID, MESSAGE_SEQUENCE_NUMBER, SESSION_ID, DEVICE_ID) VALUES (?, ?, ?, ?, ?)",
		incomingMessage.MessageId,
		incomingMessage.MessageTypeId,
		incomingMessage.MessageSequenceNumber,
		incomingMessage.SessionId,
		incomingMessage.DeviceId,
	)

	return result, err
}

func UpsertDevice(device Device) (sql.Result, error) {
	statement := `
		INSERT INTO DEVICES (DEVICE_ID, NAME, MOD_TIME)
		VALUES (?, ?, CURRENT_TIMESTAMP())
		ON DUPLICATE KEY UPDATE
		MOD_TIME = IF(NAME <> VALUES(NAME), CURRENT_TIMESTAMP(), MOD_TIME),
		NAME = IF(NAME <> VALUES(NAME), VALUES(NAME), NAME)`

	result, err := DB.Exec(statement, device.DeviceId, device.Name)

	return result, err
}

func AddDeviceLocation(deviceLocation DeviceLocation) (sql.Result, error) {
	statement := `INSERT INTO DEVICE_LOCATIONS (DEVICE_ID, TS,
		LATITUDE, LONGITUDE, ALTITUDE, SPEED,
		HOUSE, STREET, NEIGHBOURHOOD, CITY, COUNTY, 
		ZIPCODE, STATE, SUBDIVISION_CODE, COUNTRY)
		VALUES (?, FROM_UNIXTIME(?), ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := DB.Exec(
		statement,
		deviceLocation.DeviceId,
		deviceLocation.Timestamp,
		deviceLocation.Latitude,
		deviceLocation.Longitude,
		deviceLocation.Altitude,
		deviceLocation.Speed,
		deviceLocation.HouseNumber,
		deviceLocation.Road,
		deviceLocation.Neighbourhood,
		deviceLocation.City,
		deviceLocation.County,
		deviceLocation.Postcode,
		deviceLocation.State,
		deviceLocation.SubdivisionCode,
		deviceLocation.Country,
	)

	return result, err
}
