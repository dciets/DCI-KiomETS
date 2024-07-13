package serialisation

type PlayerSerialisation struct {
	Name                   string `json:"name"`
	Color                  string `json:"color"`
	NumberOfKill           int    `json:"numberOfKill"`
	PossessedTerrainsCount int    `json:"possessedTerrainsCount"`
}

func NewPlayerSerialisation(name string, color string, numberOfKill int, possessedTerrainsCount int) *PlayerSerialisation {
	return &PlayerSerialisation{Name: name, Color: color, NumberOfKill: numberOfKill, PossessedTerrainsCount: possessedTerrainsCount}
}
