package huaweihilink

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	"github.com/dylanmazurek/huawei-wifi/pkg/huawei-hilink/constants"
	"github.com/dylanmazurek/huawei-wifi/pkg/huawei-hilink/models"
)

func (c *Client) getLoginState() (*models.LoginStateResponse, error) {
	urlStr := fmt.Sprintf("%s/%s/%s", constants.BASE_URL, constants.API_PATH, "user/state-login")
	req, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.internalClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	byteValue, _ := io.ReadAll(resp.Body)

	var loginState models.LoginStateResponse
	err = xml.Unmarshal(byteValue, &loginState)

	return &loginState, err
}

// func (c *Client) postLoginChallenge() error {
// 	tokenResp, err := c.getToken()
// 	if err != nil {
// 		panic(err)
// 	}

// 	c.token = *tokenResp

// 	urlStr := fmt.Sprintf("%s/%s/%s", constants.BASE_URL, constants.API_PATH, "/user/challenge_login")
// 	req, err := http.NewRequest(http.MethodPost, urlStr, nil)
// 	if err != nil {
// 		return err
// 	}

// 	req.Header.Add("Cookie", fmt.Sprintf("SessionID=%s", c.sessionId))

// 	return nil
// }
