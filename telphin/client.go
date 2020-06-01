package telphin

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	mem "github.com/spf13/afero/mem"
)

// NewClient returns new client struct
func NewClient(clientID string, secret string, host string) (*Client, error) {
	if clientID == "" || secret == "" || host == "" {
		return nil, errors.New("ClientID, Secret and Host are required to create a Client")
	}

	return &Client{
		Client: &http.Client{
			Timeout: time.Second * 30,
		},
		ClientID: clientID,
		Secret:   secret,
		Host:     host,
	}, nil
}

// GetAccessToken returns struct of TokenResponse
// No need to call SetAccessToken to apply new access token for current Client
// Endpoint: POST /oauth/token
func (c *Client) GetAccessToken() (*TokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", c.ClientID)
	data.Set("client_secret", c.Secret)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", c.Host, "/oauth/token"), bytes.NewBufferString(data.Encode()))
	if err != nil {
		return &TokenResponse{}, err
	}

	req.Header.Set("Content-type", "application/x-www-form-urlencoded")

	t := TokenResponse{}
	err = c.Send(req, &t)

	// Set Token fur current Client
	if t.Token != "" {
		c.Token = &t
		c.tokenExpiresAt = time.Now().Add(time.Duration(t.ExpiresIn) * time.Second)
	}

	return &t, err
}

// SetHTTPClient sets *http.Client to current client
func (c *Client) SetHTTPClient(client *http.Client) {
	c.Client = client
}

// SetAccessToken sets saved token to current client
func (c *Client) SetAccessToken(token string, expiresIn expirationTime) {
	c.Token = &TokenResponse{
		Token:     token,
		ExpiresIn: expiresIn,
	}
	c.tokenExpiresAt = time.Now().Add(time.Duration(expiresIn) * time.Second)
}

// TODO
func (c *Client) ResetAccessToken() error {
	c.Lock()
	defer c.Unlock()

	// c.Token will be updated in GetAccessToken call
	_, err := c.GetAccessToken()
	return err
}

// TODO
func (c *Client) SetLogger(logger FieldLogger) {
	c.Logger = logger
}

// SendWithAuth makes a request to the API and apply OAuth2 header automatically.
func (c *Client) SendWithAuth(req *http.Request, v interface{}) error {
	c.Lock()

	if c.Token != nil {
		if !c.tokenExpiresAt.IsZero() && c.tokenExpiresAt.Sub(time.Now()) < RequestNewTokenBeforeExpiresIn {
			// c.Token will be updated in GetAccessToken call
			if _, err := c.GetAccessToken(); err != nil {
				c.Unlock()
				return err
			}
		}

		req.Header.Set("Authorization", "Bearer "+c.Token.Token)
	}

	c.Unlock()
	return c.Send(req, v)
}

// Send makes a request to the API, the response body will be
// unmarshaled into v, or if v is an io.Writer, the response will
// be written to it without decoding
func (c *Client) Send(req *http.Request, v interface{}) error {
	var (
		err  error
		resp *http.Response
		// data []byte
	)

	if req.Header.Get("Content-type") == "" {
		req.Header.Set("Content-type", "application/json")
	}

	resp, err = c.Client.Do(req)
	c.log(req, resp)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		c.Lock()
		c.tokenExpiresAt = time.Now()
		c.Unlock()
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		errResp := &ErrorResponse{
			Code: resp.StatusCode,
		}
		data, err := ioutil.ReadAll(resp.Body)

		if err == nil && len(data) > 0 {
			json.Unmarshal(data, errResp)
			errResp.Details = string(data)
		}

		return errResp
	}

	if v == nil {
		return nil
	}

	if f, ok := v.(*mem.File); ok {
		_, params, err := mime.ParseMediaType(resp.Header.Get("Content-Disposition"))
		if err != nil {
			return err
		}
		mem.ChangeFileName(f.Data(), params["filename"])
		_, err = io.Copy(f, resp.Body)
		f.Seek(0, 0)
		return err
	}

	if w, ok := v.(io.Writer); ok {
		_, err = io.Copy(w, resp.Body)
		return err
	}

	return json.NewDecoder(resp.Body).Decode(v)
}

// NewRequest constructs a request
// Convert payload to a JSON
func (c *Client) NewRequest(method, url string, payload interface{}) (*http.Request, error) {
	var buf io.Reader
	if payload != nil {
		var b []byte
		b, err := json.Marshal(&payload)
		if err != nil {
			return nil, err
		}
		buf = bytes.NewBuffer(b)
	}
	return http.NewRequest(method, url, buf)
}

// log will dump request and response
func (c *Client) log(req *http.Request, resp *http.Response) {
	if c.Logger != nil {
		fields := make(map[string]interface{})

		if req != nil {
			reqRaw, _ := httputil.DumpRequest(req, true)
			fields["request"] = map[string]interface{}{
				"method": req.Method,
				"url":    req.URL.String(),
				"raw":    string(reqRaw), // TODO: data is always empty
			}
		}
		if resp != nil {
			switch resp.Header.Get("Content-Type") {
			case "text/html":
			case "application/json":
				respRaw, _ := httputil.DumpResponse(resp, true)
				fields["response"] = map[string]interface{}{
					"code": resp.StatusCode,
					"raw":  string(respRaw),
				}
			default:
				fields["response"] = map[string]interface{}{
					"code": resp.StatusCode,
				}
			}
		}

		c.Logger.WithFields(fields).Info("Telphin request")
	}
}
