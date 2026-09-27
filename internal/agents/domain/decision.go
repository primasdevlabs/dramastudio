package domain

import "time"

type DecisionMaker string

const (
	DecisionMakerUser             DecisionMaker = "USER"
	DecisionMakerLeadDirector     DecisionMaker = "LEAD_DIRECTOR"
	DecisionMakerSpecializedAgent DecisionMaker = "SPECIALIZED_AGENT"
	DecisionMakerSystem           DecisionMaker = "SYSTEM"
)

type Decision struct {
	ID            string        `json:"id"`
	ProjectID     string        `json:"project_id"`
	EpisodeID     string        `json:"episode_id"`
	Decision      string        `json:"decision"`
	Reason        string        `json:"reason"`
	DecisionMaker DecisionMaker `json:"decision_maker"`
	Mode          string        `json:"mode"`
	Timestamp     time.Time     `json:"timestamp"`
}
