package apidoor

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

var count int = 0

func Handler(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")
	if r.Header.Get("Content-Type") != "application/json" {
		log.Print("unexpected request content")
		http.Error(w, "unexpected request content", http.StatusBadRequest)
		return
	}

	r.Header.Set("Authorization", "testtoken")
	if r.Header.Get("Authorization") == "" {
		log.Print("unauthorized request")
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	apikey := r.URL.Query().Get("apikey")
	num, err := strconv.Atoi(r.URL.Query().Get("num"))
	if err != nil {
		log.Print(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	url, err := GetAPIURL(num, apikey)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := http.Get(url)
	if err != nil {
		log.Printf("error in http get: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	contents, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("error in io.ReadAll: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	switch code := res.StatusCode; {
	case 400 <= code && code <= 499:
		log.Printf("client error: %v, status code: %d", string(contents), code)
		http.Error(w, string(contents), code)
		return
	case 500 <= code && code <= 599:
		log.Printf("server error: %v, status code: %d", string(contents), code)
		http.Error(w, string(contents), code)
		return
	}

	count++
	log.Printf("called: %d times", count)
	fmt.Fprint(w, string(contents))
}
