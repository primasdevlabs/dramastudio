package domain

// Definition is a registered agent persona: role, instructions, tools,
// permissions, model policy, and budget limits (§44).
type Definition struct {
	ID           string                 `json:"id"`
	Role         AgentRole              `json:"role"`
	Name         string                 `json:"name"`
	Instructions string                 `json:"instructions"`
	Skills       []Skill                `json:"skills"`
	Tools        []string               `json:"tools"`
	Permissions  []string               `json:"permissions"`
	ModelPolicy  map[string]interface{} `json:"model_policy"`
	BudgetLimit  *float64               `json:"budget_limit,omitempty"`
	Active       bool                   `json:"active"`
}
