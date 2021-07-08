package model

type Book struct {
	Id string `json:"id"`
	Title string `json:"title"`
	Author string `json:"author"`
	Abstract string `json:"abstract"`
}

type CreateBook struct {
	Title string
	Author string
	Abstract string
}
