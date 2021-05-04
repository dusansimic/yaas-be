package yaas

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	// ErrGeoIPPrivateRange means that event has been sent from inside a reserved range
	// https://en.wikipedia.org/wiki/Private_network#Private_IPv6_addresses
	ErrGeoIPPrivateRange = errors.New("geo ip: specified ip is inside a private range")
	// ErrGeoIPReservedRange means that event has been sent from inside a reserved range
	// https://en.wikipedia.org/wiki/Reserved_IP_addresses
	ErrGeoIPReservedRange = errors.New("geo ip: specified ip is inside a reserved range")
	// ErrGeoIPInvalidQuery means that query sent to the api was invalid
	ErrGeoIPInvalidQuery = errors.New("geo ip: invalid query")
	// ErrGeoIPUnknown means that the error that has occured is unknown
	ErrGeoIPUnknown = errors.New("geo ip: an unknown error has occured")
)

// GeoIP stores infromation about users ip location
type GeoIP struct {
	Status      string  `json:"status"`
	Message     string  `json:"message"`
	Continent   string  `json:"continent"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float32 `json:"lat"`
	Lon         float32 `json:"lon"`
	TZ          string  `json:"timezone"`
}

// GetIPLocation finds the ip location and returns it.
func GetIPLocation(ip string) (GeoIP, error) {
	r, err := http.Get(fmt.Sprintf("http://ip-api.com/json/%s", ip))
	if err != nil {
		return GeoIP{}, err
	}
	defer r.Body.Close()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return GeoIP{}, err
	}

	var g GeoIP
	if err := json.Unmarshal(b, &g); err != nil {
		return GeoIP{}, err
	}

	// ! Uncomment this
	// if g.Status == "fail" {
	// 	switch g.Message {
	// 	case "private range":
	// 		return GeoIP{}, ErrGeoIPPrivateRange
	// 	case "reserved range":
	// 		return GeoIP{}, ErrGeoIPReservedRange
	// 	case "invalid query":
	// 		return GeoIP{}, ErrGeoIPInvalidQuery
	// 	default:
	// 		return GeoIP{}, ErrGeoIPUnknown
	// 	}
	// }

	return g, nil
}
