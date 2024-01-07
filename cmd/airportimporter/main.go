package main

import (
	"bytes"
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/edplanes/test-infra/pkg/airports"
)

var airportsUrl = "https://davidmegginson.github.io/ourairports-data/airports.csv"

func main() {

	resp, err := http.Get(airportsUrl)
	if err != nil {
		log.Fatal("Failed to fetch airports database")
	}

	reader := csv.NewReader(resp.Body)

	data, err := reader.ReadAll()
	if err != nil {
		log.Fatal("Failed to load airports database")
	}

	airportsData := []airports.Airport{}
	for i := 0; i < len(data); i++ {
		if i == 0 {
			err := verifyDataStructure(data[i])
			if err != nil {
				log.Fatal("Data structure has been changed!", err)
			}

			continue
		}
		fixedData := fixData(data[i])
		airport, err := parseAirportData(fixedData)
		if err != nil {
			log.Fatal("Failed to parse data for airport", err, data[i])
		}

		if !shouldAddAirport(airport) {
			log.Print("Skipping airport", airport)
			continue
		}

		airportsData = append(airportsData, airport)
	}

	newAirports := make([]interface{}, 0)
	for _, airport := range airportsData {
		newAirports = append(newAirports, airport)
	}

	jsonAirport, err := json.Marshal(newAirports)
	if err != nil {
		log.Fatal("Failed to marshal airports to json", err)
	}

	fmt.Print(string(jsonAirport))
	log.Printf("Found %d airports, importing to system...", len(newAirports))

	token, err := authenticateAgainstSystem("admin@localhost.com", "admin")
	if err != nil {
		log.Fatal("Failed to authenticate", err)
	}

	req, err := http.NewRequest("POST", "http://localhost:8080/api/airports/import", bytes.NewReader(jsonAirport))
	if err != nil {
		log.Fatal("Failed to create import request", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal("Failed to import", err)
	}

	response, _ := io.ReadAll(resp.Body)

	log.Print(string(response))

	log.Println("Done, imported ", len(newAirports), " airports")
}

func authenticateAgainstSystem(username, password string) (string, error) {
	type authInfo struct {
		Token string `json:"token"`
	}
	data := fmt.Sprintf("%s:%s", username, password)
	encoded := base64.StdEncoding.EncodeToString([]byte(data))

	req, err := http.NewRequest("GET", "http://localhost:8080/api/auth", nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", encoded))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	var authInfoData authInfo
	err = json.NewDecoder(res.Body).Decode(&authInfoData)
	if err != nil {
		return "", err
	}

	fmt.Print(authInfoData)

	return authInfoData.Token, nil
}

func shouldAddAirport(airport airports.Airport) bool {
	r, err := regexp.Compile("[a-zA-Z]{4}")
	if err != nil {
		return false
	}
	return airport.Score != airports.Closed &&
		airport.Score != airports.BalloonPort &&
		airport.Score != airports.HeliPort &&
		airport.Score != airports.Seaplane &&
		r.MatchString(airport.ICAO)
}

func verifyDataStructure(data []string) error {
	if data[1] != "ident" {
		return fmt.Errorf("ident position changed")
	}

	if data[2] != "type" {
		return fmt.Errorf("type position changed")
	}

	if data[3] != "name" {
		return fmt.Errorf("name position changed")
	}

	if data[4] != "latitude_deg" || data[5] != "longitude_deg" {
		return fmt.Errorf("locaiton data changed")
	}

	if data[6] != "elevation_ft" {
		return fmt.Errorf("elevation position changed")
	}

	if data[10] != "municipality" {
		return fmt.Errorf("city name position changed")
	}

	if data[13] != "iata_code" {
		return fmt.Errorf("iata code position changed")
	}

	return nil
}

func fixData(data []string) []string {
	if data[6] == "" {
		data[6] = "1"
	}

	return data
}

func parseAirportData(data []string) (airports.Airport, error) {
	latitude, err := strconv.ParseFloat(data[4], 64)
	if err != nil {
		return airports.Airport{}, err
	}
	longitude, err := strconv.ParseFloat(data[4], 64)
	if err != nil {
		return airports.Airport{}, err
	}
	elevation, err := strconv.ParseInt(data[6], 10, 32)
	if err != nil {
		return airports.Airport{}, err
	}

	return airports.Airport{
		ICAO: strings.ToUpper(data[1]),
		IATA: strings.ToUpper(data[13]),
		City: data[10],
		Location: airports.AirportLocation{
			Latitude:  latitude,
			Longitude: longitude,
			Elevation: int(elevation),
		},
		Name:  data[3],
		Score: airports.NewAirportType(data[2]),
	}, nil
}
