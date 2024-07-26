package db

import (
	"database/sql"
	"log"
	"log/slog"
	"positionrecieverservice/internal/common"
	"positionrecieverservice/internal/geocoder"
	"strconv"

	"github.com/go-sql-driver/mysql"
)

func InitDB(dbName string, address string, user string, password string) {
	cfg := getConfig(dbName, address, user, password)

	var openErr error

	DB, openErr = sql.Open("mysql", cfg.FormatDSN())

	if openErr != nil {
		log.Fatalf("Error Initializing Database: %v", openErr)
	}

	slog.Info("Database Initialized")
}

func getConfig(dbName string, address string, user string, password string) mysql.Config {
	cfg := mysql.Config{
		User:   user,
		Passwd: password,
		Net:    "tcp",
		Addr:   address,
		DBName: dbName,
	}

	return cfg
}

func CreateDeviceLocationModel(deviceId string, tpv common.FlatTPV, address geocoder.Address) DeviceLocation {
	houseNumber, _ := strconv.ParseInt(address.HouseNumber, 10, 64)
	postcode, _ := strconv.ParseInt(address.Postcode, 10, 64)

	return DeviceLocation{
		DeviceId:        deviceId,
		Timestamp:       tpv.Timestamp,
		Latitude:        tpv.Latitude,
		Longitude:       tpv.Longitude,
		Altitude:        tpv.Altitude,
		Speed:           tpv.Speed,
		HouseNumber:     houseNumber,
		Road:            address.Road,
		Neighbourhood:   address.Neighbourhood,
		Suburb:          address.Suburb,
		City:            address.City,
		County:          address.County,
		State:           address.State,
		SubdivisionCode: address.ISO3166CountryCode,
		Postcode:        postcode,
		Country:         address.Country,
		CountryCode:     address.CountryCode,
	}
}

func CreateDeviceLocationModelWithUnknownAddress(deviceId string, tpv common.FlatTPV) DeviceLocation {
	return DeviceLocation{
		DeviceId:        deviceId,
		Timestamp:       tpv.Timestamp,
		Latitude:        tpv.Latitude,
		Longitude:       tpv.Longitude,
		Altitude:        tpv.Altitude,
		Speed:           tpv.Speed,
		HouseNumber:     -1,
		Road:            "UNKOWWN",
		Neighbourhood:   "UNKOWWN",
		Suburb:          "UNKOWWN",
		City:            "UNKOWWN",
		County:          "UNKOWWN",
		State:           "UNKOWWN",
		SubdivisionCode: "UNKOWWN",
		Postcode:        -1,
		Country:         "UNKOWWN",
		CountryCode:     "UNKOWWN",
	}
}
