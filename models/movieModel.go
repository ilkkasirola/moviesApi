package models

type Movie struct {
	ID          int     `json:"id"`
	Title       string  `json:"title" validate:"required"`
	ReleaseYear int     `json:"releaseYear" validate:"required,gte=1888"`
	Duration    int     `json:"duration" validate:"required,gt=0"`
	Genres      []Genre `json:"genres" validate:"dive"`
	Actors      []Actor `json:"actors" validate:"dive"`
}

type MoviePatch struct {
	Title       *string `json:"title" validate:"omitempty,min=1"`
	ReleaseYear *int    `json:"releaseYear" validate:"omitempty,gte=1888"`
	Duration    *int    `json:"duration" validate:"omitempty,min=1,gt=0"`
	GenreIDs    *[]int  `json:"genres"`
	ActorIDs    *[]int  `json:"actors"`
}
