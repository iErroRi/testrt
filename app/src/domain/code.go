package domain

type CodeRepository interface {
	Store(code Code)
	FindById(id string) Code
	Clear()
}

type Code struct {
	Id   string
	Name string
}
