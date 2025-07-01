package email

import (
	"context"
	"fmt"

	"trading-alchemist/internal/config"
	"trading-alchemist/internal/domain/entities"
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
func (r *ResendProvider) SendMagicLinkEmail(ctx context.Context, user *entities.User, magicLink *entities.MagicLink) error {
	subject := "Your Magic Link"
	html := r.buildMagicLinkEmailBody(user, magicLink)

	return r.sendEmail(ctx, user.Email, subject, html)
}

// SendWelcomeEmail sends a welcome email to new users using Resend
func (r *ResendProvider) SendWelcomeEmail(ctx context.Context, user *entities.User) error {
	subject := fmt.Sprintf("Welcome to %s!", r.config.App.Name)
	html := r.buildWelcomeEmailBody(user)

	return r.sendEmail(ctx, user.Email, subject, html)
}

// SendEmailVerificationEmail sends an email verification email using Resend
func (r *ResendProvider) SendEmailVerificationEmail(ctx context.Context, user *entities.User, magicLink *entities.MagicLink) error {
	subject := "Verify Your Email Address"
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
		return fmt.Errorf("failed to send email via Resend: %w", err)
	}

	// Log the email ID for tracking (in production, use proper logging)
	fmt.Printf("Email sent successfully via Resend, ID: %s\n", sent.Id)

	return nil
}

// buildMagicLinkEmailBody builds the magic link email body
func (r *ResendProvider) buildMagicLinkEmailBody(user *entities.User, magicLink *entities.MagicLink) string {
	magicLinkURL := fmt.Sprintf("%s/auth/verify?token=%s", r.config.App.BaseURL, magicLink.Token)
	
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Magic Link</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; line-height: 1.6; color: #333; margin: 0; padding: 0; background-color: #f6f6f6; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; background-color: #ffffff; }
        .header { text-align: center; border-bottom: 1px solid #eee; padding-bottom: 20px; margin-bottom: 20px; }
        .button { display: inline-block; padding: 12px 24px; background-color: #007bff; color: white; text-decoration: none; border-radius: 6px; margin: 20px 0; }
        .footer { margin-top: 30px; padding-top: 20px; border-top: 1px solid #eee; font-size: 12px; color: #666; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>%s</h1>
        </div>
        <h2>Hi %s!</h2>
        <p>Click the button below to sign in to your account:</p>
        <div style="text-align: center;">
            <a href="%s" class="button">Sign In</a>
        </div>
        <p>Or copy and paste this link into your browser:</p>
        <p style="word-break: break-all; background-color: #f8f9fa; padding: 10px; border-radius: 4px;"><code>%s</code></p>
        <div class="footer">
            <p><small>This link will expire in 15 minutes.</small></p>
            <p><small>If you didn't request this link, you can safely ignore this email.</small></p>
        </div>
    </div>
</body>
</html>
	`, r.config.App.Name, user.DisplayName(), magicLinkURL, magicLinkURL)
}

// buildWelcomeEmailBody builds the welcome email body
func (r *ResendProvider) buildWelcomeEmailBody(user *entities.User) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Welcome</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; line-height: 1.6; color: #333; margin: 0; padding: 0; background-color: #f6f6f6; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; background-color: #ffffff; }
        .header { text-align: center; border-bottom: 1px solid #eee; padding-bottom: 20px; margin-bottom: 20px; }
        .footer { margin-top: 30px; padding-top: 20px; border-top: 1px solid #eee; font-size: 12px; color: #666; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Welcome to %s!</h1>
        </div>
        <p>Hi %s,</p>
        <p>Thank you for joining us! We're excited to have you on board.</p>
        <p>Get started by exploring our features and let us know if you need any help.</p>
        <p>Best regards,<br>The %s Team</p>
        <div class="footer">
            <p><small>You're receiving this email because you signed up for %s.</small></p>
        </div>
    </div>
</body>
</html>
	`, r.config.App.Name, user.DisplayName(), r.config.App.Name, r.config.App.Name)
}

// buildEmailVerificationBody builds the email verification body
func (r *ResendProvider) buildEmailVerificationBody(user *entities.User, magicLink *entities.MagicLink) string {
	verificationURL := fmt.Sprintf("%s/auth/verify-email?token=%s", r.config.App.BaseURL, magicLink.Token)
	
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Verify Your Email</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; line-height: 1.6; color: #333; margin: 0; padding: 0; background-color: #f6f6f6; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; background-color: #ffffff; }
        .header { text-align: center; border-bottom: 1px solid #eee; padding-bottom: 20px; margin-bottom: 20px; }
        .button { display: inline-block; padding: 12px 24px; background-color: #28a745; color: white; text-decoration: none; border-radius: 6px; margin: 20px 0; }
        .footer { margin-top: 30px; padding-top: 20px; border-top: 1px solid #eee; font-size: 12px; color: #666; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Verify Your Email Address</h1>
        </div>
        <p>Hi %s,</p>
        <p>Please click the button below to verify your email address:</p>
        <div style="text-align: center;">
            <a href="%s" class="button">Verify Email</a>
        </div>
        <p>Or copy and paste this link into your browser:</p>
        <p style="word-break: break-all; background-color: #f8f9fa; padding: 10px; border-radius: 4px;"><code>%s</code></p>
        <div class="footer">
            <p><small>This verification link will expire in 15 minutes.</small></p>
            <p><small>If you didn't request this verification, you can safely ignore this email.</small></p>
        </div>
    </div>
</body>
</html>
	`, user.DisplayName(), verificationURL, verificationURL)
} 