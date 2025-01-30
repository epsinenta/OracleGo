package entities

type Hero struct {
	Value string
}

func (h Hero) GetValue() string {
	return h.Value
}

type Winrate struct {
	Hero    Hero
	Winrate float64
}

type CounterRate struct {
	FirstHero   Hero
	SecondHero  Hero
	CounterPick float64
}
