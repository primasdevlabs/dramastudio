package domain

// Definitions are the versioned production-intelligence vocabulary (§44–46
// production governance): who does the work (Agent), how (Skill), under what
// operational constraints (Policy), which hard constraints are machine-
// enforced (Rule), and what independently evaluates output (Evaluator).
// Definitions ship as editable catalog data, never embedded in workflow code.

// Agent is a role definition: responsibilities, reusable skills, applicable
// policies and rules, permissions, and the capability its reasoning routes to.
type Agent struct {
	ID               string            `json:"id" yaml:"id"`
	Version          int               `json:"version" yaml:"version"`
	Name             string            `json:"name" yaml:"name"`
	Role             string            `json:"role" yaml:"role"`
	Responsibilities []string          `json:"responsibilities" yaml:"responsibilities"`
	SkillIDs         []string          `json:"skills" yaml:"skills"`
	PolicyIDs        []string          `json:"policies" yaml:"policies"`
	RuleIDs          []string          `json:"rules" yaml:"rules"`
	EvaluatorIDs     []string          `json:"evaluators" yaml:"evaluators"`
	Permissions      PermissionSet     `json:"permissions" yaml:"permissions"`
	ModelCapability  string            `json:"model_capability" yaml:"model_capability"` // e.g. lead_director, vision
	ContextSpec      ContextSpec       `json:"context" yaml:"context"`                   // which slices of production state it may see
	Metadata         map[string]string `json:"metadata,omitempty" yaml:"metadata"`
}

func (d *Agent) GetID() string     { return d.ID }
func (d *Skill) GetID() string     { return d.ID }
func (d *Policy) GetID() string    { return d.ID }
func (d *Rule) GetID() string      { return d.ID }
func (d *Evaluator) GetID() string { return d.ID }

// PermissionSet scopes what an agent may read, write, and approve.
type PermissionSet struct {
	Read    []string `json:"read" yaml:"read"`
	Write   []string `json:"write" yaml:"write"`
	Approve []string `json:"approve" yaml:"approve"`
}

// ContextSpec declares which production-state slices the resolver gathers
// into the execution context (keeps model calls narrow, not the whole bible).
type ContextSpec struct {
	SeriesBible      bool `json:"series_bible" yaml:"series_bible"`
	SeasonArc        bool `json:"season_arc" yaml:"season_arc"`
	EpisodeOutline   bool `json:"episode_outline" yaml:"episode_outline"`
	Scene            bool `json:"scene" yaml:"scene"`
	ShotSpec         bool `json:"shot_spec" yaml:"shot_spec"`
	Characters       bool `json:"characters" yaml:"characters"`
	CharacterRefs    bool `json:"character_refs" yaml:"character_refs"`
	Wardrobe         bool `json:"wardrobe" yaml:"wardrobe"`
	Canon            bool `json:"canon" yaml:"canon"`
	Timeline         bool `json:"timeline" yaml:"timeline"`
	PlotThreads      bool `json:"plot_threads" yaml:"plot_threads"`
	PreviousEpisodes bool `json:"previous_episodes" yaml:"previous_episodes"`
	StoryGraph       bool `json:"story_graph" yaml:"story_graph"`
	VisualBible      bool `json:"visual_bible" yaml:"visual_bible"`
	Locations        bool `json:"locations" yaml:"locations"`
	Props            bool `json:"props" yaml:"props"`
	Assets           bool `json:"assets" yaml:"assets"`
}

// Skill is instruction + methodology + I/O contract + evaluation criteria,
// reusable across agents (cinematography serves Video Director, Storyboard
// Director, Shot Designer, and Quality Supervisor alike).
type Skill struct {
	ID         string   `json:"id" yaml:"id"`
	Version    int      `json:"version" yaml:"version"`
	Domain     string   `json:"domain" yaml:"domain"` // story|visual|production|audio|postproduction|quality
	Purpose    string   `json:"purpose" yaml:"purpose"`
	Inputs     []string `json:"inputs" yaml:"inputs"`
	Outputs    []string `json:"outputs" yaml:"outputs"`
	Method     []string `json:"method" yaml:"method"` // ordered methodology → prompt lines
	Evaluation []string `json:"evaluation" yaml:"evaluation"`
}

