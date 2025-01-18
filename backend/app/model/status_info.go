package model

type StatusInfo struct {
	Code         int    `json:"code"`
	Error        bool   `json:"error"`
	ErrorMessage string `json:"error_message"`
}
