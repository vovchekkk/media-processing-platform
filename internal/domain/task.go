package domain

type Task struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Result      string `json:"result"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}