package request

type RegisterPayload struct {
	Id          string `json:"id"`
	Password    string `json:"password"`
	OldPassword string `json:"old_password"`
}
