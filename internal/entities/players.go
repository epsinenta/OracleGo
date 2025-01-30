package entities

type Player struct {
	Value string
}

func (p Player) GetValue() string {
	return p.Value
}

type GamesCount struct {
	Player Player
	Hero   Hero
	Count  int
}

type PlayerWinrate struct {
	Player  Player
	Hero    Hero
	Winrate float64
}

type PlayerWinrateCurPatch struct {
	Player  Player
	Winrate float64
}
