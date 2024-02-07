package veikkausapi

const (
	RobotIdentifierHeaderKey   string = "X-ESA-API-KEY"
	RobotIdentifierHeaderValue string = "ROBOT"
	ContentType                string = "application/json"
	Accept                     string = "application/json"
	VeikkausApiBaseUrl         string = "https://www.veikkaus.fi/api/bff/v1"
	LoginEndpoint              string = "/sessions"
	Post                       string = "POST"
	Get                        string = "GET"
	Put                        string = "PUT"
	HttpOk                     int    = 200
	HttpMultipleChoices        int    = 300
)
