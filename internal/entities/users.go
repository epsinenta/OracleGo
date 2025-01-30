package entities

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

type User struct {
	Email    Email
	Password Password
}
