package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/tp86/legimi-go/internal/api/protocol"
	"github.com/tp86/legimi-go/internal/debug"
	"github.com/tp86/legimi-go/internal/options"
)

type Client interface {
	Exchange(protocol.Request, protocol.Response) error
}

var client Client

func GetClient(opts options.Debugging) Client {
	if client == nil {
		client = &defaultClient{opts.IsDebug()}
	}
	return client
}

type defaultClient struct {
	debugOn bool
}

func (c defaultClient) Exchange(request protocol.Request, response protocol.Response) error {
	buf := new(bytes.Buffer)
	c.debug(request)
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
	err = protocol.Decode(bytes.NewBuffer(body), response)
	c.debug(response)
	return err
}

func (c defaultClient) debug(value any) {
	if c.debugOn {
		if debugFormatter, ok := value.(debug.Formatter); ok {
			formatted := debugFormatter.DebugFormat()
			fmt.Fprintf(os.Stderr, "=== DEBUG START ===\n%s\n=== DEBUG END ===\n", formatted)
		}
	}
}
