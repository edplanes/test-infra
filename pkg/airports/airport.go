package airports

type AirportLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Elevation int     `json:"elevation"`
}

type Airport struct {
	ICAO     string          `json:"icao"`
	IATA     string          `json:"iata"`
	City     string          `json:"city"`
	Location AirportLocation `json:"location"`
	Name     string          `json:"name"`
	Score    AirportType     `json:"score"`
}

type AirportType int

const (
	Unknown AirportType = -1
	Closed  AirportType = iota
	BalloonPort
	HeliPort
	Seaplane
	Small
	Medium
	Large
)

func NewAirportType(data string) AirportType {
	switch data {
	case "closed":
		return Closed
	case "balloonport":
		return BalloonPort
	case "heliport":
		return HeliPort
	case "large_airport":
		return Large
	case "medium_airport":
		return Medium
	case "seaplane_base":
		return Seaplane
	case "small_airport":
		return Small
	default:
		return Unknown
	}
}
