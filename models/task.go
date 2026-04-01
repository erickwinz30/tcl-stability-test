package models

type Task struct {
	ID int `json:"id"`
	// add required for title
	Title string `json:"title" validate:"required"`
	Done  bool   `json:"done"`
}

type UpdateTaskRequest struct {
	Title string `json:"title" validate:"required,min=3"`
	Done  *bool  `json:"done" validate:"required"`
}
