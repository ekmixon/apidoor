package gateway

import (
	"log"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		log.Print("unexpected request content")
		http.Error(w, "unexpected request content", http.StatusBadRequest)
		return
	}

	apikey := r.Header.Get("Authorization")
	reqpath := r.URL.Path
	query := r.URL.RawQuery

	path, err := GetAPIURL(r.Context(), apikey, reqpath)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "invalid key or path", http.StatusBadRequest)
		return
	}

	res, err := http.Get("http://" + path + "?" + query)
	if err != nil {
		log.Printf("error in http get: %s", err.Error())
		http.Error(w, "invalid request", http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	if err := ResposeChecker(&w, res); err != nil {
		return
	}

	UpdateLog(apikey, path)
}
