package huaweihilink

import (
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"slices"

	"github.com/dylanmazurek/huawei-hilink-client/pkg/huawei-hilink/constants"
	"github.com/dylanmazurek/huawei-hilink-client/pkg/huawei-hilink/crypto"
	"github.com/dylanmazurek/huawei-hilink-client/pkg/huawei-hilink/models"
	"github.com/rs/zerolog/log"
)

type addAuthHeaderTransport struct {
	T       http.RoundTripper
	Session *models.Session
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

func (c *Client) Login() error {
	err := c.loadSession()
	if err == nil {
		c.internalClient, err = c.createAuthTransport()
		if err == nil {
			return nil
		}

		log.Info().Msg("session loaded but failed to create auth transport")
	}

	sessionId, err := c.newSession()
	if err != nil {
		return err
	}

	c.session.SessionId = *sessionId

	token, err := getTokenInfo(c.session.Host)
	if err != nil {
		return err
	}

	c.session.Token = *token

	err = c.getToken()
	if err != nil {
		return err
	}

	err = c.challengeToken()
	if err != nil {
		return err
	}

	err = c.authenticateToken()
	if err != nil {
		return err
	}

	err = c.saveSession()
	if err != nil {
		return err
	}

	c.internalClient, err = c.createAuthTransport()

	return nil
}

func (c *Client) getToken() error {
	urlStr := fmt.Sprintf("http://%s/%s/%s", c.session.Host, constants.API_PATH, "webserver/token")
	req, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return err
	}

	resp, err := c.internalClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	byteValue, _ := io.ReadAll(resp.Body)

	var tokenResponse models.TokenResp
	err = xml.Unmarshal(byteValue, &tokenResponse)

	return err
}

func (c *Client) challengeToken() error {
	nonce, err := c.scram.GetNonce()
	if err != nil {
		return err
	}

	reqBody := models.TokenChallengeReq{
		Username:   c.session.Username,
		FirstNonce: hex.EncodeToString(nonce),
		Mode:       1,
	}

	req, err := c.newRequest("user/challenge_login", http.MethodPost, reqBody)

	req.Header.Add("__RequestVerificationToken", c.session.Token)

	resp, err := c.do(req)
	if err != nil {
		return err
	}

	token2 := resp.header.Get("__requestverificationtoken")
	if token2 == "" {
		return fmt.Errorf("token2 not found")
	}

	c.session.Token2 = token2

	var response models.ChallengeResp
	c.parseResponse(resp, &response)

	c.scram.SetNonce(response.ServerNonce)
	c.scram.SetSalt(response.Salt)

	return err
}

func (c *Client) authenticateToken() error {
	clientProof, err := c.scram.ClientProof()
	if err != nil {
		return err
	}

	nonce, err := c.scram.GetNonce()
	if err != nil {
		return err
	}

	reqBody := models.ClientProofReq{
		ClientProof: hex.EncodeToString(clientProof),
		FinalNonce:  string(nonce),
	}

	req, err := c.newRequest("user/authentication_login", http.MethodPost, reqBody)

	req.Header.Add("__RequestVerificationToken", c.session.Token2)

	resp, err := c.do(req)
	if err != nil {
		return err
	}

	cookies := resp.cookies
	sessionIdIdx := slices.IndexFunc(cookies, func(cookie *http.Cookie) bool {
		return cookie.Name == "SessionID"
	})

	if sessionIdIdx == -1 {
		return fmt.Errorf("session id not found")
	}

	var response models.AuthenticationLoginResp
	c.parseResponse(resp, &response)

	newSession := &models.Session{
		Host:      c.session.Host,
		Username:  c.session.Username,
		LoggedIn:  true,
		SessionId: cookies[sessionIdIdx].Value,
		Token:     resp.header.Get("__requestverificationtokenone"),
		Token2:    resp.header.Get("__requestverificationtokentwo"),

		PublicKey: models.PublicKey{
			Rsan: response.RsaN,
			Rsae: response.RsaE,
		},
	}

	c.session = newSession

	return err
}

// unauthenticated methods
func getTokenInfo(host string) (*string, error) {
	urlStr := fmt.Sprintf("http://%s/%s/%s", host, constants.API_PATH, "webserver/SesTokInfo")
	req, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	byteValue, _ := io.ReadAll(resp.Body)

	var sessionTokenInfo models.SessionTokenInfo
	err = xml.Unmarshal(byteValue, &sessionTokenInfo)

	return &sessionTokenInfo.TokenInfo, err
}

func getPublicKey(host string, sessionId string, token string) (*models.PublicKey, error) {
	urlStr := fmt.Sprintf("http://%s/%s/%s", host, constants.API_PATH, "webserver/publickey")
	req, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("__RequestVerificationToken", token)

	req.AddCookie(
		&http.Cookie{
			Name:     "SessionID",
			Value:    sessionId,
			Path:     "/",
			HttpOnly: true,
		},
	)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	byteValue, _ := io.ReadAll(resp.Body)

	var publicKey models.PublicKey
	err = xml.Unmarshal(byteValue, &publicKey)

	return &publicKey, nil
}

func (c *Client) NewRSA() (*string, *models.PublicKey, error) {
	sessionId, err := c.newSession()
	if err != nil {
		return nil, nil, err
	}

	token, err := getTokenInfo(c.session.Host)

	scramOpts := []crypto.Option{}

	smsScram, err := crypto.NewScram(scramOpts...)
	if err != nil {
		return nil, nil, err
	}

	smsNonce, err := smsScram.GetNonce()
	if err != nil {
		return nil, nil, err
	}

	smsSalt, err := smsScram.GetNonce()
	if err != nil {
		return nil, nil, err
	}

	totalNonce := fmt.Sprintf("%s%s", hex.EncodeToString(smsNonce), hex.EncodeToString(smsSalt))
	smsPublicKey, err := getPublicKey(c.session.Host, *sessionId, *token)
	if err != nil {
		return nil, nil, err
	}

	return &totalNonce, smsPublicKey, nil
}
