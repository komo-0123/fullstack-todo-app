package model

type TodoTypes interface {
	Todo | []Todo
}

type TodoResponse[T TodoTypes] struct {
	Data   T          `json:"data"`
	Status StatusInfo `json:"status"`
}
