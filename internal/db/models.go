package db

type Session struct {
	SessionId    string
	Open         bool
	MessageCount int64
	OpenedTime   string
	ClosedTime   string
	ForceClosed  bool
}

type IncomingMessage struct {
	MessageId             string
	MessageTypeId         int64
	MessageSequenceNumber int64
	SessionId             string
	DeviceId              string
}

type Device struct {
	DeviceId string
	Name     string
}

type DeviceLocation struct {
	DeviceId        string
	Timestamp       int64
	Latitude        float64
	Longitude       float64
	Altitude        float64
	Speed           float64
	HouseNumber     int64
	Road            string
	Neighbourhood   string
	Suburb          string
	City            string
	County          string
	State           string
	SubdivisionCode string
	Postcode        int64
	Country         string
	CountryCode     string
}
