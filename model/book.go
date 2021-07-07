package model

type Book struct {
	Title interface{}
	Author interface{}
	Abstract interface{}
}

type CreateBook struct {
	Title string
	Author string
	Abstract string
}
