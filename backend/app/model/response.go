package model

type TodosResponse[T Todo | []Todo] struct {
	Data   T          `json:"data"`
	Status StatusInfo `json:"status"`
}
