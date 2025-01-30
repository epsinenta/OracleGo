package entities

type Team struct {
	Value string
}

func (t Team) GetValue() string {
	return t.Value
}

type TeamRoaster struct {
	Players []Player
	Team    Team
}
