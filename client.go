package httpdigest

import (
	"io"
	"net/http"
	"net/url"
	"strings"
)

type AuthHandler interface {
	HandleAuth(resp *http.Response, req *http.Request)
}

// A Client wraps http.Client as internal variable while handling HTTP
// authentication challenges and delegate to an appropriated handler.
type Client struct {
	HttpClient  http.Client
	AuthHandler AuthHandler
}

func (c *Client) Do(req *http.Request) (resp *http.Response, err error) {
	resp, err = c.HttpClient.Do(req)

	if err != nil {
		return resp, err
	}

	if resp.StatusCode == 401 && c.AuthHandler != nil {
		c.AuthHandler.HandleAuth(resp, req)

		resp, err = c.HttpClient.Do(req)
	}

	return resp, err
}

func (c *Client) Get(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	return c.Do(req)
}

func (c *Client) Head(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("HEAD", url, nil)

	if err != nil {
		return nil, err
	}

	return c.Do(req)
}

func (c *Client) Post(url string, bodyType string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", url, body)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", bodyType)

	return c.Do(req)
}

func (c *Client) PostForm(url string, data url.Values) (resp *http.Response, err error) {
	return c.Post(url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
}
