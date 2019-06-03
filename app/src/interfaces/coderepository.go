package interfaces

import (
	"domain"
)

type DbCodeRepo DbRepo

func NewDbCodeRepo(dbHandlers map[string]DbHandler) *DbCodeRepo {
	dbCodeRepo := new(DbCodeRepo)
	dbCodeRepo.dbHandlers = dbHandlers
	dbCodeRepo.dbHandler = dbHandlers["DbCodeRepo"]
	return dbCodeRepo
}

func (repo *DbCodeRepo) Store(code domain.Code) {
	repo.dbHandler.Execute("INSERT INTO code (id, name) VALUES (:id, :name)", code)
}

func (repo *DbCodeRepo) FindById(id string) domain.Code {
	row := repo.dbHandler.Query("SELECT name FROM code WHERE id = :id", map[string]interface{}{"id": id})

	var name string
	row.Next()
	row.Scan(&name)
	code := domain.Code{Id: id, Name: name}
	return code
}

func (repo *DbCodeRepo) Clear() {
	repo.dbHandler.Trunc("code")
}
