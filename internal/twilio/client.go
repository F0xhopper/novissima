package twilio

import (
	"github.com/twilio/twilio-go"
)

type Client struct {
	twilioClient *twilio.RestClient
	phoneNumber  string
}

func NewClient(accountSid, authToken, phoneNumber string) *Client {
	return &Client{
		twilioClient: twilio.NewRestClientWithParams(twilio.ClientParams{
			Username: accountSid,
			Password: authToken,
		}),
		phoneNumber: phoneNumber,
	}
}
