package client

import (
	"crypto/tls"
	"net/http"
)

type HttpClient struct {
	httpClient *http.Client
}

// NewHttpClient creates http client wrapper with helper functions for rest api call
func NewHttpClient() (*http.Client, error) {
	client := &http.Client{Transport:  &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
	}
	return client, nil
}
