package response

type Profile struct {
	SessionID string `json:"session_id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
}

type LoginRes struct {
	Profile Profile `json:"profile"`
	Type    string  `json:"type"`
	Token   string  `json:"token"`
}
