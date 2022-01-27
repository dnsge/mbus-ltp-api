package mbus

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type stopsResponse struct {
	Wrapper struct {
		Stops []Stop `json:"stops"`
	} `json:"bustime-response"`
}

type Stop struct {
	ID        string  `json:"stpid"`
	Name      string  `json:"stpnm"`
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
}

func (a *APIClient) GetStops(routeID string) ([]Stop, error) {
	directions, err := a.GetDirections(routeID)
	if err != nil {
		return nil, err
	}

	var allStops []Stop
	for _, direction := range directions {
		stops, err := a.GetStopsInDirection(routeID, direction.ID)
		if err != nil {
			return nil, err
		}

		allStops = append(allStops, stops...)
	}

	return DeduplicateStops(allStops), nil
}

func (a *APIClient) GetStopsInDirection(routeID string, directionID string) ([]Stop, error) {
	req, err := http.NewRequest("GET", stopsAPIURL, nil)
	if err != nil {
		return nil, err
	}

	args := req.URL.Query()
	args.Set("rt", routeID)
	args.Set("dir", directionID)
	if a.DataFeed != "" {
		args.Set("rtpidatafeed", a.DataFeed)
	}
	args.Set("format", "json")
	args.Set("locale", "en")
	req.URL.RawQuery = args.Encode()

	req.Header.Set("Accept", "application/json")

	res, err := a.doApiRequest(req)
	if err != nil {
		return nil, err
	}

	if ok, bErr, err := checkApiResponse(res); err != nil {
		return nil, err
	} else if !ok {
		// Bustime error occurred, take first error message
		errorMessage := bErr[0].Message
		if det := determineError(errorMessage); det != nil {
			return nil, det
		}
		return nil, fmt.Errorf("got bustime error: %s", errorMessage)
	}

	defer res.Body.Close()

	// Parse JSON response
	var data stopsResponse
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, err
	}

	return data.Wrapper.Stops, nil
}

type directionsResponse struct {
	Wrapper struct {
		Directions []Direction `json:"directions"`
	} `json:"bustime-response"`
}

type Direction struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (a *APIClient) GetDirections(routeID string) ([]Direction, error) {
	req, err := http.NewRequest("GET", directionsAPIURL, nil)
	if err != nil {
		return nil, err
	}

	args := req.URL.Query()
	args.Set("rt", routeID)
	if a.DataFeed != "" {
		args.Set("rtpidatafeed", a.DataFeed)
	}
	args.Set("format", "json")
	args.Set("locale", "en")
	req.URL.RawQuery = args.Encode()

	req.Header.Set("Accept", "application/json")

	res, err := a.doApiRequest(req)
	if err != nil {
		return nil, err
	}

	if ok, bErr, err := checkApiResponse(res); err != nil {
		return nil, err
	} else if !ok {
		// Bustime error occurred, take first error message
		errorMessage := bErr[0].Message
		if det := determineError(errorMessage); det != nil {
			return nil, det
		}
		return nil, fmt.Errorf("got bustime error: %s", errorMessage)
	}

	defer res.Body.Close()

	// Parse JSON response
	var data directionsResponse
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, err
	}

	return data.Wrapper.Directions, nil
}

// DeduplicateStops removes duplicate stops from the input slice by checking stop IDs.
func DeduplicateStops(stops []Stop) (ret []Stop) {
	stopsByID := make(map[string]*Stop, len(stops))

	for i := range stops {
		s := &stops[i]
		if _, ok := stopsByID[s.ID]; !ok {
			stopsByID[s.ID] = s
		}
	}

	// Preallocate space for output
	ret = make([]Stop, len(stopsByID))
	i := 0
	for _, val := range stopsByID {
		ret[i] = *val
		i++
	}

	return ret
}
