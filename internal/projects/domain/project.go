package domain

import (
	"fmt"
	"time"
)

// ProjectStatus is the project-level production state machine (§32).
type ProjectStatus string

const (
	ProjectStatusCreated    ProjectStatus = "PROJECT_CREATED"
	ProjectStatusDeveloping ProjectStatus = "DEVELOPING_SERIES"
	ProjectStatusApproved   ProjectStatus = "SERIES_APPROVED"
	ProjectStatusPlanning   ProjectStatus = "PLANNING_SEASON"
	ProjectStatusPlanningEp ProjectStatus = "PLANNING_EPISODES"
	ProjectStatusProducing  ProjectStatus = "PRODUCING_EPISODE"
	ProjectStatusValidating ProjectStatus = "VALIDATING_EPISODE"
	ProjectStatusPostProd   ProjectStatus = "POST_PRODUCTION"
	ProjectStatusReady      ProjectStatus = "READY"
	ProjectStatusPublished  ProjectStatus = "PUBLISHED"
)

// validTransitions encodes the §32 state machine.
var validTransitions = map[ProjectStatus][]ProjectStatus{
	ProjectStatusCreated:    {ProjectStatusDeveloping},
	ProjectStatusDeveloping: {ProjectStatusApproved, ProjectStatusDeveloping},
	ProjectStatusApproved:   {ProjectStatusPlanning},
	ProjectStatusPlanning:   {ProjectStatusPlanningEp},
	ProjectStatusPlanningEp: {ProjectStatusProducing},
	ProjectStatusProducing:  {ProjectStatusValidating, ProjectStatusPlanningEp},
	ProjectStatusValidating: {ProjectStatusPostProd, ProjectStatusProducing},
	ProjectStatusPostProd:   {ProjectStatusReady},
	ProjectStatusReady:      {ProjectStatusPublished, ProjectStatusPlanningEp},
	ProjectStatusPublished:  {ProjectStatusPlanningEp},
}

// CanTransition reports whether the state machine allows from→to.
func CanTransition(from, to ProjectStatus) bool {
	for _, next := range validTransitions[from] {
		if next == to {
			return true
		}
	}
	return false
}

// Transition validates and applies a status change.
func (p *Project) Transition(to ProjectStatus) error {
	if !CanTransition(p.Status, to) {
		return fmt.Errorf("invalid project transition %s → %s", p.Status, to)
	}
	p.Status = to
	p.UpdatedAt = time.Now().UTC()
	return nil
}

type Budget struct {
	MaxCostPerGeneration float64 `json:"max_cost_per_generation"`
	TotalBudget          float64 `json:"total_budget"`
	CurrentSpent         float64 `json:"current_spent"`
	Currency             string  `json:"currency"`
}

// CanSpend reports whether a generation cost fits the budget (§85).
func (b Budget) CanSpend(cost float64) bool {
	if b.TotalBudget > 0 && b.CurrentSpent+cost > b.TotalBudget {
		return false
	}
	if b.MaxCostPerGeneration > 0 && cost > b.MaxCostPerGeneration {
		return false
	}
	return true
}

type Project struct {
	ID          ProjectID        `json:"id"`
	OrgID       string           `json:"org_id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Genre       string           `json:"genre"`
	Language    string           `json:"language"`
	Mode        ProductionMode   `json:"mode"`
	Policy      ProductionPolicy `json:"policy"`
	Settings    Settings         `json:"settings"`
	Budget      Budget           `json:"budget"`
	Status      ProjectStatus    `json:"status"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

// Settings carries the project-level configuration of §9.
type Settings struct {
	TargetPlatforms []string `json:"target_platforms,omitempty"`
	AspectRatio     string   `json:"aspect_ratio,omitempty"`
	AutonomyLevel   string   `json:"autonomy_level,omitempty"` // monitored|autonomous
	CreativeRules   []string `json:"creative_rules,omitempty"`
	PublishingRules []string `json:"publishing_rules,omitempty"`
}

func NewProject(id ProjectID, orgID, name, description, genre, language string, mode ProductionMode) *Project {
	now := time.Now().UTC()
	return &Project{
		ID:          id,
		OrgID:       orgID,
		Name:        name,
		Description: description,
		Genre:       genre,
		Language:    language,
		Mode:        mode,
		Status:      ProjectStatusCreated,
		Budget:      Budget{Currency: "USD"},
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
