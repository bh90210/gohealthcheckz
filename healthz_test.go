package healthz

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLiveness(t *testing.T) {
	t.Parallel()

	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080", bytes.NewReader(pl))
	if err != nil {
		t.Fatalf("could not create test request: %v", err)
	}
	rec := httptest.NewRecorder()
	handle(rec, req)
	res := rec.Result()

	if res.StatusCode != http.StatusOK {
		t.Errorf("unexpected status code %s", res.Status)
	}
	defer res.Body.Close()

	msg, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("could not read result payload: %v", err)
	}

	if exp := "pull request id: 191568743"; string(msg) != exp {
		t.Fatalf("expected message %q; got %q", exp, msg)
	}
}

func TestReadiness(t *testing.T) {
	t.Parallel()
}

func TestTerminating(t *testing.T) {
	t.Parallel()
}

// // https://golang.org/src/net/http/httptest/example_test.go
// func TestRouter(t *testing.T) {
// 	t.Parallel()
// 	ts := httptest.NewServer(router())
// 	defer ts.Close()

// 	res, err := http.Get(ts.URL)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	if res.StatusCode != 200 {
// 		t.Errorf("handler returned unexpected status code: got %v want 200",
// 			res.StatusCode)
// 	}
// }

// func TestGet(t *testing.T) {
// 	t.Parallel()
// 	req := httptest.NewRequest("GET", "/", nil)
// 	w := httptest.NewRecorder()
// 	get(w, req)
// 	resp := w.Result()
// 	if status := resp.StatusCode; status != http.StatusOK {
// 		t.Errorf("handler returned wrong status code: got %v want %v",
// 			status, http.StatusOK)
// 	}

// 	if resp.StatusCode != 200 {
// 		t.Errorf("handler returned unexpected status code: got\n %v \nwant\n 200",
// 			resp.StatusCode)
// 	}
// }

// func TestHandle(t *testing.T) {
// 	pl, err := ioutil.ReadFile("testdata/payload.json")
// 	if err != nil {
// 		t.Fatalf("could not read payload.json: %v", err)
// 	}

// 	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080", bytes.NewReader(pl))
// 	if err != nil {
// 		t.Fatalf("could not create test request: %v", err)
// 	}
// 	rec := httptest.NewRecorder()
// 	handle(rec, req)
// 	res := rec.Result()

// 	if res.StatusCode != http.StatusOK {
// 		t.Errorf("unexpected status code %s", res.Status)
// 	}
// 	defer res.Body.Close()

// 	msg, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		t.Fatalf("could not read result payload: %v", err)
// 	}

// 	if exp := "pull request id: 191568743"; string(msg) != exp {
// 		t.Fatalf("expected message %q; got %q", exp, msg)
// 	}
// }

// func BenchmarkHandle(b *testing.B) {
// 	b.StopTimer()

// 	pl, err := ioutil.ReadFile("testdata/payload.json")
// 	if err != nil {
// 		b.Fatalf("could not read payload.json: %v", err)
// 	}

// 	for i := 0; i < b.N; i++ {
// 		req, err := http.NewRequest(http.MethodPost, "http://localhost:8080", bytes.NewReader(pl))
// 		if err != nil {
// 			b.Fatalf("could not create test request: %v", err)
// 		}
// 		rec := httptest.NewRecorder()

// 		b.StartTimer()
// 		handle(rec, req)
// 		b.StopTimer()
// 	}
// }
