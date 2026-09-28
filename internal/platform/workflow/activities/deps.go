package activities

import (
	aiapp "dramastudio/internal/ai/application"
	continuitysvc "dramastudio/internal/continuity/application/services"
	mediasvc "dramastudio/internal/media/application/services"
	postsvc "dramastudio/internal/postproduction/application/services"
	prodsvc "dramastudio/internal/production/application/services"
	storysvc "dramastudio/internal/story/application/services"
)

// Deps wires the application services Temporal activities call. Bound at
// worker construction; activities must never hold per-request state.
type Deps struct {
	Story       *storysvc.StoryService
	Production  *prodsvc.ProductionService
	Media       *mediasvc.MediaService
	Continuity  *continuitysvc.ContinuityService
	Postprod    *postsvc.PostproductionService
	AI          *aiapp.ExecuteGenerationHandler
	CallbackURL string // base URL for provider webhooks (§61)
}
