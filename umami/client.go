package umami

import (
	"crypto/tls"
	"net/http"
	"time"

	logger "github.com/rtfmkiesel/kisslog"

	"github.com/rtfmkiesel/umami-importer/config"
)

var log = logger.New("umami")

type Client struct {
	config     *config.UmamiConfig
	httpClient *http.Client
}

func NewClient(config *config.UmamiConfig) *Client {
	transportConfig := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: config.IgnoreTLS,
		},
	}

	client := &Client{
		config: config,
		httpClient: &http.Client{
			Timeout:   time.Duration(config.Timeout) * time.Second,
			Transport: transportConfig,
		},
	}

	return client
}
