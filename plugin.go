package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/http-wasm/http-wasm-guest-tinygo/handler"
	"github.com/http-wasm/http-wasm-guest-tinygo/handler/api"
)

func main() {
	var config Config
	if err := json.Unmarshal(handler.Host.GetConfig(), &config); err != nil {
		handler.Host.Log(api.LogLevelError, fmt.Sprintf("could not parse config: %v", err))
		os.Exit(1)
	}

	mw, err := New(config)
	if err != nil {
		handler.Host.Log(api.LogLevelError, fmt.Sprintf("invalid config: %v", err))
		os.Exit(1)
	}
	handler.HandleRequestFn = mw.handleRequest
}

// Config the plugin configuration.
type Config struct {
	Copy []CopyHeader `json:"copy,omitempty"`
}

// CopyHeader copies all values of request header From into request header To.
type CopyHeader struct {
	From string `json:"from,omitempty"`
	To   string `json:"to,omitempty"`
}

// HeadersExtra a HeadersExtra plugin.
type HeadersExtra struct {
	config Config
}

// New validates the config and creates a new HeadersExtra plugin.
func New(config Config) (*HeadersExtra, error) {
	for i, c := range config.Copy {
		if strings.TrimSpace(c.From) == "" || strings.TrimSpace(c.To) == "" {
			return nil, fmt.Errorf("copy[%d]: from and to must be set", i)
		}
		if strings.EqualFold(c.From, c.To) {
			return nil, fmt.Errorf("copy[%d]: from and to must differ", i)
		}
	}
	return &HeadersExtra{config: config}, nil
}

func (a *HeadersExtra) handleRequest(req api.Request, _ api.Response) (next bool, reqCtx uint32) {
	headers := req.Headers()
	for _, c := range a.config.Copy {
		values := headers.GetAll(c.From)
		// Remove first so a client-supplied destination header never survives
		// when the source header is absent.
		headers.Remove(c.To)
		for _, v := range values {
			headers.Add(c.To, v)
		}
	}
	return true, 0
}
