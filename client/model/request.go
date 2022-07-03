package model

type WriteRequest struct {
	Raw string `json:"raw"`
	Key string `json:"key"`
}
