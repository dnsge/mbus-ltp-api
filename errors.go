package mbus

import (
	"errors"
	"strings"
)

const (
	mNoArrivalTimes       = "No arrival times"
	mNoServiceScheduled   = "No service scheduled"
	mNoParameterDataFound = "No data found for parameter"
	mNoRTPIDataFeed       = "No RTPI Data Feed parameter provided"
	mInvalidDataFeed      = "Invalid RTPI Data Feed parameter"
	mUnsupportedDataFeed  = "The rtpidatafeed does not support this function"
)

var (
	ErrParameterNotFound       = errors.New("no data found for parameter")
	ErrNoRTPIDataFeed          = errors.New("no RTPI data feed parameter provided")
	ErrInvalidRTPIDataFeed     = errors.New("invalid RTPI data feed")
	ErrUnsupportedRTPIDataFeed = errors.New("call does not support RTPI data feed")
)

func errorMessageIs(message string, anyOf ...string) bool {
	for i := range anyOf {
		if anyOf[i] == message || strings.Index(message, anyOf[i]) != -1 {
			return true
		}
	}

	return false
}

var checkErrorArr = map[string]error{
	mNoParameterDataFound: ErrParameterNotFound,
	mNoRTPIDataFeed:       ErrNoRTPIDataFeed,
	mInvalidDataFeed:      ErrInvalidRTPIDataFeed,
	mUnsupportedDataFeed:  ErrUnsupportedRTPIDataFeed,
}

func determineError(message string) error {
	for text, err := range checkErrorArr {
		if errorMessageIs(message, text) {
			return err
		}
	}
	return nil
}
