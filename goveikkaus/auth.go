package goveikkaus

// Service type: Auth
type AuthService service

// Request Payload Types for AuthService Endpoints
type LoginPayload struct {
	Type     string `json:"type"`
	User     string `json:"login"`
	Password string `json:"password"`
}

// End of Request payload types

// Response Types for AuthService Endpoints
type LoginSuccessful = map[string]interface{}

// Account Balance types
type Balances struct {
	Cash Cash `json:"CASH"`
}
type Cash struct {
	Currency      string `json:"currency"`
	Type          string `json:"type"`
	Balance       int    `json:"balance"`
	UsableBalance int    `json:"usableBalance"`
	FrozenBalance int    `json:"frozenBalance"`
	HoldBalance   int    `json:"holdBalance"`
}
type AccountBalance struct {
	Status        string   `json:"status"`
	TimerInterval int      `json:"timerInterval"`
	Balances      Balances `json:"balances"`
}

// End of Response Types for AuthService Endpoints
