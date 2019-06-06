package infrastructure

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func NewServer(logger Log) *Server {
	server := new(Server)
	server.Instance = mux.NewRouter()
	server.logger = logger
	return server
}

type Server struct {
	Instance *mux.Router
	logger   Log
}

func (s Server) AddRoute(method string, url string, handler func(res http.ResponseWriter, req *http.Request)) {
	s.logger.Info(fmt.Sprintf("Add route method '%v' url '%v", method, url))
	s.Instance.HandleFunc(url, handler).Methods(method)
}

func (s Server) Params(req *http.Request) map[string]string {
	return mux.Vars(req)
}

func (s Server) ListenAndServe() {
	s.logger.Info("Start web server")
	http.Handle("/", s.Instance)

	s.Instance.NotFoundHandler = http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusNotFound)
		s.logger.Info("Not found " + req.URL.Path)
	})

	s.Instance.MethodNotAllowedHandler = http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusMethodNotAllowed)
		s.logger.Info("Not allow " + req.URL.Path)
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		s.logger.Info("Failed start web server " + err.Error())
		panic(err)
	}
}
