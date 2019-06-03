package main

import (
	"infrastructure"
	"interfaces"
	"net/http"
	"usecases"
)

func main() {
	logger := new(infrastructure.Logger)

	dbHandler := infrastructure.NewMysqlHandler("user:1234@(mysql_db:3306)/testrt", logger)

	handlers := make(map[string]interfaces.DbHandler)
	handlers["DbCountryRepo"] = dbHandler
	handlers["DbCodeRepo"] = dbHandler

	codeInteractor := new(usecases.CodeInteractor)
	codeInteractor.CountryRepository = interfaces.NewDbCountryRepo(handlers)
	codeInteractor.CodeRepository = interfaces.NewDbCodeRepo(handlers)
	codeInteractor.Logger = logger

	reloadInteractor := new(usecases.ReloadInteractor)
	reloadInteractor.CountryRepository = interfaces.NewDbCountryRepo(handlers)
	reloadInteractor.CodeRepository = interfaces.NewDbCodeRepo(handlers)
	reloadInteractor.Logger = logger

	webserviceHandler := interfaces.WebserviceHandler{}
	webserviceHandler.CodeInteractor = codeInteractor
	webserviceHandler.ReloadInteractor = reloadInteractor

	if _, err := reloadInteractor.Reload(); err != nil {
		logger.Error("Ошибка обновления базы данных городов и кодов")
	}

	server := infrastructure.NewServer(logger)

	server.AddRoute("GET", "/code/{name:[a-zA-Z ]+}", func(res http.ResponseWriter, req *http.Request) {
		webserviceHandler.Code(res, req, server.Params(req))
	})

	server.AddRoute("POST", "/reload", func(res http.ResponseWriter, req *http.Request) {
		webserviceHandler.Reload(res, req, server.Params(req))
	})

	server.ListenAndServe()
}
