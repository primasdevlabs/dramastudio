package domain

type ModelPolicy struct {
	ID           string                  `json:"id"`
	ProjectID    string                  `json:"project_id"`
	ProductionID string                  `json:"production_id"`
	Mapping      map[AICapability]ModelID `json:"mapping"`
}

func NewModelPolicy(projectID string) *ModelPolicy {
	return &ModelPolicy{
		ProjectID: projectID,
		Mapping:   make(map[AICapability]ModelID),
	}
}
