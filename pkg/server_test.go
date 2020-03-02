package roboot

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matryer/is"
)

func Test_server_404(t *testing.T) {
	is := is.New(t)
	srv := NewServer()

	r := httptest.NewRequest(http.MethodGet, "/404", nil)
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, r)

	is.Equal(w.Result().StatusCode, http.StatusNotFound)
}

func Test_server_handleHealth(t *testing.T) {
	is := is.New(t)
	srv := NewServer()

	r := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	srv.handleHealth()(w, r)

	is.Equal(w.Result().StatusCode, http.StatusOK)
}

func Test_server_onlyFleetLockProtocol(t *testing.T) {
	is := is.New(t)
	srv := NewServer()

	n := 0
	h := func(w http.ResponseWriter, r *http.Request) {
		n++
	}
	h = srv.onlyFleetLockProtocol(h)

	// test: no header
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	w := httptest.NewRecorder()
	h(w, r)
	is.Equal(n, 0)
	is.Equal(w.Result().StatusCode, http.StatusBadRequest)

	// test: wrong header
	r = httptest.NewRequest(http.MethodPost, "/", nil)
	w = httptest.NewRecorder()

	r.Header.Set("fleet-lock-protocol", "wrong")

	h(w, r)
	is.Equal(n, 0)
	is.Equal(w.Result().StatusCode, http.StatusBadRequest)

	// test: correct header
	r = httptest.NewRequest(http.MethodPost, "/", nil)
	w = httptest.NewRecorder()

	r.Header.Set("fleet-lock-protocol", "true")

	h(w, r)
	is.Equal(n, 1)
}
