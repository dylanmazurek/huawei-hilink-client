package huaweihilink

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"slices"

	"github.com/dylanmazurek/huawei-hilink-client/pkg/huawei-hilink/models"
	"github.com/rs/zerolog/log"
)

func (c *Client) newSession() (*string, error) {
	urlStr := fmt.Sprintf("http://%s", c.session.Host)
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

	cookies := resp.Cookies()
	sessionIdIdx := slices.IndexFunc(cookies, func(cookie *http.Cookie) bool {
		return cookie.Name == "SessionID"
	})

	if sessionIdIdx == -1 {
		return nil, fmt.Errorf("session id not found")
	}

	sessionId := cookies[sessionIdIdx]

	return &sessionId.Value, err
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

	var session *models.Session
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
