package db

import (
	"reflect"
)

type NamedEntity interface {
	GetName() string
}

type Hero struct {
	Name string
}

func (h Hero) GetName() string {
	return h.Name
}

type Team struct {
	Name string
}

func (t Team) GetName() string {
	return t.Name
}

type Player struct {
	Name string
}

func (p Player) GetName() string {
	return p.Name
}

func NamesFromAny(entities interface{}) []string {
	v := reflect.ValueOf(entities)

	if v.Kind() != reflect.Slice {
		panic("expected a slice")
	}

	names := make([]string, v.Len())

	for i := 0; i < v.Len(); i++ {
		entity := v.Index(i).Interface().(NamedEntity)
		names[i] = entity.GetName()
	}

	return names
}

type TeamRoaster struct {
	Players []Player
	Team    Team
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

type Winrate struct {
	Hero    Hero
	Winrate float64
}

type CounterRate struct {
	FirstHero   Hero
	SecondHero  Hero
	CounterPick float64
}
