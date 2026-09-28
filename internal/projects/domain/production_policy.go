package domain

type ProductionPolicy struct {
	AllowAutoPublish     bool     `json:"allow_auto_publish"`
	MaxDailyGenerations  int      `json:"max_daily_generations"`
	AllowedAIPrompting   bool     `json:"allowed_ai_prompting"`
	RequireHumanApproval []string `json:"require_human_approval"`
}
