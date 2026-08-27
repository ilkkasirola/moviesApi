package models

type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name" validate:"omitempty, min=1"`
}

type GenrePatch struct {
	Name *string `json:"name"`
}
