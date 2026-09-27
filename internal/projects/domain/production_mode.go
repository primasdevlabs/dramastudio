package domain

type ProductionMode string

const (
	ModeMonitored  ProductionMode = "monitored"
	ModeAutonomous ProductionMode = "autonomous"
)
