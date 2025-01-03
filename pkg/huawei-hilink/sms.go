package huaweihilink

import (
	"encoding/xml"
	"fmt"
	"io"
	"math"
	"net/http"

	"github.com/dylanmazurek/huawei-hilink-client/pkg/huawei-hilink/constants"
	"github.com/dylanmazurek/huawei-hilink-client/pkg/huawei-hilink/models"
)

func (c *Client) GetSMSCount() (*int, error) {
	urlStr := fmt.Sprintf("http://%s/%s/%s", c.session.Host, constants.API_PATH, "sms/sms-count")
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

	var smsCountResp models.SMSCountResp
	err = xml.Unmarshal(byteValue, &smsCountResp)

	sumCountFloat := float64(smsCountResp.LocalInbox + smsCountResp.LocalOutbox)
	var number = int(math.Ceil(sumCountFloat / 21))

	return &number, err
}
