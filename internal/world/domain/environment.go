package domain

type Environment struct {
	Lighting  string `json:"lighting"`
	Weather   string `json:"weather"`
	TimeOfDay string `json:"time_of_day"`
}
