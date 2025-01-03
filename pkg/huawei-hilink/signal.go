package huaweihilink

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	"github.com/dylanmazurek/huawei-hilink-client/pkg/huawei-hilink/constants"
	"github.com/dylanmazurek/huawei-hilink-client/pkg/huawei-hilink/models"
)

func (c *Client) NetMode() (*string, error) {
	urlStr := fmt.Sprintf("http://%s/%s/%s", c.session.Host, constants.API_PATH, "net/net-mode")
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

	return nil, err
}

func (c *Client) SignalInfo() (*models.SignalInfoResp, error) {
	urlStr := fmt.Sprintf("http://%s/%s/%s", c.session.Host, constants.API_PATH, "device/signal")
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

	var signalInfoResp models.SignalInfoResp
	err = xml.Unmarshal(byteValue, &signalInfoResp)

	return &signalInfoResp, err
}
