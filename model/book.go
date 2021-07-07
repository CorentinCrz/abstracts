package model

type Book struct {
	Id interface{}
	Title interface{}
	Author interface{}
	Abstract interface{}
}

type CreateBook struct {
	Title string
	Author string
	Abstract string
}
