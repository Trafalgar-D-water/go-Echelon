package email

import (
	"context"
	"fmt"
	"log"

	brevo "github.com/getbrevo/brevo-go/lib"
	"github.com/go-Echelon/go-Echelon/pkg/core/config"
)

// SendOTP sends an OTP email directly using the Brevo API client
func SendOTP(recipientEmail, otp string) error {
	ctx := context.Background()
	cfg := brevo.NewConfiguration()
	cfg.AddDefaultHeader("api-key", config.GetConfig().BrevoAPIKey)

	client := brevo.NewAPIClient(cfg)

	// Fetch the verified account email directly from Brevo
	acct, _, err := client.AccountApi.GetAccount(ctx)
	if err != nil {
		log.Println("Error when fetching Brevo Account details: ", err.Error())
		return err
	}

	senderName := "Go-Echelon"
	if acct.CompanyName != "" {
		senderName = acct.CompanyName
	}

	// The email template
	htmlContent := fmt.Sprintf("<html><body><h1>Welcome to %s!</h1><p>Your OTP for registration is: <strong>%s</strong></p></body></html>", senderName, otp)

	sendSmtpEmail := brevo.SendSmtpEmail{
		Sender: &brevo.SendSmtpEmailSender{
			Name:  senderName,
			Email: acct.Email, // Dynamically use the verified email from the Brevo account
		},
		To: []brevo.SendSmtpEmailTo{
			{
				Email: recipientEmail,
			},
		},
		Subject:     fmt.Sprintf("Verify your %s Account", senderName),
		HtmlContent: htmlContent,
	}

	result, _, err := client.TransactionalEmailsApi.SendTransacEmail(ctx, sendSmtpEmail)
	if err != nil {
		log.Println("Error when calling Brevo API: ", err.Error())
		return err
	}

	log.Println("OTP Email sent successfully, MessageId: ", result.MessageId)
	return nil
}
