package models

type Actor struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	BirthDate Date   `json:"birthDate"`
}
