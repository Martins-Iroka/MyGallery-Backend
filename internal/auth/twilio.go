package auth

import (
	"fmt"
	"log"
	"time"

	"github.com/twilio/twilio-go"
	verify "github.com/twilio/twilio-go/rest/verify/v2"
)

type TwilioVerification struct {
	client    *twilio.RestClient
	serviceID string
}

func NewTwilioVerification(accountSID, authToken, serviceID string) *TwilioVerification {
	clientParam := twilio.ClientParams{
		Username: accountSID,
		Password: authToken,
	}
	client := twilio.NewRestClientWithParams(clientParam)
	return &TwilioVerification{client: client, serviceID: serviceID}
}

func (twilio *TwilioVerification) SendVerificationCode(userEmail string) error {
	channelEmail := "email"
	params := &verify.CreateVerificationParams{
		To:      &userEmail,
		Channel: &channelEmail,
	}
	for i := range maxRetries {
		resp, err := twilio.client.VerifyV2.CreateVerification(twilio.serviceID, params)
		if err != nil {
			log.Printf("failed to send email attempt %d of %d", i+1, maxRetries)
			log.Printf("Error: %v", err.Error())

			time.Sleep(time.Second * time.Duration(i+1))
			continue
		} else if resp.Sid != nil {
			return nil
		}
	}

	return fmt.Errorf("failed to send email after %d attempts", maxRetries)
}

func (twilio *TwilioVerification) VerifyCode(email, code string) error {
	params := &verify.CreateVerificationCheckParams{
		To:   &email,
		Code: &code,
	}
	for i := range maxRetries {
		resp, err := twilio.client.VerifyV2.CreateVerificationCheck(twilio.serviceID, params)
		if err != nil {
			log.Printf("failed to verify code attempt %d of %d", i+1, maxRetries)
			log.Printf("Error: %v", err.Error())

			time.Sleep(time.Second * time.Duration(i+1))
			continue
		} else if resp.Sid != nil {
			return nil
		}
	}

	return fmt.Errorf("failed to verify code after %d attempts", maxRetries)
}
