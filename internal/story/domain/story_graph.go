package domain

type EdgeType string

const (
	EdgeCauses      EdgeType = "causes"
	EdgeReveals     EdgeType = "reveals"
	EdgeDependsOn   EdgeType = "depends_on"
	EdgeContradicts EdgeType = "contradicts"
	EdgeResolves    EdgeType = "resolves"
	EdgeForeshadows EdgeType = "foreshadows"
	EdgeFollows     EdgeType = "follows"
)

type StoryNode struct {
	ID        string `json:"id"`
	Type      string `json:"type"` // Event, Reveal, Conflict, Decision
	Title     string `json:"title"`
	EpisodeID string `json:"episode_id"`
	SceneID   string `json:"scene_id"`
}

type StoryEdge struct {
	FromID   string   `json:"from_id"`
	ToID     string   `json:"to_id"`
	Relation EdgeType `json:"relation"`
}

type StoryGraph struct {
	Nodes map[string]*StoryNode `json:"nodes"`
	Edges []*StoryEdge          `json:"edges"`
}

func NewStoryGraph() *StoryGraph {
	return &StoryGraph{
		Nodes: make(map[string]*StoryNode),
		Edges: make([]*StoryEdge, 0),
	}
}
