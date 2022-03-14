package smsru

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	StatusOK    = "OK"
	StatusError = "ERROR"

	smsRuUrl = "https://sms.ru"
)

type request func(string, url.Values) (*http.Response, error)

type SmsRu struct {
	apiID string
}

var smsRuCodeStatus = map[int]string{
	-1:  "not found",
	100: "success",
	101: "the message is sent to the operator",
	102: "message has been sent (in transit)",
	103: "message delivered",
	104: "cannot be delivered: life time expired",
	105: "cannot be delivered: deleted by operator",
	106: "cannot be delivered: phone failure",
	107: "cannot be delivered: unknown reason",
	108: "cannot be delivered: rejected",
	110: "message read",
	150: "cannot be delivered: no route found to this number",
	200: "wrong api_id",
	201: "payment required",
	202: "the recipient's phone number is incorrect, or there is no route to it",
	203: "empty message",
	204: "the sender's name was not agreed with the administration",
	205: "message is too long (more than 8 sms)",
	206: "daily message limit exceeded",
	207: "there is no route for message delivery to this number",
	208: "the time parameter is incorrect",
	209: "this number is in the stop list",
	210: "use POST method instead GET",
	211: "method not found",
	212: "the message text must be encoded in UTF-8",
	213: "more than 100 numbers in the list of recipients",
	220: "the service is temporarily unavailable, please try again later",
	230: "daily message limit on this number was exceeded",
	231: "exceeded the limit of identical messages to this number per minute",
	232: "exceeded the limit of identical messages to this number per day",
	300: "invalid token (maybe it was expired or your IP was changed)",
	301: "invalid password or user does not exist",
	302: "the user is authorized, but the account is not activated",
	303: "confirmation code is incorrect",
	304: "too many verification codes have been sent, please try again later",
	305: "too many wrong attempts to enter the verification code, please try again later",
	400: "the number has not yet been confirmed",
	401: "the number is confirmed",
	402: "the allotted time for checking has expired, or the check_id specified is incorrect",
	500: "internal gate error",
	901: "invalid Url (should begin with 'http://')",
	902: "callback is not defined",
}

func MessageByCode(id int) string {
	return smsRuCodeStatus[id]
}

func New(apiID string) *SmsRu {
	return &SmsRu{apiID}
}

func (sms *SmsRu) request(path string, val url.Values) (*http.Response, error) {
	u, err := url.Parse(smsRuUrl)
	if err != nil {
		return nil, err
	}

	u.Path = path

	val.Set("api_id", sms.apiID)
	val.Set("json", "1")
	u.RawQuery = val.Encode()

	req, err := http.NewRequest("", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "smsru-go")

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read body: %s", err)
		}

		return nil, fmt.Errorf("expect status code 200 got %d: %s", resp.StatusCode, b)
	}

	return resp, nil
}
