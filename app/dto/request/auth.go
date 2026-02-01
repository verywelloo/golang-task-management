package request

type RegisterPayload struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
