package email

// Sender is the interface that wraps the basic email sending methods.
// This allows you to easily switch to SendGrid, AWS SES, or use a mock sender in tests.
type Sender interface {
	SendOTP(recipientEmail, otp string) error
}
