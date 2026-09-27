package domain

type WardrobeItem struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Wardrobe struct {
	Outfits []WardrobeItem `json:"outfits"`
}
