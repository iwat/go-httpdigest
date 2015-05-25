package httpdigest

import (
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	Client http.Client
}

func (client *Client) Do(req *http.Request) (resp *http.Response, err error) {
	resp, err = client.Client.Do(req)

	if err != nil {
		return resp, err
	}

	username, password, ok := req.BasicAuth()

	if resp.StatusCode == 401 && ok {
		challenge := ChallengeFromResponse(resp)
		challenge.username = username
		challenge.password = password
		challenge.Path = req.URL.RequestURI()
		challenge.ApplyAuth(req)

		resp, err = client.Client.Do(req)
	}

	return resp, err
}

func (c *Client) Get(url string) (resp *http.Response, err error) {
	panic("httpdigest.Client.Get is not implemented")
	return
}

func (c *Client) Head(url string) (resp *http.Response, err error) {
	panic("httpdigest.Client.Head is not implemented")
	return
}

func (c *Client) Post(url string, bodyType string, body io.Reader) (resp *http.Response, err error) {
	panic("httpdigest.Client.Post is not implemented")
	return
}

func (c *Client) PostForm(url string, data url.Values) (resp *http.Response, err error) {
	panic("httpdigest.Client.PostForm is not implemented")
	return
}
