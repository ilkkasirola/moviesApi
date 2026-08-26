package models

type Actor struct {
	ID        int    `json:"id"`
	Name      string `json:"name" validate:"required"`
	BirthDate string `json:"birthDate" validate:"required,datetime=2006-01-02"`
}
type ActorPatch struct {
	Name      *string `json:"name" validate:"omitempty,min=1"`
	BirthDate *string `json:"birthDate" validate:"omitempty,datetime=2006-01-02"`
}
