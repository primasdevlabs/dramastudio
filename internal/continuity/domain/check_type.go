package domain

type CheckType string

const (
	CheckStory     CheckType = "story"
	CheckTimeline  CheckType = "timeline"
	CheckCharacter CheckType = "character"
	CheckWardrobe  CheckType = "wardrobe"
	CheckVisual    CheckType = "visual"
)
