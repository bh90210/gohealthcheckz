package healthz

import (
	"net/http"
	"net/http/httptest"
	"os"
	"syscall"
	"testing"
	"time"
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
	if res1.StatusCode != http.StatusOK {
		t.Errorf("handler returned unexpected status code: got %v want 200",
			res1.Status)
	}

	r = httptest.NewRequest(http.MethodGet, "/ready", nil)
	w = httptest.NewRecorder()
	h.router().ServeHTTP(w, r)
	h.readyHandler(w, r)
	res2 := w.Result()
	if res2.StatusCode != http.StatusOK {
		t.Errorf("handler returned unexpected status code: got %v want 200",
			res2.Status)
	}
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
	if res.StatusCode != http.StatusOK {
		t.Errorf("handler returned unexpected status code: got %v want 200",
			res.Status)
	}

	r = httptest.NewRequest(http.MethodGet, "/readyz", nil)
	w = httptest.NewRecorder()
	h.router().ServeHTTP(w, r)
	h.readyHandler(w, r)
	res = w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("handler returned unexpected status code: got %v want 200",
			res.Status)
	}
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
	if res.StatusCode != http.StatusOK {
		t.Errorf("handler returned unexpected status code: got %v want 200",
			res.Status)
	}

	r = httptest.NewRequest(http.MethodGet, "/readyz", nil)
	w = httptest.NewRecorder()
	h.router().ServeHTTP(w, r)
	h.readyHandler(w, r)
	res = w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("handler returned unexpected status code: got %v want 200",
			res.Status)
	}
}

func TestLiveness(t *testing.T) {
	t.Parallel()

	h := NewCheck(OptionsLivePath("livez"), OptionsPort("8086"))

	r := httptest.NewRequest(http.MethodGet, "/livez", nil)
	w := httptest.NewRecorder()
	h.router().ServeHTTP(w, r)
	h.liveHandler(w, r)
	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("handler returned unexpected status code: got %v want 200",
			res.Status)
	}

	r = httptest.NewRequest(http.MethodPost, "/livez", nil)
	w = httptest.NewRecorder()
	h.router().ServeHTTP(w, r)
	h.liveHandler(w, r)
	res = w.Result()
	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("handler returned unexpected status code: got %v want 405",
			res.Status)
	}
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
	if res.StatusCode != http.StatusOK {
		t.Errorf("handler returned unexpected status code: got %v want 200",
			res.Status)
	}

	r = httptest.NewRequest(http.MethodPost, "/readyz", nil)
	w = httptest.NewRecorder()
	h.router().ServeHTTP(w, r)
	h.readyHandler(w, r)
	res = w.Result()
	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("handler returned unexpected status code: got %v want 405",
			res.Status)
	}

	// test notready
	h.NotReady()

	r = httptest.NewRequest(http.MethodGet, "/readyz", nil)
	w = httptest.NewRecorder()
	h.router().ServeHTTP(w, r)
	h.readyHandler(w, r)
	res = w.Result()
	if res.StatusCode != http.StatusServiceUnavailable {
		t.Errorf("handler returned unexpected status code: got %v want 503",
			res.Status)
	}
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
	if term != true {
		t.Errorf("termination return: got %v want true",
			term)
	}
}
