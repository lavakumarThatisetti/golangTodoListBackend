package models

//Todo Struct (Model)
type Todo struct {
	ID        string `json:"id"`
	UUID      string `json:"uuid"`
	Title     string `json:"title"`
	Text      string `json:"text"`
	Completed bool   `json:"completed"`
}
