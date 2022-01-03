package mbus

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type predictionsResponse struct {
	Wrapper struct {
		Predictions []BusPrediction `json:"prd"`
	} `json:"bustime-response"`
}

type BusPrediction struct {
	StopName            string `json:"stpnm"`
	StopID              string `json:"stpid"`
	PredictionType      string `json:"typ"` // "A" for arrival and "D" for departure
	VehicleID           string `json:"vid"`
	RouteID             string `json:"rt"`
	RouteDisplay        string `json:"rtdd"`
	RouteDirection      string `json:"rtdir"`
	FinalDestination    string `json:"des"`
	PredictedArrival    string `json:"prdtm"` // Formatted as 20211226 13:22
	Delayed             bool   `json:"dly"`
	DynamicActionMode   int    `json:"dyn"`
	PredictionCountdown string `json:"prdctdn"` // Time in minutes until arrival
}

func GetStopPredictions(stopID string, routeIDs []string) ([]BusPrediction, error) {
	req, err := http.NewRequest("GET", predictionsAPIURL, nil)
	if err != nil {
		return nil, err
	}

	args := req.URL.Query()
	args.Set("stpid", stopID)
	args.Set("format", "json")
	args.Set("locale", "en")
	args.Set("tmres", "s") // Seconds time resolution
	if routeIDs != nil {
		args.Set("rt", strings.Join(routeIDs, ","))
	}
	req.URL.RawQuery = args.Encode()

	req.Header.Set("Accept", "application/json")
	prepareRequestWithV3Auth(req)

	res, err := doApiRequest(req)
	if err != nil {
		return nil, err
	}

	if ok, bErr, err := checkApiResponse(res); err != nil {
		return nil, err
	} else if !ok {
		// Bustime error occurred, take first error message
		errorMessage := bErr[0].Message
		if MessageIs(errorMessage, mNoArrivalTimes, mNoServiceScheduled) {
			// Handle errors that aren't really errors
			return []BusPrediction{}, nil
		} else if MessageIs(errorMessage, mNoParameterDataFound) {
			return nil, ErrParameterNotFound
		}
		return nil, fmt.Errorf("got bustime error: %s", errorMessage)
	}

	defer res.Body.Close()

	// Parse JSON response
	var data predictionsResponse
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, err
	}

	return data.Wrapper.Predictions, nil
}
