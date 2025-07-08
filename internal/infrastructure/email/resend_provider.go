package email

import (
	"context"
	"fmt"
	"log"
	"trading-alchemist/internal/config"
	"trading-alchemist/internal/domain/auth"
	"trading-alchemist/internal/domain/services"

	"github.com/resend/resend-go/v2"
)
type ResendProvider struct {
	client *resend.Client
	config *config.Config
}
// NewResendProvider creates a new Resend email provider
func NewResendProvider(cfg *config.Config) services.EmailService {
	client := resend.NewClient(cfg.Email.ResendAPIKey)
	return &ResendProvider{
		client: client,
		config: cfg,
	}
}
// SendMagicLinkEmail sends a magic link email to the user using Resend
func (r *ResendProvider) SendMagicLinkEmail(ctx context.Context, user *auth.User, magicLink *auth.MagicLink) error {
	subject := "Your Magic Link to " + r.config.App.Name
	html := r.buildMagicLinkEmailBody(user, magicLink)
	return r.sendEmail(ctx, user.Email, subject, html)
}
// SendWelcomeEmail sends a welcome email to new users using Resend
func (r *ResendProvider) SendWelcomeEmail(ctx context.Context, user *auth.User) error {
	subject := fmt.Sprintf("Welcome to %s!", r.config.App.Name)
	html := r.buildWelcomeEmailBody(user)
	return r.sendEmail(ctx, user.Email, subject, html)
}
// SendEmailVerificationEmail sends an email verification email using Resend
func (r *ResendProvider) SendEmailVerificationEmail(ctx context.Context, user *auth.User, magicLink *auth.MagicLink) error {
	subject := "Verify Your Email Address for " + r.config.App.Name
	html := r.buildEmailVerificationBody(user, magicLink)
	return r.sendEmail(ctx, user.Email, subject, html)
}
// sendEmail sends an email using the Resend API
func (r *ResendProvider) sendEmail(ctx context.Context, to, subject, html string) error {
	params := &resend.SendEmailRequest{
		From:    fmt.Sprintf("%s <%s>", r.config.Email.FromName, r.config.Email.FromEmail),
		To:      []string{to},
		Subject: subject,
		Html:    html,
	}
	sent, err := r.client.Emails.Send(params)
	if err != nil {
		// The Resend SDK does not export typed errors for all API responses,
		// so we log the full error string provided, which contains useful details.
		log.Printf("Failed to send email via Resend: %v", err)
		return fmt.Errorf("failed to send email: %w", err)
	}
	// Log the email ID for tracking (in production, use proper logging)
	log.Printf("Email sent successfully via Resend, ID: %s\n", sent.Id)
	return nil
}
// buildMagicLinkEmailBody builds the magic link email body
func (r *ResendProvider) buildMagicLinkEmailBody(user *auth.User, magicLink *auth.MagicLink) string {
	magicLinkURL := fmt.Sprintf("%s/auth/verify?token=%s", r.config.App.FrontendBaseURL, magicLink.Token)
	
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Your Magic Link</title>
    <style>
        body { 
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; 
            line-height: 1.6; 
            color: #333333; /* Dark gray for text */
            margin: 0; 
            padding: 0; 
            background-color: #f6f6f6; /* Light gray background */
        }
        .container { 
            max-width: 600px; 
            margin: 20px auto; 
            padding: 30px; 
            background-color: #ffffff; /* White container background */
            border-radius: 8px; /* Slightly rounded corners */
            box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05); /* Subtle shadow */
        }
        .header { 
            text-align: center; 
            padding-bottom: 25px; 
            margin-bottom: 25px; 
            border-bottom: 1px solid #eeeeee; /* Lighter border */
        }
        .header h1 {
            color: #6A0DAD; /* Deep purple for brand name */
            font-size: 28px;
            margin: 0;
            padding: 0;
        }
        h2 {
            color: #333333; /* Dark gray for headings */
            font-size: 24px;
            margin-top: 0;
            margin-bottom: 15px;
        }
        p {
            margin-bottom: 15px;
            font-size: 16px;
            color: #333333;
        }
        .button-container { 
            text-align: center; 
            margin: 30px 0;
        }
        .button { 
            display: inline-block; 
            padding: 15px 30px; 
            background-color: #6A0DAD; /* Deep purple button */
            color: white; 
            text-decoration: none; 
            border-radius: 6px; 
            font-weight: bold;
            font-size: 18px;
            transition: background-color 0.3s ease;
        }
        .button:hover {
            background-color: #5C0AA0; /* Slightly darker purple on hover */
        }
        .link-text {
            word-break: break-all; 
            background-color: #f8f9fa; /* Very light gray for link background */
            padding: 12px; 
            border-radius: 4px;
            font-family: monospace;
            font-size: 14px;
            color: #333333;
            border: 1px dashed #cccccc; /* Dashed border for the link box */
            margin-top: 20px;
        }
        .footer { 
            margin-top: 35px; 
            padding-top: 25px; 
            border-top: 1px solid #eeeeee; 
            font-size: 13px; 
            color: #666666; /* Medium gray for footer text */
            text-align: center;
        }
        .footer p {
            margin: 5px 0;
        }
        @media only screen and (max-width: 600px) {
            .container {
                margin: 10px;
                padding: 20px;
            }
            .header h1 {
                font-size: 24px;
            }
            h2 {
                font-size: 20px;
            }
            .button {
                padding: 12px 25px;
                font-size: 16px;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>%s</h1>
        </div>
        <h2>Hi %s!</h2>
        <p>To sign in to your account, simply click the button below:</p>
        <div class="button-container">
            <a href="%s" class="button">Sign In to Your Account</a>
        </div>
        <p>Alternatively, you can copy and paste this secure link into your browser:</p>
        <p class="link-text"><code>%s</code></p>
        <div class="footer">
            <p>This link will expire in <strong>15 minutes</strong> for security reasons.</p>
            <p>If you didn't request this link, please disregard this email.</p>
            <p>Thank you,<br>The %s Team</p>
        </div>
    </div>
</body>
</html>
	`, r.config.App.Name, user.DisplayName(), magicLinkURL, magicLinkURL, r.config.App.Name)
}
// buildWelcomeEmailBody builds the welcome email body
func (r *ResendProvider) buildWelcomeEmailBody(user *auth.User) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Welcome to %s!</title>
    <style>
        body { 
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; 
            line-height: 1.6; 
            color: #333333; 
            margin: 0; 
            padding: 0; 
            background-color: #f6f6f6; 
        }
        .container { 
            max-width: 600px; 
            margin: 20px auto; 
            padding: 30px; 
            background-color: #ffffff; 
            border-radius: 8px; 
            box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05); 
        }
        .header { 
            text-align: center; 
            padding-bottom: 25px; 
            margin-bottom: 25px; 
            border-bottom: 1px solid #eeeeee; 
        }
        .header h1 {
            color: #6A0DAD; 
            font-size: 28px;
            margin: 0;
            padding: 0;
        }
        p {
            margin-bottom: 15px;
            font-size: 16px;
            color: #333333;
        }
        .footer { 
            margin-top: 35px; 
            padding-top: 25px; 
            border-top: 1px solid #eeeeee; 
            font-size: 13px; 
            color: #666666; 
            text-align: center;
        }
        .footer p {
            margin: 5px 0;
        }
        @media only screen and (max-width: 600px) {
            .container {
                margin: 10px;
                padding: 20px;
            }
            .header h1 {
                font-size: 24px;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Welcome to %s!</h1>
        </div>
        <p>Hi %s,</p>
        <p>We are thrilled to have you join the %s community! Your journey with us begins now.</p>
        <p>To help you get started, here are a few things you can do:</p>
        <ul>
            <li><strong>Explore your dashboard:</strong> Dive into your personalized space and discover all our features.</li>
            <li><strong>Set up your profile:</strong> Make sure your account is complete and ready to go.</li>
            <li><strong>Check out our guide:</strong> [Link to Onboarding Guide/FAQ - *Highly Recommended Addition*]</li>
        </ul>
        <p>We're here to support you every step of the way. If you have any questions, don't hesitate to reach out to our support team.</p>
        <p>Best regards,<br>The %s Team</p>
        <div class="footer">
            <p>You're receiving this email because you signed up for %s.</p>
            <p>If you have any questions, please contact us at [Support Email/Link].</p>
        </div>
    </div>
</body>
</html>
	`, r.config.App.Name, r.config.App.Name, user.DisplayName(), r.config.App.Name, r.config.App.Name, r.config.App.Name)
}
// buildEmailVerificationBody builds the email verification body
func (r *ResendProvider) buildEmailVerificationBody(user *auth.User, magicLink *auth.MagicLink) string {
	verificationURL := fmt.Sprintf("%s/auth/verify?token=%s", r.config.App.FrontendBaseURL, magicLink.Token)
	
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Verify Your Email Address</title>
    <style>
        body { 
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; 
            line-height: 1.6; 
            color: #333333; 
            margin: 0; 
            padding: 0; 
            background-color: #f6f6f6; 
        }
        .container { 
            max-width: 600px; 
            margin: 20px auto; 
            padding: 30px; 
            background-color: #ffffff; 
            border-radius: 8px; 
            box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05); 
        }
        .header { 
            text-align: center; 
            padding-bottom: 25px; 
            margin-bottom: 25px; 
            border-bottom: 1px solid #eeeeee; 
        }
        .header h1 {
            color: #6A0DAD; 
            font-size: 28px;
            margin: 0;
            padding: 0;
        }
        p {
            margin-bottom: 15px;
            font-size: 16px;
            color: #333333;
        }
        .button-container { 
            text-align: center; 
            margin: 30px 0;
        }
        .button { 
            display: inline-block; 
            padding: 15px 30px; 
            background-color: #6A0DAD; /* Deep purple button */
            color: white; 
            text-decoration: none; 
            border-radius: 6px; 
            font-weight: bold;
            font-size: 18px;
            transition: background-color 0.3s ease;
        }
        .button:hover {
            background-color: #5C0AA0; /* Slightly darker purple on hover */
        }
        .link-text {
            word-break: break-all; 
            background-color: #f8f9fa; 
            padding: 12px; 
            border-radius: 4px;
            font-family: monospace;
            font-size: 14px;
            color: #333333;
            border: 1px dashed #cccccc;
            margin-top: 20px;
        }
        .footer { 
            margin-top: 35px; 
            padding-top: 25px; 
            border-top: 1px solid #eeeeee; 
            font-size: 13px; 
            color: #666666; 
            text-align: center;
        }
        .footer p {
            margin: 5px 0;
        }
        @media only screen and (max-width: 600px) {
            .container {
                margin: 10px;
                padding: 20px;
            }
            .header h1 {
                font-size: 24px;
            }
            .button {
                padding: 12px 25px;
                font-size: 16px;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Verify Your Email Address</h1>
        </div>
        <p>Hi %s,</p>
        <p>To complete your registration and verify your email address for %s, please click the button below:</p>
        <div class="button-container">
            <a href="%s" class="button">Verify My Email</a>
        </div>
        <p>If the button doesn't work, you can copy and paste this link into your browser:</p>
        <p class="link-text"><code>%s</code></p>
        <div class="footer">
            <p>This verification link will expire in <strong>15 minutes</strong>.</p>
            <p>If you didn't request this verification, you can safely ignore this email.</p>
            <p>Thank you,<br>The %s Team</p>
        </div>
    </div>
</body>
</html>
	`, user.DisplayName(), r.config.App.Name, verificationURL, verificationURL, r.config.App.Name)
}