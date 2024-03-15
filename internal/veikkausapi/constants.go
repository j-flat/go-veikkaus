package veikkausapi

const (
	VeikkausAPIVersion         string = "v1"
	RobotIdentifierHeaderKey   string = "X-ESA-API-KEY"
	RobotIdentifierHeaderValue string = "ROBOT"
	UserAgent                  string = "goveikkaus-client"
	ContentType                string = "application/json"
	Accept                     string = "application/json"

	// Endpoint paths, there is some variance in the paths on Veikkaus API
	LoginEndpoint          string = "bff/v1/sessions"
	AccountBalanceEndpoint string = "v1/players/self/account"
)

// SessionTimeoutSeconds is half-hour as shown here: https://github.com/VeikkausOy/sport-games-robot/issues/160
var SessionTimeoutSeconds int = 1800
var BaseURL string = "https://www.veikkaus.fi/api/"
