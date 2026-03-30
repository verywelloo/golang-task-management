package request

type CreateTask struct {
	TaskName  string `json:"task_name"`
	ProjectID string `json:"project_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}
