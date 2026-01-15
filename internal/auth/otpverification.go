package auth

const maxRetries = 3

type OTPVerification interface {
	SendVerificationCode(email string) (string, error)
	VerifyCode(email, code string) error
}

type MockOTPVerification struct{}

func (m MockOTPVerification) SendVerificationCode(userEmail string) (string, error) {
	return "", nil
}

func (m MockOTPVerification) VerifyCode(email, code string) error {
	return nil
}
