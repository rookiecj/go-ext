package httpx

import (
	"io"
	"net/http"
	"net/http/httptest"
)

var testServer *httptest.Server

func startServer() (url string) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//fmt.Fprintln(w, "Hello, client")
		reqBytes, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusOK)
			contentType := r.Header.Get("Content-Type")
			w.Header().Set("Content-Type", contentType)
			w.Write(reqBytes)
		}
	})
	testServer = httptest.NewServer(handler)
	return testServer.URL
}

func stopServer() {
	if testServer != nil {
		testServer.Close()
		testServer = nil
	}
}
