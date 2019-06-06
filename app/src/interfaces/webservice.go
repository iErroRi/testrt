package interfaces

import (
	"domain"
	"encoding/json"
	"fmt"
	"net/http"
)

type CodeInteractor interface {
	Code(countryName string) (domain.Code, error)
}

type ReloadInteractor interface {
	Reload() (int, error)
}

type WebserviceHandler struct {
	CodeInteractor   CodeInteractor
	ReloadInteractor ReloadInteractor
}

func (handler WebserviceHandler) Code(res http.ResponseWriter, req *http.Request, params map[string]string) {
	code, _ := handler.CodeInteractor.Code(params["name"])

	if code.Id == "" {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	res.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(res).Encode(code); err != nil {
		fmt.Println(err)
	}
}

func (handler WebserviceHandler) Reload(res http.ResponseWriter, req *http.Request, params map[string]string) {
	res.Header().Set("Content-Type", "application/json")

	count, err := handler.ReloadInteractor.Reload()
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
	}

	if err := json.NewEncoder(res).Encode(fmt.Sprint("Complete loaded ", count)); err != nil {
		fmt.Println(err)
	}
}
