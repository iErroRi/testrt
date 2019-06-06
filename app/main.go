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
	reloadInteractor := new(usecases.ReloadInteractor)

	countryRepo := interfaces.NewDbCountryRepo(handlers)
	codeRepo := interfaces.NewDbCodeRepo(handlers)

	codeInteractor.CountryRepository = countryRepo
	codeInteractor.CodeRepository = codeRepo
	codeInteractor.Logger = logger

	reloadInteractor.CountryRepository = countryRepo
	reloadInteractor.CodeRepository = codeRepo
	reloadInteractor.HttpClient = infrastructure.NewHttpClient()
	reloadInteractor.Logger = logger

	webserviceHandler := interfaces.WebserviceHandler{}
	webserviceHandler.CodeInteractor = codeInteractor
	webserviceHandler.ReloadInteractor = reloadInteractor

	if _, err := reloadInteractor.Reload(); err != nil {
		logger.Error("Ошибка обновления базы данных городов и кодов")
	}

	server := Server(webserviceHandler, logger)
	server.ListenAndServe()
}

func Server(webserviceHandler interfaces.WebserviceHandler, logger infrastructure.Log) *infrastructure.Server {
	server := infrastructure.NewServer(logger)

	server.AddRoute("GET", "/code/{name:[a-zA-Z ]+}", func(res http.ResponseWriter, req *http.Request) {
		webserviceHandler.Code(res, req, server.Params(req))
	})

	server.AddRoute("POST", "/reload", func(res http.ResponseWriter, req *http.Request) {
		webserviceHandler.Reload(res, req, server.Params(req))
	})

	return server
}
