package request

type CreateTask struct {
	TaskName  string   `json:"task_name"`
	Assignee  []string `json:"assignee"`
	ProjectID string   `json:"project_id"`
	StartDate string   `json:"start_date"`
	EndDate   string   `json:"end_date"`
}
