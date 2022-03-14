package smsru

import (
	"encoding/json"
	"errors"
	"net/url"
)

const (
	smsRuCallcheckAddApiPath    = "/callcheck/add"
	smsRuCallcheckStatusApiPath = "/callcheck/status"
)

type ResponseCallCheckData struct {
	CallPhone       string  `json:"call_phone"`
	CallPhoneHTML   string  `json:"call_phone_html"`
	CallPhonePretty string  `json:"call_phone_pretty"`
	CheckID         string  `json:"check_id"`
	Cost            float64 `json:"cost"`
	Status          string  `json:"status"`
	StatusCode      int     `json:"code"`
}

type ResponseCallCheckStatusData struct {
	CheckStatus     string `json:"check_status"`
	CheckStatusText string `json:"check_status_text"`
	Status          string `json:"status"`
	StatusCode      int    `json:"code"`
}

type CallCheck struct {
	checkID string
	req     request
}

type CallCheckStatus struct {
	req request
}

func (sms *SmsRu) NewCallCheck() *CallCheck {
	return &CallCheck{req: sms.request}
}

func (sms *SmsRu) NewCallCheckStatus() *CallCheckStatus {
	return &CallCheckStatus{req: sms.request}
}

func (c *CallCheck) Send(phone string) (*ResponseCallCheckData, error) {
	params := url.Values{}
	params.Set("phone", phone)

	resp, err := c.req(smsRuCallcheckAddApiPath, params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data := new(ResponseCallCheckData)
	if err = json.NewDecoder(resp.Body).Decode(data); err != nil {
		return nil, err
	}

	c.checkID = data.CheckID

	return data, nil
}

func (c *CallCheck) Status() (*ResponseCallCheckStatusData, error) {
	if len(c.checkID) == 0 {
		return nil, errors.New("empty check_id. First call Send()")
	}

	return getStatus(c.req, c.checkID)
}

func (c *CallCheckStatus) Get(checkID string) (*ResponseCallCheckStatusData, error) {
	if len(checkID) == 0 {
		return nil, errors.New("empty check_id")
	}

	return getStatus(c.req, checkID)
}

func getStatus(req request, checkID string) (*ResponseCallCheckStatusData, error) {
	params := url.Values{}
	params.Set("check_id", checkID)

	resp, err := req(smsRuCallcheckStatusApiPath, params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data := new(ResponseCallCheckStatusData)
	if err = json.NewDecoder(resp.Body).Decode(data); err != nil {
		return nil, err
	}

	return data, nil
}
