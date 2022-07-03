package models

type WriteRequest struct {
	Raw        string `json:"raw" validate:"required"`
	Key        string `json:"key" validate:"required"`
	TimeStamp  int64  `json:"-"`
	ClientName string `json:"-"`
}
