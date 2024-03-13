package veikkausapi

const (
	VeikkausAPIVersion         string = "v1"
	RobotIdentifierHeaderKey   string = "X-ESA-API-KEY"
	RobotIdentifierHeaderValue string = "ROBOT"
	UserAgent                  string = "goveikkaus-client"
	ContentType                string = "application/json"
	Accept                     string = "application/json"
	// BaseURL                    string = "https://www.veikkaus.fi"
	LoginEndpoint          string = "bff/v1/sessions"
	AccountBalanceEndpoint string = "v1/players/self/account"
)

var SessionTimeoutSeconds int = 1800
var BaseURL string = "https://www.veikkaus.fi/api/"
