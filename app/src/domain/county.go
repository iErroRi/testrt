package domain

type CountryRepository interface {
	Store(country Country)
	FindById(id string) Country
	FindByName(name string) Country
	Clear()
}

type Country struct {
	Id   string
	Name string
}
