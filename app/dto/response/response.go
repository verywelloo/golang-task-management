package response

type Result struct {
	Status  int         `json:"status,omitempty"`
	Message string      `json:"message,omitempty"`
	Details interface{} `json:"details,omitempty"`
}
