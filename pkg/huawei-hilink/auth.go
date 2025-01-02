package huaweihilink

import (
	"bytes"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"slices"

	"github.com/dylanmazurek/huawei-wifi/pkg/huawei-hilink/constants"
	"github.com/dylanmazurek/huawei-wifi/pkg/huawei-hilink/models"
	"github.com/rs/zerolog/log"
)

func (c *Client) Login() error {
	err := c.newSession()
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

	c.internalClient, err = c.createAuthTransport()

	return nil
}

func (c *Client) newSession() error {
	urlStr := fmt.Sprintf("%s", constants.BASE_URL)
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
	urlStr := fmt.Sprintf("%s/%s/%s", constants.BASE_URL, constants.API_PATH, "webserver/SesTokInfo")
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

	c.session.Token = sessionTokenInfo.TokInfo

	return err
}

func (c *Client) getToken() error {
	urlStr := fmt.Sprintf("%s/%s/%s", constants.BASE_URL, constants.API_PATH, "webserver/token")
	req, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return err
	}

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
	urlStr := fmt.Sprintf("%s/%s/%s", constants.BASE_URL, constants.API_PATH, "user/challenge_login")
	nonce, err := c.scram.GetNonce()
	if err != nil {
		return err
	}

	reqBody := models.TokenChallengeReq{
		Username:   "admin",
		FirstNonce: hex.EncodeToString(nonce),
		Mode:       1,
	}

	w := &bytes.Buffer{}
	xmlHeader := `<?xml version: "1.0" encoding="UTF-8"?>`

	w.WriteString(xmlHeader)
	xml.NewEncoder(w).Encode(reqBody)

	req, err := http.NewRequest(http.MethodPost, urlStr, bytes.NewBuffer(w.Bytes()))
	if err != nil {
		return err
	}

	req.AddCookie(
		&http.Cookie{
			Name:     "SessionID",
			Value:    c.session.SessionId,
			Path:     "/",
			HttpOnly: true,
		},
	)

	req.Header.Add("__RequestVerificationToken", c.session.Token)
	req.Header.Add("_ResponseSource", "Broswer")
	req.Header.Add("DNT", "1")

	resp, err := c.internalClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	byteValue, _ := io.ReadAll(resp.Body)

	var response models.ChallengeResp
	err = xml.Unmarshal(byteValue, &response)

	token2 := resp.Header.Get("__requestverificationtoken")
	if token2 == "" {
		return fmt.Errorf("token2 not found")
	}

	c.session.Token2 = token2

	c.scram.SetNonce(response.ServerNonce)
	c.scram.SetSalt(response.Salt)

	return err
}

func (c *Client) authenticateToken() error {
	clientProof, err := c.scram.ClientProof()
	if err != nil {
		return err
	}

	urlStr := fmt.Sprintf("%s/%s/%s", constants.BASE_URL, constants.API_PATH, "user/authentication_login")
	nonce, err := c.scram.GetNonce()
	if err != nil {
		return err
	}

	reqBody := models.ClientProofReq{
		ClientProof: hex.EncodeToString(clientProof),
		FinalNonce:  string(nonce),
	}

	w := &bytes.Buffer{}
	xmlHeader := `<?xml version: "1.0" encoding="UTF-8"?>`

	w.WriteString(xmlHeader)
	xml.NewEncoder(w).Encode(reqBody)

	req, err := http.NewRequest(http.MethodPost, urlStr, bytes.NewBuffer(w.Bytes()))
	if err != nil {
		return err
	}

	req.AddCookie(
		&http.Cookie{
			Name:     "SessionID",
			Value:    c.session.SessionId,
			Path:     "/",
			HttpOnly: true,
		},
	)

	req.Header.Add("__RequestVerificationToken", c.session.Token2)
	req.Header.Add("_ResponseSource", "Broswer")
	req.Header.Add("DNT", "1")

	resp, err := c.internalClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	byteValue, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	log.Info().Msgf("%s", byteValue)

	cookies := resp.Cookies()
	sessionIdIdx := slices.IndexFunc(cookies, func(cookie *http.Cookie) bool {
		return cookie.Name == "SessionID"
	})

	if sessionIdIdx == -1 {
		return fmt.Errorf("session id not found")
	}

	var authLoginResp models.AuthenticationLoginResp
	err = xml.Unmarshal(byteValue, &authLoginResp)

	newSession := &Session{
		LoggedIn:  true,
		SessionId: cookies[sessionIdIdx].Value,
		Token:     resp.Header.Get("__requestverificationtokenone"),
		Token2:    resp.Header.Get("__requestverificationtokentwo"),

		PublicKey: PublicKey{
			rsan: authLoginResp.RsaN,
			rsae: authLoginResp.RsaE,
		},
	}

	c.session = newSession

	return err
}

func (c *Client) PrintDebug() {
	fmt.Printf("SessionID: %s\n", c.session.SessionId)

	finalNonce, _ := c.scram.GetNonce()
	fmt.Printf("FinalNonce: %s\n", finalNonce)
	fmt.Printf("Token: %s\n", c.session.Token)
	fmt.Printf("Token2: %s\n", c.session.Token2)
}
