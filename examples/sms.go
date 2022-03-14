package main

import (
	"fmt"

	"github.com/myback/smsru-go"
)

func main() {
	smsRuApi := smsru.New("FD48334A-BC98-9E70-C196-23456789ABCD")

	resp, err := smsRuApi.NewSMS().
		MultipleRecipients("hello world", "12345678900", "19876543200").
		ConsiderTimezone().
		Send()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", resp)

	resp, err = smsRuApi.NewSMS().
		Add("code: 1234", "12345678900").
		TTL(5).
		Send()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", resp)

	resp, err = smsRuApi.NewSMS().
		Add("Hello Alice!", "12345678900").
		Add("Hello Bob!", "19876543200").
		Send()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", resp)

	status, err := smsRuApi.NewSMSStatus().Get("000000-000001", "000000-000002", "000000-000003")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", status)
}
