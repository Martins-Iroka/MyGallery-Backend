package auth

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/stytchauth/stytch-go/v16/stytch/consumer/otp"
	"github.com/stytchauth/stytch-go/v16/stytch/consumer/otp/email"
	"github.com/stytchauth/stytch-go/v16/stytch/consumer/stytchapi"
)

type StytchVerification struct {
	client *stytchapi.API
}

func NewStytchVerification(projectID, secret string) (*StytchVerification, error) {
	client, err := stytchapi.NewClient(projectID, secret)
	if err != nil {
		return nil, err
	}
	return &StytchVerification{client: client}, nil
}

func (sv *StytchVerification) SendVerificationCode(userEmail string) (string, error) {
	params := &email.LoginOrCreateParams{
		Email:             userEmail,
		ExpirationMinutes: 10,
	}

	for i := range maxRetries {
		resp, err := sv.client.OTPs.Email.LoginOrCreate(context.Background(), params)
		if err != nil {
			log.Printf("failed to send email attempt %d of %d", i+1, maxRetries)
			log.Printf("Error: %v", err.Error())

			time.Sleep(time.Second * time.Duration(i+1))
			continue
		} else if resp.StatusCode != 200 {
			return "", fmt.Errorf("error sending OTP. status code is %d", resp.StatusCode)
		} else {
			return resp.EmailID, nil
		}
	}

	return "", fmt.Errorf("failed to send email after %d attempts", maxRetries)
}

func (sv *StytchVerification) VerifyCode(email, code string) error {
	params := &otp.AuthenticateParams{
		MethodID:               email,
		Code:                   code,
		SessionDurationMinutes: 60,
	}
	for i := range maxRetries {
		resp, err := sv.client.OTPs.Authenticate(context.Background(), params)
		if err != nil {
			log.Printf("failed to verify code attempt %d of %d", i+1, maxRetries)
			log.Printf("Error: %v", err.Error())

			time.Sleep(time.Second * time.Duration(i+1))
			continue
		} else if resp.StatusCode != 200 {
			return fmt.Errorf("error verifying OTP. status code is %d", resp.StatusCode)
		} else {
			return nil
		}
	}

	return fmt.Errorf("failed to verify code after %d attempts", maxRetries)
}
