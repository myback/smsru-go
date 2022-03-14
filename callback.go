package smsru

import (
	"encoding/json"
	"net/url"
)

const (
	smsRuCallbackAddPath    = "/callback/add"
	smsRuCallbackDeletePath = "/callback/del"
	smsRuCallbackListPath   = "/callback/get"
)

type ResponseCallback struct {
	Callback   []string `json:"callback"`
	Status     string   `json:"status"`
	StatusCode int      `json:"status_code"`
}

type Callback struct {
	req request
}

func (sms *SmsRu) NewCallback() *Callback {
	return &Callback{}
}

func (c *Callback) request(p, u string) (*ResponseCallback, error) {
	params := url.Values{}
	if len(u) > 0 {
		params.Set("url", u)
	}

	resp, err := c.req(p, params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data := new(ResponseCallback)
	if err = json.NewDecoder(resp.Body).Decode(data); err != nil {
		return nil, err
	}

	return data, nil
}

func (c *Callback) Add(url string) (*ResponseCallback, error) {
	return c.request(smsRuCallbackAddPath, url)
}

func (c *Callback) Delete(url string) (*ResponseCallback, error) {
	return c.request(smsRuCallbackDeletePath, url)
}

func (c *Callback) List() (*ResponseCallback, error) {
	return c.request(smsRuCallbackListPath, "")
}
