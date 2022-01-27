package mbus

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type rtpiResponse struct {
	Wrapper struct {
		Feeds []RTPIDataFeed `json:"rtpidatafeeds"`
	} `json:"bustime-response"`
}

type strBool bool

func (strB *strBool) UnmarshalJSON(data []byte) error {
	asStr := strings.ToLower(strings.Trim(string(data), `"`))
	if asStr == "true" {
		*strB = true
	} else if asStr == "false" {
		*strB = false
	} else {
		return fmt.Errorf("strBool unmarshal: invalid input %s", asStr)
	}
	return nil
}

type RTPIDataFeed struct {
	Name        string  `json:"name"`
	Source      string  `json:"source"`
	DisplayName string  `json:"displayname"`
	Enabled     strBool `json:"enabled"`
	Visible     strBool `json:"visible"`
}

func (a *APIClient) GetRTPIDataFeeds() ([]RTPIDataFeed, error) {
	req, err := http.NewRequest("GET", rtpiURL, nil)
	if err != nil {
		return nil, err
	}

	args := req.URL.Query()
	args.Set("format", "json")
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
	var data rtpiResponse
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, err
	}

	return data.Wrapper.Feeds, nil
}

// AutoConfigureDataFeed queries the available RTPIDataFeeds and selects
// the first one with the enabled field.
func (a *APIClient) AutoConfigureDataFeed() error {
	allFeeds, err := a.GetRTPIDataFeeds()
	if err != nil {
		return err
	}

	for _, feed := range allFeeds {
		if feed.Enabled {
			a.DataFeed = feed.Name
			return nil
		}
	}

	// No enabled data feeds
	a.DataFeed = ""
	return nil
}
