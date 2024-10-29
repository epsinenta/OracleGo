package db

import (
	"reflect"
)

type NamedEntity interface {
	GetValue() string
}

type Hero struct {
	Value string
}

func (h Hero) GetValue() string {
	return h.Value
}

type Team struct {
	Value string
}

func (t Team) GetValue() string {
	return t.Value
}

type Player struct {
	Value string
}

func (p Player) GetValue() string {
	return p.Value
}

type Email struct {
	Value string
}

func (e Email) GetValue() string {
	return e.Value
}

type Password struct {
	Value string
}

func (p Password) GetValue() string {
	return p.Value
}

func ValuesFromAny(entities interface{}) []string {
	v := reflect.ValueOf(entities)

	if v.Kind() != reflect.Slice {
		panic("expected a slice")
	}

	values := make([]string, v.Len())

	for i := 0; i < v.Len(); i++ {
		entity := v.Index(i).Interface().(NamedEntity)
		values[i] = entity.GetValue()
	}

	return values
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

type User struct {
	Email    Email
	Password Password
}
