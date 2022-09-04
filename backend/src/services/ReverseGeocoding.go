package services

import (
    "fmt"
	"os"
    "io/ioutil"
    "log"
    "net/http"
	"encoding/json"
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

    url := "https://revgeocode.search.hereapi.com/v1/revgeocode?apiKey=" + apikey + "&at=" + fmt.Sprint(*lat) + "," + fmt.Sprint(*lng)

    res, err := http.Get(url)
    if err != nil {
        fmt.Println("Couldn't get response from ReverseGeocode server:, err")
		return "", "", err
    }
    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        log.Fatalln(err)
    }

	var data GeoLocation
	json.Unmarshal(body, &data)

	var city, country string

	for _, item := range data.Items {
		city = item.Address.City
		country = item.Address.CountryName
		fmt.Println("CITY:",item.Address.City)
		fmt.Println("COUNTRY:",item.Address.CountryName)
	}


	return city, country, nil
}