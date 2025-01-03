package huaweihilink

import (
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"

	"github.com/dylanmazurek/huawei-hilink-client/pkg/huawei-hilink/constants"
	"github.com/dylanmazurek/huawei-hilink-client/pkg/huawei-hilink/models"
	"github.com/rs/zerolog/log"
)

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

func (c *Client) checkSession() (*int, error) {
	if c.session.LoggedIn {
		return nil, nil
	}

	req, err := c.newRequest("user/heartbeat", http.MethodGet, nil)

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}

	var response models.HeartbeatResp
	c.parseResponse(resp, &response)

	return &response.Userlevel, err
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

	err = c.newSession()
	if err != nil {
		return err
	}

	err = c.newTokenInfo()
	if err != nil {
		return err
	}

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

func (c *Client) newSession() error {
	urlStr := fmt.Sprintf("http://%s", c.session.Host)
	req, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	cookies := resp.Cookies()
	sessionIdIdx := slices.IndexFunc(cookies, func(cookie *http.Cookie) bool {
		return cookie.Name == "SessionID"
	})

	if sessionIdIdx == -1 {
		return fmt.Errorf("session id not found")
	}

	sessionId := cookies[sessionIdIdx]
	c.session.SessionId = sessionId.Value

	return err
}

func (c *Client) newTokenInfo() error {
	urlStr := fmt.Sprintf("http://%s/%s/%s", c.session.Host, constants.API_PATH, "webserver/SesTokInfo")
	req, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	byteValue, _ := io.ReadAll(resp.Body)

	var sessionTokenInfo models.SessionTokenInfo
	err = xml.Unmarshal(byteValue, &sessionTokenInfo)

	c.session.Token = sessionTokenInfo.TokenInfo

	return err
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

	token2 := resp.Header.Get("__requestverificationtoken")
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

	cookies := resp.Cookies()
	sessionIdIdx := slices.IndexFunc(cookies, func(cookie *http.Cookie) bool {
		return cookie.Name == "SessionID"
	})

	if sessionIdIdx == -1 {
		return fmt.Errorf("session id not found")
	}

	var response models.AuthenticationLoginResp
	c.parseResponse(resp, &response)

	newSession := &Session{
		Host:      c.session.Host,
		Username:  c.session.Username,
		LoggedIn:  true,
		SessionId: cookies[sessionIdIdx].Value,
		Token:     resp.Header.Get("__requestverificationtokenone"),
		Token2:    resp.Header.Get("__requestverificationtokentwo"),

		PublicKey: PublicKey{
			Rsan: response.RsaN,
			Rsae: response.RsaE,
		},
	}

	c.session = newSession

	return err
}

func (c *Client) saveSession() error {
	session := c.session
	sessionJson, err := json.MarshalIndent(session, "", "  ")
	if err != nil {
		return err
	}

	sessionFile, err := os.OpenFile("session.json", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer sessionFile.Close()
	_, err = sessionFile.Write(sessionJson)

	log.Trace().Msg("session saved")

	return err
}

func (c *Client) loadSession() error {
	sessionFile, err := os.Open("session.json")
	if err != nil {
		if os.IsNotExist(err) {
			return ErrSessionFileNotFound
		}

		return err
	}

	defer sessionFile.Close()

	session := &Session{}
	err = json.NewDecoder(sessionFile).Decode(session)
	if err != nil {
		return err
	}

	c.session = session

	userLevel, err := c.checkSession()
	if err != nil {
		return err
	}

	if userLevel == nil || *userLevel < 1 {
		return ErrSessionExpired
	}

	return nil
}
