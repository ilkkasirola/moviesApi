package models

type Actor struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	BirthDate string `json:"birthDate" validate:"required,datetime=2006-01-02"`
}
