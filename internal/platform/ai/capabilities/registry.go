package capabilities

import (
	"context"
	"fmt"
	"sync"
)

type ModelPolicy struct {
	Provider string  `json:"provider"`
	Model    string  `json:"model"`
	MaxCost  float64 `json:"max_cost"`
}

type ModelRegistry struct {
	mu       sync.RWMutex
	policies map[CapabilityType]ModelPolicy
}

func NewModelRegistry() *ModelRegistry {
	r := &ModelRegistry{
		policies: make(map[CapabilityType]ModelPolicy),
	}
	r.SetPolicy(CapVideoGen, ModelPolicy{Provider: "wan", Model: "wan-2.1-t2v", MaxCost: 0.20})
	r.SetPolicy(CapImageGen, ModelPolicy{Provider: "flux", Model: "flux-1-dev", MaxCost: 0.05})
	r.SetPolicy(CapScriptWriting, ModelPolicy{Provider: "openai", Model: "gpt-4o", MaxCost: 0.02})
	r.SetPolicy(CapDialogueWriting, ModelPolicy{Provider: "openai", Model: "gpt-4o", MaxCost: 0.02})
	return r
}

func (r *ModelRegistry) SetPolicy(cap CapabilityType, policy ModelPolicy) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.policies[cap] = policy
}

func (r *ModelRegistry) GetPolicy(cap CapabilityType) (ModelPolicy, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	policy, ok := r.policies[cap]
	if !ok {
		return ModelPolicy{}, fmt.Errorf("no model policy registered for capability: %s", cap)
	}
	return policy, nil
}

type CapabilityRouter struct {
	registry *ModelRegistry
}

func NewCapabilityRouter(registry *ModelRegistry) *CapabilityRouter {
	return &CapabilityRouter{registry: registry}
}

func (r *CapabilityRouter) Execute(ctx context.Context, req GenerationRequest) (*GenerationResponse, error) {
	policy, err := r.registry.GetPolicy(req.Capability)
	if err != nil {
		policy = ModelPolicy{Provider: "wan", Model: "wan-2.1-t2v", MaxCost: 0.10}
	}

	outputURL := fmt.Sprintf("https://storage.dramastudio.ai/gen/%s_%s.mp4", policy.Provider, req.Capability)
	return &GenerationResponse{
		OutputURL: outputURL,
		Provider:  policy.Provider,
		Model:     policy.Model,
		Cost:      policy.MaxCost,
	}, nil
}
