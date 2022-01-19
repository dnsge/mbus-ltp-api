package mbus

import "errors"

const (
	mNoArrivalTimes       = "No arrival times"
	mNoServiceScheduled   = "No service scheduled"
	mNoParameterDataFound = "No data found for parameter"
)

var (
	ErrParameterNotFound = errors.New("no data found for parameter")
)

func MessageIs(message string, anyOf ...string) bool {
	for i := range anyOf {
		if anyOf[i] == message {
			return true
		}
	}

	return false
}
