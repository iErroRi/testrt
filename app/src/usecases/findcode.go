package usecases

import (
	"domain"
	"infrastructure"
)

type CodeInteractor struct {
	CountryRepository domain.CountryRepository
	CodeRepository    domain.CodeRepository
	Logger            infrastructure.Log
}

func (interactor *CodeInteractor) Code(countryName string) (domain.Code, error) {
	interactor.Logger.Info("Start find code")

	country := interactor.CountryRepository.FindByName(countryName)
	code := interactor.CodeRepository.FindById(country.Id)

	interactor.Logger.Info("Stop find code")
	return code, nil
}
