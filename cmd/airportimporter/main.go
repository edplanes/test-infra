package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/edplanes/test-infra/pkg/airports"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var airportsUrl = "https://davidmegginson.github.io/ourairports-data/airports.csv"

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("Failed to connect to mongodb")
	}
	defer client.Disconnect(context.TODO())

	coll := client.Database("edplanes").Collection("airports")

	resp, err := http.Get(airportsUrl)
	if err != nil {
		log.Fatal("Failed to fetch airports database")
	}

	reader := csv.NewReader(resp.Body)

	data, err := reader.ReadAll()
	if err != nil {
		log.Fatal("Failed to load airports database")
	}

	airports := []airports.Airport{}
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

		airports = append(airports, airport)
	}

	newAirports := make([]interface{}, 0)
	for _, airport := range airports {
		newAirports = append(newAirports, airport)
	}

	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "icao", Value: "text"},
			{Key: "name", Value: "text"},
			{Key: "city", Value: "text"},
		},
	}

	err = coll.Drop(context.TODO())
	if err != nil {
		log.Fatal("Cannot drop current airports")
	}

	_, err = coll.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		log.Fatal("Cannot recreate indexes")
	}

	result, err := coll.InsertMany(context.TODO(), newAirports)
	if err != nil {
		log.Fatal("Cannot insert new airports")
	}

	fmt.Println(result, len(airports))
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
		ICAO: data[1],
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
