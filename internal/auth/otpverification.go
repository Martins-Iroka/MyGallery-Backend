package auth

type OTPVerification interface {
	SendVerificationCode(email string) error
	VerifyCode(email, code string) error
}

type MockOTPVerification struct{}

func (m MockOTPVerification) SendVerificationCode(email string) error {
	return nil
}

func (m MockOTPVerification) VerifyCode(email, code string) error {
	return nil
}
