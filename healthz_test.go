package healthz

import (
	"net/http"
	"net/http/httptest"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/matryer/is"
)

func TestNewCheckEmptyLive(t *testing.T) {
	t.Parallel()

	h := NewCheck()
	h.Ready()

	r := httptest.NewRequest(http.MethodGet, "/live", nil)
	w := httptest.NewRecorder()
	h.router().ServeHTTP(w, r)
	h.liveHandler(w, r)
	res1 := w.Result()

	is := is.New(t)
	is.Equal(res1.StatusCode, http.StatusOK)

	r = httptest.NewRequest(http.MethodGet, "/ready", nil)
	w = httptest.NewRecorder()
	h.router().ServeHTTP(w, r)
	h.readyHandler(w, r)
	res2 := w.Result()

	is = is.New(t)
	is.Equal(res2.StatusCode, http.StatusOK)
}

func TestNewCheckValues(t *testing.T) {
	t.Parallel()

	h := NewCheck(OptionsLivePath("livez"),
		OptionsReadyPath("readyz"), OptionsPort("8081"))
	h.Ready()

	r := httptest.NewRequest(http.MethodGet, "/livez", nil)
	w := httptest.NewRecorder()
	h.router().ServeHTTP(w, r)
	h.liveHandler(w, r)
	res := w.Result()

	is := is.New(t)
	is.Equal(res.StatusCode, http.StatusOK)

	r = httptest.NewRequest(http.MethodGet, "/readyz", nil)
	w = httptest.NewRecorder()
	h.router().ServeHTTP(w, r)
	h.readyHandler(w, r)
	res = w.Result()

	is = is.New(t)
	is.Equal(res.StatusCode, http.StatusOK)
}

func TestNewCheckPrefixes(t *testing.T) {
	t.Parallel()

	h := NewCheck(OptionsLivePath("/livez"),
		OptionsReadyPath("/readyz"), OptionsPort(":8082"))
	h.Ready()

	r := httptest.NewRequest(http.MethodGet, "/livez", nil)
	w := httptest.NewRecorder()
	h.router().ServeHTTP(w, r)
	h.liveHandler(w, r)
	res := w.Result()

	is := is.New(t)
	is.Equal(res.StatusCode, http.StatusOK)

	r = httptest.NewRequest(http.MethodGet, "/readyz", nil)
	w = httptest.NewRecorder()
	h.router().ServeHTTP(w, r)
	h.readyHandler(w, r)
	res = w.Result()

	is = is.New(t)
	is.Equal(res.StatusCode, http.StatusOK)
}

func TestLiveness(t *testing.T) {
	t.Parallel()

	h := NewCheck(OptionsLivePath("livez"), OptionsPort("8086"))

	r := httptest.NewRequest(http.MethodGet, "/livez", nil)
	w := httptest.NewRecorder()
	h.router().ServeHTTP(w, r)
	h.liveHandler(w, r)
	res := w.Result()

	is := is.New(t)
	is.Equal(res.StatusCode, http.StatusOK)
}

func TestReadiness(t *testing.T) {
	t.Parallel()

	h := NewCheck(OptionsReadyPath("readyz"), OptionsPort("8087"))

	// test ready
	h.Ready()

	r := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	w := httptest.NewRecorder()
	h.router().ServeHTTP(w, r)
	h.readyHandler(w, r)
	res := w.Result()

	is := is.New(t)
	is.Equal(res.StatusCode, http.StatusOK)

	r = httptest.NewRequest(http.MethodPost, "/readyz", nil)
	w = httptest.NewRecorder()
	h.router().ServeHTTP(w, r)
	h.readyHandler(w, r)
	res = w.Result()

	is = is.New(t)
	is.Equal(res.StatusCode, http.StatusMethodNotAllowed)

	// test notready
	h.NotReady()

	r = httptest.NewRequest(http.MethodGet, "/readyz", nil)
	w = httptest.NewRecorder()
	h.router().ServeHTTP(w, r)
	h.readyHandler(w, r)
	res = w.Result()

	is = is.New(t)
	is.Equal(res.StatusCode, http.StatusServiceUnavailable)
}

func TestTerminating(t *testing.T) {
	t.Parallel()

	h := NewCheck()

	go func() {
		time.Sleep(250 * time.Millisecond)
		proc, err := os.FindProcess(os.Getpid())
		if err != nil {
			t.Error(err)
		}

		proc.Signal(syscall.SIGINT)
	}()

	term := h.Terminating()

	is := is.New(t)
	is.True(term)
}
