package goveikkaus

type SessionCookies struct {
	JSessionID string `json:"JSESSIONID"`
}

type LoginPayload struct {
	Type     string `json:"type"`
	User     string `json:"login"`
	Password string `json:"password"`
}

type LoginSuccessful = map[string]interface{}

type AuthService service
