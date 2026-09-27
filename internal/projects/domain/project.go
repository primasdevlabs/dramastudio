package domain

type Project struct {
	ID     ProjectID        `json:"id"`
	Name   string           `json:"name"`
	Mode   ProductionMode   `json:"mode"`
	Policy ProductionPolicy `json:"policy"`
}
