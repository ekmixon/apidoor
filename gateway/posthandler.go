package gateway

import (
	"io"
	"log"
	"net/http"
)

func PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		log.Print("unexpected request content")
		http.Error(w, "unexpected request content", http.StatusBadRequest)
		return
	}

	apikey := r.Header.Get("Authorization")
	reqpath := r.URL.Path

	path, err := GetAPIURL(r.Context(), apikey, reqpath)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "invalid key or path", http.StatusBadRequest)
		return
	}

	res, err := http.Post("http://"+path, "application/json", r.Body)
	if err != nil {
		log.Printf("error in http post: %s", err.Error())
		http.Error(w, "invalid request", http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	switch code := res.StatusCode; {
	case 400 <= code && code <= 499:
		log.Printf("client error: %v, status code: %d", res.Body, code)
		http.Error(w, "client error", code)
		return
	case 500 <= code && code <= 599:
		log.Printf("server error: %v, status code: %d", res.Body, code)
		http.Error(w, "server error", code)
		return
	}

	UpdateLog(apikey, path)
	if _, err := io.Copy(w, res.Body); err != nil {
		log.Printf("error occur while writing response: %s", err.Error())
		http.Error(w, "error occur while writing response", http.StatusInternalServerError)
	}
}
