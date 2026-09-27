package persistence

type CharacterModel struct {
	ID       string `gorm:"primaryKey"`
	Name     string
	IsLocked bool
}
