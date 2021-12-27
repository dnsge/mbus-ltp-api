package mbus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.82 Safari/537.36"
)

var httpClient = &http.Client{
	Transport:     nil,
	CheckRedirect: nil,
	Jar:           nil,
	Timeout:       time.Second * 5,
}

type BustimeError struct {
	Wrapper struct {
		Error []BustimeErrorMessage `json:"error"`
	} `json:"bustime-response"`
}

type BustimeErrorMessage struct {
	Message string `json:"msg"`
}

func doApiRequest(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", userAgent)
	return httpClient.Do(req)
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
