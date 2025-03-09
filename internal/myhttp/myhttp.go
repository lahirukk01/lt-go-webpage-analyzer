package myhttp

import (
	"lt-app/internal/constants"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

type HTTPClient interface {
	Get(url string) (*resty.Response, error)
}

type RestyClient struct {
	client *resty.Client
}

func NewRestyClient() *RestyClient {
	client := resty.New().SetTimeout(constants.REQUEST_TIMEOUT_SECONDS * time.Second)
	client.SetDoNotParseResponse(true)
	client.SetContentLength(true)
	return &RestyClient{client: client}
}

func (c *RestyClient) Get(url string) (*resty.Response, error) {
	resp, err := c.client.R().Get(url)
	return resp, err
}

func (r *RestyClient) SetTransport(transport http.RoundTripper) {
	r.client.SetTransport(transport)
}
