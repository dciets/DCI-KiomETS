package serialisation

type PlayerSerialisation struct {
	Name              string   `json:"name"`
	Color             string   `json:"color"`
	NumberOfKill      int      `json:"numberOfKill"`
	PossessedTerrains []string `json:"possessedTerrains"`
}

func NewPlayerSerialisation(name string, color string, numberOfKill int, possessedTerrains []string) *PlayerSerialisation {
	return &PlayerSerialisation{Name: name, Color: color, NumberOfKill: numberOfKill, PossessedTerrains: possessedTerrains}
}
