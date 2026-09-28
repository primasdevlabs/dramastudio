package capability

// TextGeneration covers all reasoning/writing tasks: bible, architecture,
// episodes, dialogue, continuity analysis, QA. A model's capabilities
// decide which task-level policies it can serve.
const TextGeneration = "text_generation"

// TextSpec is the typed parameter vocabulary adapters map onto provider
// requests (chat completions, messages APIs, etc.).
type TextSpec struct {
	SystemPrompt string                 `json:"system_prompt,omitempty"`
	Prompt       string                 `json:"prompt"`
	MaxTokens    int                    `json:"max_tokens,omitempty"`
	Temperature  float64                `json:"temperature,omitempty"`
	JSONMode     bool                   `json:"json_mode,omitempty"`
	Tools        []map[string]any       `json:"tools,omitempty"`
	Extra        map[string]interface{} `json:"extra,omitempty"`
}
