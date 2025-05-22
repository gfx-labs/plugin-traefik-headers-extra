// Package plugindemo a demo plugin.
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/http-wasm/http-wasm-guest-tinygo/handler"
	"github.com/http-wasm/http-wasm-guest-tinygo/handler/api"
)

func main() {
	var config Config
	err := json.Unmarshal(handler.Host.GetConfig(), &config)
	if err != nil {
		handler.Host.Log(api.LogLevelError, fmt.Sprintf("Could not load config %v", err))
		os.Exit(1)
	}

	mw, err := New(config)
	if err != nil {
		handler.Host.Log(api.LogLevelError, fmt.Sprintf("Could not load config %v", err))
		os.Exit(1)
	}
	handler.HandleRequestFn = mw.handleRequest
}

// Config the plugin configuration.
type Config struct {
	Copy []CopyHeader `json:"copy,omitempty"`
}

type CopyHeader struct {
	From string `json:"from,omitempty"`
	To   string `json:"to,omitempty"`
}

// HeadersExtra a HeadersExtra plugin.
type HeadersExtra struct {
	config Config
}

// New created a new HeadersExtra plugin.
func New(config Config) (*HeadersExtra, error) {
	return &HeadersExtra{
		config: config,
	}, nil
}

func (a *HeadersExtra) handleRequest(req api.Request, resp api.Response) (next bool, reqCtx uint32) {
	for _, copyInst := range a.config.Copy {
		if val, ok := req.Headers().Get(copyInst.From); ok {
			req.Headers().Set(copyInst.To, val)
		}
	}
	return true, 0
}
