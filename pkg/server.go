package roboot

// https://www.youtube.com/watch?v=rWBSMsLG8po

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type server struct {
	router        *httprouter.Router
	locker        locker
	scheduler     scheduler
	healthChecker healthChecker
}

func NewServer() *server {
	s := &server{
		router: httprouter.New(),
	}
	s.routes()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) routes() {
	s.router.HandlerFunc(http.MethodGet, "/health", s.handleHealth())
	s.router.HandlerFunc(http.MethodPost, "/v1/pre-reboot", s.onlyFleetLockProtocol(s.handlePreReboot()))
	s.router.HandlerFunc(http.MethodPost, "/v1/steady-state", s.onlyFleetLockProtocol(s.handleSteadyState()))
}

func (s *server) handleHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// returns status 200 OK by default
	}
}

func (s *server) onlyFleetLockProtocol(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("fleet-lock-protocol") != "true" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		h(w, r)
	}
}
