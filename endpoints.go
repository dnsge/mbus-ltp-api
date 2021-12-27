package mbus

const (
	v2RestrictedURL = "https://mbus.ltp.umich.edu/bustime/api/restricted/v2"
	v3BaseURL       = "https://mbus.ltp.umich.edu/bustime/api/v3"

	stopsAPIURL       = v3BaseURL + "/getstops"
	directionsAPIURL  = v3BaseURL + "/getdirections"
	predictionsAPIURL = v3BaseURL + "/getpredictions"
	routesAPIURL      = v3BaseURL + "/getroutes"
	nearbyAPIURL      = v2RestrictedURL + "/stops/nearby"
)
