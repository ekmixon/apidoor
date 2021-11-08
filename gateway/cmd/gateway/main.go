package main

import (
	"encoding/csv"
	"github.com/future-architect/apidoor/gateway"
	"github.com/future-architect/apidoor/gateway/logger"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"os"
)

func main() {

	// open log file
	file, err := os.OpenFile(os.Getenv("LOG_PATH"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0200)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer file.Close()

	// write to log file
	writer := csv.NewWriter(file)
	defer writer.Flush()

	h := gateway.DefaultHandler{
		Appender: logger.CSVAppender{
			Writer: writer,
		},
	}

	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Get("/*", h.Handle)
		r.Put("/*", h.Handle)
		r.Delete("/*", h.Handle)
		r.Post("/*", h.Handle)
	})

	s := &http.Server{
		Addr:    ":3000",
		Handler: r,
	}

	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
