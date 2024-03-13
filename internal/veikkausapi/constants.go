package veikkausapi

const (
	VeikkausAPIVersion         string = "v1"
	RobotIdentifierHeaderKey   string = "X-ESA-API-KEY"
	RobotIdentifierHeaderValue string = "ROBOT"
	UserAgent                  string = "goveikkaus-client"
	ContentType                string = "application/json"
	Accept                     string = "application/json"
	// BaseURL                    string = "https://www.veikkaus.fi"
	VeikkausAPIBaseURL string = "api/bff/"
	LoginEndpoint      string = "sessions"
)

var SessionTimeoutSeconds int = 1800
var BaseURL string = "https://www.veikkaus.fi"
var OverWriteBaseURL bool = false
