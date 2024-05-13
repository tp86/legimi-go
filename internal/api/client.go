package api

import (
	"bytes"
	"io"
	"net/http"

	"github.com/tp86/legimi-go/internal/api/protocol"
)

type Client interface {
	Exchange(protocol.Request, protocol.Response) error
}

var client Client

func GetClient() Client {
	if client == nil {
		client = &defaultClient{}
	}
	return client
}

type defaultClient struct{}

func (c defaultClient) Exchange(request protocol.Request, response protocol.Response) error {
	buf := new(bytes.Buffer)
	err := protocol.Encode(buf, request)
	if err != nil {
		return err
	}
	resp, err := http.Post(protocol.APIUrl, "application/octet-stream", buf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return protocol.Decode(bytes.NewBuffer(body), response)
}
