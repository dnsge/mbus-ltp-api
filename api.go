package mbus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

var DefaultClient = &http.Client{
	Timeout: time.Second * 5,
}

type Config struct {
	// Client is the http.Client to use for requests. If nil, DefaultClient is used.
	Client *http.Client

	// Auth is an AuthApplier that prepares requests with the required authorization.
	Auth AuthApplier

	// UserAgent is an optional User-Agent header to set on requests.
	UserAgent string
}

type Client interface {
	GetStops(routeID string) ([]Stop, error)
	GetStopsInDirection(routeID string, directionID string) ([]Stop, error)
	GetDirections(routeID string) ([]Direction, error)
	GetRoutes() ([]Route, error)
	GetStopPredictions(stopID string, routeIDs []string) ([]BusPrediction, error)
}

var _ Client = &APIClient{}

type APIClient struct {
	client    *http.Client
	auth      AuthApplier
	userAgent string
}

// NewAPIClient returns an APIClient instance initialized with the provided Config.
func NewAPIClient(config *Config) *APIClient {
	client := config.Client
	if client == nil {
		client = DefaultClient
	}

	return &APIClient{
		client:    client,
		auth:      config.Auth,
		userAgent: config.UserAgent,
	}
}

type BustimeError struct {
	Wrapper struct {
		Error []BustimeErrorMessage `json:"error"`
	} `json:"bustime-response"`
}

type BustimeErrorMessage struct {
	Message string `json:"msg"`
}

// doApiRequest prepares a request with the User-Agent header and authorization
// before executing the request.
func (a *APIClient) doApiRequest(req *http.Request) (*http.Response, error) {
	if a.userAgent != "" {
		req.Header.Set("User-Agent", a.userAgent)
	}
	err := a.auth.ApplyAuth(req)
	if err != nil {
		return nil, err
	}
	return a.client.Do(req)
}

// checkApiResponse checks if a bustime-response error occurred.
// Returns whether the API call is OK, the list of error messages if not OK,
// or an error that occurred while decoding the response.
func checkApiResponse(res *http.Response) (bool, []BustimeErrorMessage, error) {
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return false, nil, fmt.Errorf("bad status code: %d", res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return false, nil, err
	}

	// Close original body and replace with nop-closer copy
	_ = res.Body.Close()
	res.Body = io.NopCloser(bytes.NewBuffer(data))

	var bErr BustimeError
	if err := json.Unmarshal(data, &bErr); err != nil {
		return false, nil, err
	}

	if len(bErr.Wrapper.Error) != 0 {
		return false, bErr.Wrapper.Error, nil
	}

	return true, nil, nil
}
