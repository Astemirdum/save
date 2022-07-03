package models

type File struct {
	Filename string `json:"filename" validate:"required"`
}
