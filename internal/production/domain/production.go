package domain

type Production struct {
	ID        string        `json:"id"`
	ProjectID string        `json:"project_id"`
	Runs      []ProductionRun `json:"runs"`
}
