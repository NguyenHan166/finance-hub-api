package utils

import (
	"bytes"
	"finance-hub-api/internal/config"
	"fmt"
	"html/template"
	"net/smtp"
)

// EmailService handles email sending
type EmailService struct {
	smtpHost     string
	smtpPort     string
	smtpUsername string
	smtpPassword string
	fromEmail    string
	fromName     string
	frontendURL  string
}

// NewEmailService creates a new email service
func NewEmailService(cfg *config.Config) *EmailService {
	return &EmailService{
		smtpHost:     cfg.Email.SMTPHost,
		smtpPort:     cfg.Email.SMTPPort,
		smtpUsername: cfg.Email.SMTPUsername,
		smtpPassword: cfg.Email.SMTPPassword,
		fromEmail:    cfg.Email.FromEmail,
		fromName:     cfg.Email.FromName,
		frontendURL:  cfg.Server.FrontendURL,
	}
}

// EmailData represents email template data
type EmailData struct {
	RecipientName string
	RecipientEmail string
	Subject       string
	Body          string
	Link          string
	Token         string
}

// SendVerificationEmail sends email verification email
func (s *EmailService) SendVerificationEmail(toEmail, toName, token string) error {
	verifyURL := fmt.Sprintf("%s/auth/verify-email?token=%s", 
		s.frontendURL, token)

	emailData := EmailData{
		RecipientName:  toName,
		RecipientEmail: toEmail,
		Subject:        "Xác thực email - Finance Hub",
		Link:           verifyURL,
		Token:          token,
	}

	htmlBody := s.getVerificationEmailTemplate(emailData)

	return s.sendEmail(toEmail, emailData.Subject, htmlBody)
}

// SendPasswordResetEmail sends password reset email
func (s *EmailService) SendPasswordResetEmail(toEmail, toName, token string) error {
	resetURL := fmt.Sprintf("%s/auth/reset-password?token=%s", 
		s.frontendURL, token)

	emailData := EmailData{
		RecipientName:  toName,
		RecipientEmail: toEmail,
		Subject:        "Đặt lại mật khẩu - Finance Hub",
		Link:           resetURL,
		Token:          token,
	}

	htmlBody := s.getPasswordResetEmailTemplate(emailData)

	return s.sendEmail(toEmail, emailData.Subject, htmlBody)
}

