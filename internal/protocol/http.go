package protocol

import (
	"bytes"
	"net/http"
)

func Exchange(request Request, response Response) error {
	buf := new(bytes.Buffer)
	err := Encode(buf, request)
	if err != nil {
		return err
	}
	resp, err := http.Post(APIUrl, "application/octet-stream", buf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return Decode(resp.Body, response)
}
