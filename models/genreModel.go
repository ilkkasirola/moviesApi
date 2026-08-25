package models

type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type GenrePatch struct {
	Name *string `json:"name"`
}
