package domain

import (
	"time"
)

type ProjectStatus string

const (
	ProjectStatusCreated    ProjectStatus = "PROJECT_CREATED"
	ProjectStatusDeveloping ProjectStatus = "DEVELOPING_SERIES"
	ProjectStatusApproved   ProjectStatus = "SERIES_APPROVED"
	ProjectStatusPlanning   ProjectStatus = "PLANNING_SEASON"
	ProjectStatusProducing  ProjectStatus = "PRODUCING_EPISODE"
	ProjectStatusValidating ProjectStatus = "VALIDATING_EPISODE"
	ProjectStatusPostProd   ProjectStatus = "POST_PRODUCTION"
	ProjectStatusReady      ProjectStatus = "READY"
	ProjectStatusPublished  ProjectStatus = "PUBLISHED"
)

type Budget struct {
	MaxCostPerGeneration float64 `json:"max_cost_per_generation"`
	TotalBudget          float64 `json:"total_budget"`
	CurrentSpent         float64 `json:"current_spent"`
	Currency             string  `json:"currency"`
}

type Project struct {
	ID          ProjectID        `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Genre       string           `json:"genre"`
	Language    string           `json:"language"`
	Mode        ProductionMode   `json:"mode"`
	Policy      ProductionPolicy `json:"policy"`
	Budget      Budget           `json:"budget"`
	Status      ProjectStatus    `json:"status"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

func NewProject(id ProjectID, name, description, genre, language string, mode ProductionMode) *Project {
	now := time.Now().UTC()
	return &Project{
		ID:          id,
		Name:        name,
		Description: description,
		Genre:       genre,
		Language:    language,
		Mode:        mode,
		Status:      ProjectStatusCreated,
		Budget: Budget{
			Currency: "USD",
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
}
