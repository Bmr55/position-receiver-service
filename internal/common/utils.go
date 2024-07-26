package common

func CreateFlatTPV(tpv TPV) FlatTPV {
	ts, _ := tpv.Timestamp.Int64()
	lat, _ := tpv.Position.Latitude.Float64()
	lon, _ := tpv.Position.Longitude.Float64()
	alt, _ := tpv.Position.Altitude.Float64()
	speed, _ := tpv.Speed.Float64()

	return FlatTPV{
		Timestamp: ts,
		Latitude:  lat,
		Longitude: lon,
		Altitude:  alt,
		Speed:     speed,
	}
}
