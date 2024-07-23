package geocoder

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

type Geocoder interface {
	ReverseGeocode(latitude float64, longitude float64) (ReverseGeocodeResponse, error)
}

type NominatimGeocoder struct {
	BaseUrl string
}

type Address struct {
	HouseNumber        string `json:"house_number"`
	Road               string `json:"road"`
	Neighbourhood      string `json:"neighbourhood"`
	Suburb             string `json:"suburb"`
	City               string `json:"city"`
	County             string `json:"county"`
	State              string `json:"state"`
	ISO3166CountryCode string `json:"ISO3166-2-lvl4"`
	Postcode           string `json:"postcode"`
	Country            string `json:"country"`
	CountryCode        string `json:"country_code"`
}

type ReverseGeocodeResponse struct {
	PlaceId     json.Number `json:"place_id"`
	Licence     string      `json:"licence"`
	OsmType     string      `json:"osm_type"`
	OsmId       json.Number `json:"osm_id"`
	Lat         json.Number `json:"lat"`
	Lon         json.Number `json:"lon"`
	Category    string      `json:"category"`
	PlaceType   string      `json:"type"`
	PlaceRank   json.Number `json:"place_rank"`
	Importance  json.Number `json:"house_number"`
	AddressType string      `json:"addresstype"`
	Name        string      `json:"name"`
	DisplayName string      `json:"display_name"`
	Address     Address     `json:"address"`
	BoundingBox []string    `json:"boundingbox"`
}

func (n NominatimGeocoder) ReverseGeocode(latitude float64, longitude float64) (ReverseGeocodeResponse, error) {
	latStr := strconv.FormatFloat(latitude, 'f', -1, 64)
	lonStr := strconv.FormatFloat(longitude, 'f', -1, 64)
	url := n.BaseUrl + "reverse?lat=" + latStr + "&lon=" + lonStr + "&format=jsonv2"
	byteResp, reqErr := getRequest(url)

	if reqErr != nil {
		return ReverseGeocodeResponse{}, reqErr
	}

	var reverseGeocodeResp ReverseGeocodeResponse

	json.Unmarshal(byteResp, &reverseGeocodeResp)

	return reverseGeocodeResp, nil
}

func getRequest(url string) ([]byte, error) {
	resp, getErr := http.Get(url)

	if getErr != nil {
		return nil, getErr
	}

	body, readErr := io.ReadAll(resp.Body)

	if readErr != nil {
		return nil, readErr
	}

	return body, nil
}
