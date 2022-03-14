package smsru

import (
	"encoding/json"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	smsRuSmsSendPath   = "/sms/send"
	smsRuSmsStatusPath = "/sms/status"
)

type SMS struct {
	params url.Values
	req    request
}

type SMSStatus struct {
	req request
}

type ResponseSMSSent struct {
	Balance    float64                `json:"balance"`
	SMS        map[string]SMSSentInfo `json:"sms"`
	Status     string                 `json:"status"`
	StatusCode int                    `json:"status_code"`
}

type SMSSentInfo struct {
	SmsID      string `json:"sms_id,omitempty"`
	Status     string `json:"status"`
	StatusCode int    `json:"status_code"`
	StatusText string `json:"status_text,omitempty"`
}

type SMSStatuses struct {
	SMS        map[string]SMSStatusesInfo `json:"sms"`
	Status     string                     `json:"status"`
	StatusCode int                        `json:"status_code"`
}

type SMSStatusesInfo struct {
	Cost       float64 `json:"cost"`
	Status     string  `json:"status"`
	StatusCode int     `json:"status_code"`
	StatusText string  `json:"status_text,omitempty"`
}

func (sms *SmsRu) NewSMS() *SMS {
	return &SMS{
		params: url.Values{},
		req:    sms.request,
	}
}

func (sms *SmsRu) NewSMSStatus() *SMSStatus {
	return &SMSStatus{sms.request}
}

func (sms *SMS) Add(msg, to string) *SMS {
	sms.params.Set("to["+to+"]", msg)

	return sms
}

func (sms *SMS) AddClientIP(ip string) *SMS {
	sms.params.Set("ip", ip)

	return sms
}

func (sms *SMS) ConsiderTimezone() *SMS {
	sms.params.Set("daytime", "1")

	return sms
}

func (sms *SMS) DelayedSend(t time.Time) *SMS {
	timeStr := strconv.FormatUint(uint64(t.Unix()), 10)
	sms.params.Set("time", timeStr)

	return sms
}

func (sms *SMS) PartnerID(id string) *SMS {
	sms.params.Set("partner_id", id)

	return sms
}

func (sms *SMS) TTL(minutes int) *SMS {
	if minutes < 1 {
		//warn
		minutes = 1
	} else if minutes > 1440 {
		//warn
		minutes = 1440
	}
	sms.params.Set("ttl", strconv.Itoa(minutes))

	return sms
}

func (sms *SMS) UseTranslit() *SMS {
	sms.params.Set("translit", "1")

	return sms
}

func (sms *SMS) Test() *SMS {
	sms.params.Set("test", "1")

	return sms
}

func (sms *SMS) MultipleRecipients(msg string, to ...string) *SMS {
	sms.params.Set("to", strings.Join(to, ","))
	sms.params.Set("msg", msg)

	return sms
}

func (sms *SMS) Send() (*ResponseSMSSent, error) {
	resp, err := sms.req(smsRuSmsSendPath, sms.params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data := new(ResponseSMSSent)
	if err = json.NewDecoder(resp.Body).Decode(data); err != nil {
		return nil, err
	}

	return data, nil
}

func (sms *SMSStatus) Get(smsID ...string) (*SMSStatuses, error) {
	params := url.Values{}
	params.Set("sms_id", strings.Join(smsID, ","))

	resp, err := sms.req(smsRuSmsStatusPath, params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data := new(SMSStatuses)
	if err = json.NewDecoder(resp.Body).Decode(data); err != nil {
		return nil, err
	}

	return data, nil
}
