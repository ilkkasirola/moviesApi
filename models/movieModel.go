package models

type Movie struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	ReleaseYear int     `json:"releaseYear"`
	Duration    int     `json:"duration"`
	Genres      []Genre `json:"genres"`
	Actors      []Actor `json:"actors" validate:"dive"`
}
