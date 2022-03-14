package smsru

import (
	"encoding/json"
	"net/url"
)

const smsRuCallUrl = "/code/call"

type ResponseCallData struct {
	Balance float64 `json:"balance"`
	CallID  string  `json:"call_id"`
	Code    string  `json:"code"`
	Cost    float64 `json:"cost"`
	Status  string  `json:"status"`
}

type Call struct {
	req request
}

func (sms *SmsRu) NewCall() *Call {
	return &Call{sms.request}
}

func (c *Call) Send(phone string) (*ResponseCallData, error) {
	params := url.Values{}
	params.Set("phone", phone)

	resp, err := c.req(smsRuCallUrl, params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data := new(ResponseCallData)
	if err = json.NewDecoder(resp.Body).Decode(data); err != nil {
		return nil, err
	}

	return data, nil
}
