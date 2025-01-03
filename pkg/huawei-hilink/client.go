package huaweihilink

import (
	"context"
	"net/http"
	"os"

	"github.com/dylanmazurek/huawei-hilink-client/pkg/huawei-hilink/crypto"
)

type Client struct {
	internalClient *http.Client

	session *Session
	scram   *crypto.Scram
}

type addAuthHeaderTransport struct {
	T       http.RoundTripper
	Session *Session
}

func (adt *addAuthHeaderTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8;enc")
	req.Header.Set("__RequestVerificationToken", adt.Session.Token2)
	req.Header.Set("_ResponseSource", "Broswer")

	req.AddCookie(
		&http.Cookie{
			Name:  "SessionID",
			Value: adt.Session.SessionId,
		},
	)

	return adt.T.RoundTrip(req)
}

func New(ctx context.Context) (*Client, error) {
	newServiceClient := &Client{
		internalClient: http.DefaultClient,

		session: &Session{
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
