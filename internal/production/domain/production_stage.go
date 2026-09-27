package domain

type ProductionStage string

const (
	StagePreProduction  ProductionStage = "pre_production"
	StageGeneration     ProductionStage = "generation"
	StagePostProduction ProductionStage = "post_production"
)