// PolicyLayer is the precedence level a policy applies at; lower layers may
// narrow but never violate keys the higher layer marked protected.
type PolicyLayer string

const (
	LayerSystem  PolicyLayer = "system"
	LayerStudio  PolicyLayer = "studio"
	LayerProject PolicyLayer = "project"
	LayerSeries  PolicyLayer = "series"
	LayerSeason  PolicyLayer = "season"
	LayerEpisode PolicyLayer = "episode"
	LayerTask    PolicyLayer = "task"
)

// layerOrder ranks precedence: index 0 is the strongest layer.
var layerOrder = []PolicyLayer{
	LayerSystem, LayerStudio, LayerProject, LayerSeries,
	LayerSeason, LayerEpisode, LayerTask,
}

// Policy is a versioned operational constraint set scoped to a layer.
type Policy struct {
	ID           string                 `json:"id" yaml:"id"`
	Version      int                    `json:"version" yaml:"version"`
	Layer        PolicyLayer            `json:"layer" yaml:"layer"`
	ScopeID      string                 `json:"scope_id,omitempty" yaml:"scope_id"`     // "" = wildcard for the layer
	AppliesTo    []string               `json:"applies_to,omitempty" yaml:"applies_to"` // actions/capabilities/scopes
	Requirements map[string]interface{} `json:"requirements,omitempty" yaml:"requirements"`
	Protected    []string               `json:"protected,omitempty" yaml:"protected"`     // requirement keys lower layers cannot weaken
	Constraints  []string               `json:"constraints,omitempty" yaml:"constraints"` // instruction lines injected into prompts
}

// Severity classes for rule findings (same ladder continuity uses).
type Severity string

const (
	SeverityInfo     Severity = "INFO"
	SeverityWarning  Severity = "WARNING"
	SeverityError    Severity = "ERROR"
	SeverityBlocking Severity = "BLOCKING"
)

// Enforcement decides what a violated rule does to the run.
type Enforcement string

const (
	EnforceAdvisory      Enforcement = "ADVISORY"       // record only
	EnforceRequireReview Enforcement = "REQUIRE_REVIEW" // flags human approval
	EnforceBlock         Enforcement = "BLOCK"          // halts the execution
	EnforceAutoFix       Enforcement = "AUTO_FIX"       // records fix-needed marker
)

// Rule is a small deterministic constraint checked by the rules engine —
// never left to model interpretation. Predicates evaluate over the task's
// fact map (e.g. "character.canonical_reference != null").
type Rule struct {
	ID          string      `json:"id" yaml:"id"`
	Version     int         `json:"version" yaml:"version"`
	Name        string      `json:"name" yaml:"name"`
	When        RuleTrigger `json:"when" yaml:"when"`
	Requires    []string    `json:"requires" yaml:"requires"` // predicate expressions
	Severity    Severity    `json:"severity" yaml:"severity"`
	Enforcement Enforcement `json:"enforcement" yaml:"enforcement"`
	OnFailure   RuleFailure `json:"on_failure" yaml:"on_failure"`
}

// RuleTrigger selects which actions a rule guards ("*" = every action).
type RuleTrigger struct {
	Action string `json:"action" yaml:"action"`
	Phase  string `json:"phase,omitempty" yaml:"phase"` // pre|post; empty = both
}

type RuleFailure struct {
	Action  string `json:"action" yaml:"action"` // block|flag|annotate
	Message string `json:"message" yaml:"message"`
}

// Evaluator independently assesses output quality — separated from the
// generator so generation and judgement don't share a prompt. Checks are
// rule IDs evaluated deterministically; ModelCapability, when set, routes a
// review call through the model policy (e.g. quality_evaluation).
type Evaluator struct {
	ID              string   `json:"id" yaml:"id"`
	Version         int      `json:"version" yaml:"version"`
	Name            string   `json:"name" yaml:"name"`
	Target          string   `json:"target" yaml:"target"` // script|image|video|audio|episode
	RuleIDs         []string `json:"rules" yaml:"rules"`
	SkillIDs        []string `json:"skills" yaml:"skills"` // evaluation criteria reused from skills
	ModelCapability string   `json:"model_capability,omitempty" yaml:"model_capability"`
	Rubric          []string `json:"rubric,omitempty" yaml:"rubric"` // review prompt criteria
}
