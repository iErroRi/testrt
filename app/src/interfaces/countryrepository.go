package interfaces

import (
	"domain"
)

type DbCountryRepo DbRepo

func NewDbCountryRepo(dbHandlers map[string]DbHandler) *DbCountryRepo {
	dbCountryRepo := new(DbCountryRepo)
	dbCountryRepo.dbHandlers = dbHandlers
	dbCountryRepo.dbHandler = dbHandlers["DbCountryRepo"]
	return dbCountryRepo
}

func (repo *DbCountryRepo) Store(country domain.Country) {
	repo.dbHandler.Execute("INSERT INTO country (id, name) VALUES (:id, :name)", country)
}

func (repo *DbCountryRepo) FindById(id string) domain.Country {
	row := repo.dbHandler.Query("SELECT name FROM code WHERE id = :id", map[string]interface{}{"id": id})

	var name string
	row.Next()
	row.Scan(&name)
	return domain.Country{Id: id, Name: name}
}

func (repo *DbCountryRepo) FindByName(findName string) domain.Country {
	row := repo.dbHandler.Query("SELECT id, name FROM country WHERE name = :name", map[string]interface{}{"name": findName})

	var id, name string
	row.Next()
	row.Scan(&id, &name)
	return domain.Country{Id: id, Name: name}
}

func (repo *DbCountryRepo) Clear() {
	repo.dbHandler.Trunc("country")
}
