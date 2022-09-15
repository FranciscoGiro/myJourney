package services

import (
    "fmt"
	"os"
	"errors"
    "io/ioutil"
    "net/http"
	"encoding/json"
)

var (
	unableToReadLocation = errors.New("Unable to read image location. Please insert one")
)

type GeoLocation struct {
	Items []struct {
		Address struct {
			CountryName string `json:"countryName"`
			City string `json:"city"`
		} `json:"address"`
		Title string `json:"title"`
	} `json:"items"`
}


func ReverseGeocode(lat, lng *float64) (string, string, error){
	apikey := os.Getenv("GEOCODING_KEY")

    url := fmt.Sprintf("https://revgeocode.search.hereapi.com/v1/revgeocode?apiKey=%s&at=%s,%s", apikey, *lat, *lng)

	res, err := http.Get(url)
    if err != nil {
		fmt.Println("Unable to access Reverse Geocode url. Error:", err)
		return "", "", unableToReadLocation
    }
    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
		fmt.Println("Unable to read response. Error:", err)
		return "", "", unableToReadLocation
    }

	var data GeoLocation
	json.Unmarshal(body, &data)

	var city, country string

	for _, item := range data.Items {
		city = item.Address.City
		country = item.Address.CountryName
	}


	return city, country, nil
}