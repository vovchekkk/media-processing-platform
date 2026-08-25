package dto

type Task struct{
	Filter ImageFilter `json:"filter"`
	Image  string      `json:"image"`
}

type ImageFilter struct {
	Name string `json:"name"`
}