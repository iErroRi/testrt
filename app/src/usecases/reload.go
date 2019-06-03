package usecases

import (
	"domain"
	"encoding/json"
	"infrastructure"
	"net/http"
	"time"
)

type ReloadInteractor struct {
	CountryRepository domain.CountryRepository
	CodeRepository    domain.CodeRepository
	Logger            *infrastructure.Logger
}

func (interactor *ReloadInteractor) getJson(url string, target interface{}) error {
	var myClient = &http.Client{Timeout: 10 * time.Second}
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func (interactor *ReloadInteractor) Reload() (int, error) {
	interactor.Logger.Info("Start reload")

	countries := make(map[string]string, 0)
	if err := interactor.getJson("http://country.io/names.json", &countries); err != nil {
		interactor.Logger.Info("Error get names.json " + err.Error())
		return 0, err
	}

	interactor.CountryRepository.Clear()

	for key, val := range countries {
		country := domain.Country{Id: key, Name: val}
		interactor.CountryRepository.Store(country)
	}
	interactor.Logger.Info("Load countries")

	codes := make(map[string]string, 0)

	if err := interactor.getJson("http://country.io/phone.json", &codes); err != nil {
		interactor.Logger.Info("Error get phone.json " + err.Error())
		return 0, err
	}

	interactor.CodeRepository.Clear()
	for key, val := range codes {
		code := domain.Code{Id: key, Name: val}
		interactor.CodeRepository.Store(code)
	}
	interactor.Logger.Info("Load codes")

	interactor.Logger.Info("Stop reload")
	return len(countries) + len(codes), nil
}
