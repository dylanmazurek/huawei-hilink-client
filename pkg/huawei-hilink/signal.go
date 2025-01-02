package huaweihilink

import (
	"fmt"
	"io"
	"net/http"

	"github.com/dylanmazurek/huawei-wifi/pkg/huawei-hilink/constants"
	"github.com/rs/zerolog/log"
)

func (c *Client) netMode() (*string, error) {
	urlStr := fmt.Sprintf("%s/%s/%s", constants.BASE_URL, constants.API_PATH, "net/net-mode")
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

func (c *Client) SignalInfo() (*string, error) {
	// _, err := c.netMode()
	// if err != nil {
	// 	return nil, err
	// }

	urlStr := fmt.Sprintf("%s/%s/%s", constants.BASE_URL, constants.API_PATH, "device/signal")
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
	body := string(byteValue)
	log.Info().Msgf("SignalInfo: %s", body)

	return nil, err
}
