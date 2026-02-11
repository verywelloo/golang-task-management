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

type CacheSession struct {
	ID    string `json:"_id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Ip    string `json:"ip"`
	Agent string `json:"agent"`
}
