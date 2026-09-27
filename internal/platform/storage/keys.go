package storage

import "fmt"

// Keys centralizes the object-storage layout (Backend.md §44) so path
// conventions are not duplicated across the codebase.
type Keys struct{}

func NewKeys() Keys { return Keys{} }

func (Keys) ProjectRoot(projectID string) string {
	return fmt.Sprintf("projects/%s", projectID)
}

func (k Keys) Characters(projectID string) string {
	return fmt.Sprintf("%s/characters", k.ProjectRoot(projectID))
}

func (k Keys) Locations(projectID string) string {
	return fmt.Sprintf("%s/locations", k.ProjectRoot(projectID))
}

func (k Keys) EpisodeRoot(projectID, episodeID string) string {
	return fmt.Sprintf("%s/episodes/%s", k.ProjectRoot(projectID), episodeID)
}

func (k Keys) Scene(projectID, episodeID, sceneID string) string {
	return fmt.Sprintf("%s/scenes/%s", k.EpisodeRoot(projectID, episodeID), sceneID)
}

func (k Keys) Shot(projectID, episodeID, sceneID, shotID string) string {
	return fmt.Sprintf("%s/shots/%s", k.Scene(projectID, episodeID, sceneID), shotID)
}

func (k Keys) Audio(projectID, episodeID string) string {
	return fmt.Sprintf("%s/audio", k.EpisodeRoot(projectID, episodeID))
}

func (k Keys) Render(projectID, episodeID, renderID string) string {
	return fmt.Sprintf("%s/renders/%s", k.EpisodeRoot(projectID, episodeID), renderID)
}
