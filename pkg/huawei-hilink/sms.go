package huaweihilink

import (
	"bytes"
	"io"
	"math"
	"net/http"
	"time"

	"github.com/dylanmazurek/huawei-hilink-client/pkg/huawei-hilink/crypto"
	"github.com/dylanmazurek/huawei-hilink-client/pkg/huawei-hilink/models"
	"github.com/rs/zerolog/log"
)

func (c *Client) GetSMSCount() (*int, error) {
	req, err := c.newRequest("sms/sms-count", http.MethodGet, nil)

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}

	smsCountResp := &models.SMSCountResp{}
	err = c.parseResponse(resp, smsCountResp)
	if err != nil {
		return nil, err
	}

	sumCountFloat := float64(smsCountResp.LocalInbox + smsCountResp.LocalOutbox)
	var number = int(math.Ceil(sumCountFloat / 21))

	return &number, err
}

func (c *Client) GetSMSPages() (*models.SMSListResp, error) {
	totalNonce, publicKey, err := c.NewRSA()
	encryptedSmsNonce, err := crypto.RSAEncrypt(*totalNonce, publicKey.Rsan)
	if err != nil {
		return nil, err
	}

	reqBody := models.SMSListReq{
		Phone:     "+61482053833",
		PageIndex: 1,
		ReadCount: 20,
		Nonce:     *encryptedSmsNonce,
	}

	newRequest, err := c.newRequest("sms/sms-list-phone", http.MethodPost, reqBody)

	bodyBytes, err := io.ReadAll(newRequest.Body)
	if err != nil {
		return nil, err
	}

	bodyStr := string(bodyBytes)
	encryptedData, err := crypto.RSAEncrypt(bodyStr, c.session.PublicKey.Rsan)
	if err != nil {
		return nil, err
	}

	newRequest.Body = io.NopCloser(bytes.NewBuffer([]byte(*encryptedData)))

	resp, err := c.do(newRequest)
	if err != nil {
		return nil, err
	}

	smsListResp := &models.SMSListResp{}
	err = c.parseResponse(resp, smsListResp)
	if err != nil {
		return nil, err
	}

	return smsListResp, err
}

func (c *Client) PostSMS(number string, message string) (*int, error) {
	totalNonce, publicKey, err := c.NewRSA()
	encryptedSmsNonce, err := crypto.RSAEncrypt(*totalNonce, publicKey.Rsan)
	if err != nil {
		return nil, err
	}

	reqBody := models.SMSSendReq{
		Index: -1,
		Phones: []models.Phone{
			{Phone: number},
		},
		Content:  message,
		Length:   len(message),
		Reserved: 1,
		Date:     time.Now().Format(time.DateTime),
		Nonce:    *encryptedSmsNonce,
	}

	newRequest, err := c.newRequest("sms/send-sms", http.MethodPost, reqBody)

	bodyBytes, err := io.ReadAll(newRequest.Body)
	if err != nil {
		return nil, err
	}

	bodyStr := string(bodyBytes)
	encryptedData, err := crypto.RSAEncrypt(bodyStr, c.session.PublicKey.Rsan)
	if err != nil {
		return nil, err
	}

	newRequest.Body = io.NopCloser(bytes.NewBuffer([]byte(*encryptedData)))

	resp, err := c.do(newRequest)
	if err != nil {
		return nil, err
	}

	log.Info().Msgf("%s", resp.body)

	return nil, err
}
