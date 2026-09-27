package domain

type ProductionJob struct {
	ID        string    `json:"id"`
	ShotID    string    `json:"shot_id"`
	Status    JobStatus `json:"status"`
	ResultURL string    `json:"result_url"`
}
