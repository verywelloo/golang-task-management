package response

type GetTaskResponse struct {
	ID        string `json:"_id"`
	ProjectID string `json:"project_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}