// sendEmail sends an email using SMTP
func (s *EmailService) sendEmail(to, subject, htmlBody string) error {
	// If SMTP not configured, just log and return (for development)
	if s.smtpHost == "" || s.smtpUsername == "" {
		fmt.Printf("📧 [EMAIL] Would send to: %s\nSubject: %s\nBody:\n%s\n", to, subject, htmlBody)
		return nil
	}

	// Prepare email headers
	from := fmt.Sprintf("%s <%s>", s.fromName, s.fromEmail)
	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	// Build email message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + htmlBody

	// Setup authentication
	auth := smtp.PlainAuth("", s.smtpUsername, s.smtpPassword, s.smtpHost)

	// Send email
	addr := fmt.Sprintf("%s:%s", s.smtpHost, s.smtpPort)
	err := smtp.SendMail(addr, auth, s.fromEmail, []string{to}, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// getVerificationEmailTemplate returns HTML template for email verification
func (s *EmailService) getVerificationEmailTemplate(data EmailData) string {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Subject}}</title>
</head>
<body style="margin: 0; padding: 0; font-family: Arial, sans-serif; background-color: #f4f4f7;">
    <table width="100%" cellpadding="0" cellspacing="0" border="0">
        <tr>
            <td align="center" style="padding: 40px 0;">
                <table width="600" cellpadding="0" cellspacing="0" border="0" style="background-color: #ffffff; border-radius: 8px; box-shadow: 0 2px 8px rgba(0,0,0,0.1);">
                    <!-- Header -->
                    <tr>
                        <td style="padding: 40px 40px 30px; text-align: center; background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); border-radius: 8px 8px 0 0;">
                            <h1 style="margin: 0; color: #ffffff; font-size: 28px; font-weight: bold;">Finance Hub</h1>
                        </td>
                    </tr>
                    <!-- Content -->
                    <tr>
                        <td style="padding: 40px;">
                            <h2 style="margin: 0 0 20px; color: #333333; font-size: 24px;">Xin chào {{.RecipientName}}!</h2>
                            <p style="margin: 0 0 20px; color: #666666; font-size: 16px; line-height: 1.5;">
                                Cảm ơn bạn đã đăng ký tài khoản Finance Hub. Để hoàn tất quá trình đăng ký, vui lòng xác thực địa chỉ email của bạn bằng cách nhấp vào nút bên dưới.
                            </p>
                            <div style="text-align: center; margin: 30px 0;">
                                <a href="{{.Link}}" style="display: inline-block; padding: 14px 40px; background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: #ffffff; text-decoration: none; border-radius: 6px; font-size: 16px; font-weight: bold;">
                                    Xác thực Email
                                </a>
                            </div>
                            <p style="margin: 20px 0 0; color: #999999; font-size: 14px; line-height: 1.5;">
                                Hoặc copy link sau vào trình duyệt:<br>
                                <a href="{{.Link}}" style="color: #667eea; word-break: break-all;">{{.Link}}</a>
                            </p>
                            <p style="margin: 30px 0 0; color: #999999; font-size: 14px; line-height: 1.5;">
                                Link này sẽ hết hạn sau 24 giờ. Nếu bạn không tạo tài khoản này, vui lòng bỏ qua email này.
                            </p>
                        </td>
                    </tr>
                    <!-- Footer -->
                    <tr>
                        <td style="padding: 30px; background-color: #f8f9fa; border-radius: 0 0 8px 8px; text-align: center;">
                            <p style="margin: 0; color: #999999; font-size: 12px;">
                                © 2026 Finance Hub. All rights reserved.
                            </p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
    </table>
</body>
</html>
`
	t := template.Must(template.New("email").Parse(tmpl))
	var buf bytes.Buffer
	t.Execute(&buf, data)
	return buf.String()
}

// getPasswordResetEmailTemplate returns HTML template for password reset
func (s *EmailService) getPasswordResetEmailTemplate(data EmailData) string {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Subject}}</title>
</head>
<body style="margin: 0; padding: 0; font-family: Arial, sans-serif; background-color: #f4f4f7;">
    <table width="100%" cellpadding="0" cellspacing="0" border="0">
        <tr>
            <td align="center" style="padding: 40px 0;">
                <table width="600" cellpadding="0" cellspacing="0" border="0" style="background-color: #ffffff; border-radius: 8px; box-shadow: 0 2px 8px rgba(0,0,0,0.1);">
                    <!-- Header -->
                    <tr>
                        <td style="padding: 40px 40px 30px; text-align: center; background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%); border-radius: 8px 8px 0 0;">
                            <h1 style="margin: 0; color: #ffffff; font-size: 28px; font-weight: bold;">Finance Hub</h1>
                        </td>
                    </tr>
                    <!-- Content -->
                    <tr>
                        <td style="padding: 40px;">
                            <h2 style="margin: 0 0 20px; color: #333333; font-size: 24px;">Đặt lại mật khẩu</h2>
                            <p style="margin: 0 0 20px; color: #666666; font-size: 16px; line-height: 1.5;">
                                Xin chào {{.RecipientName}},
                            </p>
                            <p style="margin: 0 0 20px; color: #666666; font-size: 16px; line-height: 1.5;">
                                Chúng tôi nhận được yêu cầu đặt lại mật khẩu cho tài khoản của bạn. Nhấp vào nút bên dưới để tạo mật khẩu mới.
                            </p>
                            <div style="text-align: center; margin: 30px 0;">
                                <a href="{{.Link}}" style="display: inline-block; padding: 14px 40px; background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%); color: #ffffff; text-decoration: none; border-radius: 6px; font-size: 16px; font-weight: bold;">
                                    Đặt lại mật khẩu
                                </a>
                            </div>
                            <p style="margin: 20px 0 0; color: #999999; font-size: 14px; line-height: 1.5;">
                                Hoặc copy link sau vào trình duyệt:<br>
                                <a href="{{.Link}}" style="color: #f5576c; word-break: break-all;">{{.Link}}</a>
                            </p>
                            <p style="margin: 30px 0 0; color: #999999; font-size: 14px; line-height: 1.5;">
                                Link này sẽ hết hạn sau 1 giờ. Nếu bạn không yêu cầu đặt lại mật khẩu, vui lòng bỏ qua email này hoặc liên hệ với chúng tôi nếu bạn có thắc mắc.
                            </p>
                        </td>
                    </tr>
                    <!-- Footer -->
                    <tr>
                        <td style="padding: 30px; background-color: #f8f9fa; border-radius: 0 0 8px 8px; text-align: center;">
                            <p style="margin: 0; color: #999999; font-size: 12px;">
                                © 2026 Finance Hub. All rights reserved.
                            </p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
    </table>
</body>
</html>
`
	t := template.Must(template.New("email").Parse(tmpl))
	var buf bytes.Buffer
	t.Execute(&buf, data)
	return buf.String()
}
