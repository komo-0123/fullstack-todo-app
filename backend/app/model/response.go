package model

type TodoResponse[T Todo | []Todo] struct {
	Data   T          `json:"data"`
	Status StatusInfo `json:"status"`
}
