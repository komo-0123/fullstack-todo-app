package model

type TodoResponse struct {
	Data   *Todo      `json:"data"`
	Status StatusInfo `json:"status"`
}

type TodosResponse struct {
	Data   []Todo     `json:"data"`
	Status StatusInfo `json:"status"`
}
