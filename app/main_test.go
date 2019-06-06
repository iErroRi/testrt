package main

import (
	"domain"
	"infrastructure"
	"interfaces"
	"net/http"
	"net/http/httptest"
	"testing"
	"usecases"
)

type MockDbCountryRepo interfaces.DbRepo

func (repo *MockDbCountryRepo) Store(country domain.Country) {

}

func (repo *MockDbCountryRepo) FindById(id string) domain.Country {
	if id == "ru" {
		return domain.Country{Id: "ru", Name: "russia"}
	}

	return domain.Country{}
}

func (repo *MockDbCountryRepo) FindByName(findName string) domain.Country {
	if findName == "russia" {
		return domain.Country{Id: "ru", Name: "russia"}
	}

	return domain.Country{}
}

func (repo *MockDbCountryRepo) Clear() {}

type MockDbCodeRepo interfaces.DbRepo

func (repo *MockDbCodeRepo) Store(code domain.Code) {
}

func (repo *MockDbCodeRepo) FindById(id string) domain.Code {
	if id == "ru" {
		return domain.Code{Id: "ru", Name: "7"}
	}

	return domain.Code{}
}

func (repo *MockDbCodeRepo) Clear() {
}

type MockLogger struct{}

func (MockLogger MockLogger) Info(mess string) {
}

func (MockLogger MockLogger) Error(mess string) {
}

func getServer() *infrastructure.Server {
	logger := new(MockLogger)

	codeInteractor := new(usecases.CodeInteractor)

	codeInteractor.CountryRepository = new(MockDbCountryRepo)
	codeInteractor.CodeRepository = new(MockDbCodeRepo)
	codeInteractor.Logger = logger

	webserviceHandler := interfaces.WebserviceHandler{}
	webserviceHandler.CodeInteractor = codeInteractor

	return Server(webserviceHandler, logger)
}

func TestGetFound(t *testing.T) {
	server := getServer()

	req, err := http.NewRequest("GET", "/code/russia", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	server.Instance.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "{\"Id\":\"ru\",\"Name\":\"7\"}\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got -%v- want %v", rr.Body.String(), expected)
	}
}

func TestGetNotFound(t *testing.T) {
	server := getServer()

	req, err := http.NewRequest("GET", "/code/notfound", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	server.Instance.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

func TestGetNotAllow(t *testing.T) {
	server := getServer()

	req, err := http.NewRequest("POST", "/code/russia", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	server.Instance.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}
