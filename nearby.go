package mbus

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type nearbyStopsResponse struct {
	Data struct {
		Stops []NearbyStop `json:"stops"`
	} `json:"data"`
}

type NearbyStop struct {
	Distance float64 `json:"dist"`
	ID       string  `json:"stpid"`
}

func GetNearbyBusStops(lat string, lon string) ([]NearbyStop, error) {
	req, err := http.NewRequest("GET", nearbyAPIURL, nil)
	if err != nil {
		return nil, err
	}

	args := req.URL.Query()
	args.Set("lat", lat)
	args.Set("lon", lon)
	args.Set("rad", "0.5")
	args.Set("max", "10")
	req.URL.RawQuery = args.Encode()

	req.Header.Set("Accept", "application/json")
	prepareRequestWithV2Auth(req)

	res, err := doApiRequest(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Check response status code
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("bad status code: %d", res.StatusCode)
	}

	// Parse JSON response
	var data nearbyStopsResponse
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, err
	}

	return data.Data.Stops, nil
}

func ClosestStop(stops []NearbyStop, threshold float64) (*NearbyStop, error) {
	if len(stops) == 0 {
		return nil, ErrNoCloseStop
	}

	minDist := stops[0].Distance
	var closestStop *NearbyStop

	for i := range stops {
		stop := &stops[i]
		if stop.Distance <= minDist && stop.Distance <= threshold {
			minDist = stop.Distance
			closestStop = stop
		}
	}

	if closestStop == nil {
		return nil, ErrNoCloseStop
	}

	return closestStop, nil
}
