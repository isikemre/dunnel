package network

import (
	"net/http"
	"time"

	"github.com/tv42/httpunix"
)

func unixSocketClient(socketPath string) *http.Client {
	// This example shows handling all net/http requests for the
	// http+unix URL scheme.
	u := &httpunix.Transport{
		DialTimeout:           2 * time.Minute,
		RequestTimeout:        2 * time.Minute,
		ResponseHeaderTimeout: 2 * time.Minute,
	}
	u.RegisterLocation("docker", socketPath)

	// If you want to use http: with the same client:
	t := &http.Transport{}
	t.RegisterProtocol(httpunix.Scheme, u)
	var client = http.Client{
		Transport: t,
		Timeout: 5 * time.Minute,
	}

	return &client
}
