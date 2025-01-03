package model

type Todo struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	IsComplete bool   `json:"is_complete"`
}
