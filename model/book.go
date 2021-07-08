package model

type Book struct {
	Id string
	Title string
	Author string
	Abstract string
}

type CreateBook struct {
	Title string
	Author string
	Abstract string
}
