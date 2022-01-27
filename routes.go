package mbus

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type routesResponse struct {
	Wrapper struct {
		Routes []Route `json:"routes"`
	} `json:"bustime-response"`
}

type Route struct {
	ID      string `json:"rt"`
	Name    string `json:"rtnm"`
	Color   string `json:"rtclr"`
	Display string `json:"rtdd"`
}

func (a *APIClient) GetRoutes() ([]Route, error) {
	req, err := http.NewRequest("GET", routesAPIURL, nil)
	if err != nil {
		return nil, err
	}

	args := req.URL.Query()
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
		return nil, fmt.Errorf("got bustime error: %s", errorMessage)
	}

	defer res.Body.Close()

	// Parse JSON response
	var data routesResponse
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, err
	}

	return data.Wrapper.Routes, nil
}
