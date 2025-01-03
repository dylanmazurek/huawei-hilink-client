package huaweihilink

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/dylanmazurek/huawei-hilink-client/pkg/huawei-hilink/constants"
	"github.com/dylanmazurek/huawei-hilink-client/pkg/huawei-hilink/crypto"
)

type Client struct {
	internalClient *http.Client

	session *Session
	scram   *crypto.Scram
}

func New(ctx context.Context) (*Client, error) {
	host := os.Getenv("HOST")
	if host == "" {
		host = constants.DEFAULT_HOST
	}

	username := os.Getenv("USERNAME")
	if username == "" {
		return nil, fmt.Errorf("missing USERNAME env var")
	}

	newServiceClient := &Client{
		internalClient: http.DefaultClient,

		session: &Session{
			Host:     host,
			Username: username,

			LoggedIn: false,
		},
	}

	password := os.Getenv("PASSWORD")
	scramOpts := []crypto.Option{
		crypto.WithPassword(password),
	}

	scram, err := crypto.NewScram(scramOpts...)
	if err != nil {
		return nil, err
	}

	newServiceClient.scram = scram

	return newServiceClient, nil
}

func (c *Client) createAuthTransport() (*http.Client, error) {
	authClient := &http.Client{
		Transport: &addAuthHeaderTransport{
			T:       http.DefaultTransport,
			Session: c.session,
		},
	}

	return authClient, nil
}

func (c *Client) newRequest(path string, method string, body any) (*http.Request, error) {
	urlStr := fmt.Sprintf("http://%s/%s/%s", c.session.Host, constants.API_PATH, path)
	req, err := http.NewRequest(method, urlStr, nil)
	if err != nil {
		return nil, err
	}

	if c.session.Token != "" {
		req.Header.Add("_ResponseSource", "Broswer")
		req.Header.Add("DNT", "1")

		req.AddCookie(
			&http.Cookie{
				Name:     "SessionID",
				Value:    c.session.SessionId,
				Path:     "/",
				HttpOnly: true,
			},
		)
	}

	if body != nil {
		w := &bytes.Buffer{}
		xmlHeader := `<?xml version: "1.0" encoding="UTF-8"?>`

		w.WriteString(xmlHeader)
		xml.NewEncoder(w).Encode(body)

		bodyBytes := w.Bytes()
		req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	}

	return req, nil
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	resp, err := c.internalClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) parseResponse(resp *http.Response, respObj any) error {
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	byteValue, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = xml.Unmarshal(byteValue, respObj)
	if err != nil {
		return err
	}

	return nil
}
