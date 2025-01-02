package huaweihilink

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	"github.com/dylanmazurek/huawei-wifi/pkg/huawei-hilink/constants"
	"github.com/dylanmazurek/huawei-wifi/pkg/huawei-hilink/models"
)

func (c *Client) Status() (*models.LoginStateResponse, error) {
	urlStr := fmt.Sprintf("%s/%s/%s", constants.BASE_URL, constants.API_PATH, "monitoring/status")
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
