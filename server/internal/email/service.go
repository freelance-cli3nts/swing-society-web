package email

import (
	"fmt"
	"log"
	"net/smtp"
	"swing-society-website/server/internal/config"
)

// Service handles sending emails via SMTP.
type Service struct {
	cfg *config.Config
}

// NewService creates an email service using the loaded application config.
func NewService(cfg *config.Config) *Service {
	return &Service{cfg: cfg}
}

// send is the core send function. It builds a plain MIME message and sends via SMTP.
func (s *Service) send(to, subject, htmlBody string) error {
	if !s.cfg.Email.Enabled {
		log.Printf("Email service not configured — skipping send to %s", to)
		return nil
	}

	auth := smtp.PlainAuth("", s.cfg.Email.SMTPUser, s.cfg.Email.SMTPPass, s.cfg.Email.SMTPHost)
	addr := s.cfg.Email.SMTPHost + ":" + s.cfg.Email.SMTPPort

	msg := []byte(fmt.Sprintf(
		"From: Swing Society <%s>\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s",
		s.cfg.Email.SMTPUser, to, subject, htmlBody,
	))

	if err := smtp.SendMail(addr, auth, s.cfg.Email.SMTPUser, []string{to}, msg); err != nil {
		return fmt.Errorf("smtp send to %s: %w", to, err)
	}
	return nil
}

// SendWelcome sends a registration confirmation email to the new registrant.
func (s *Service) SendWelcome(name, email string) error {
	body, err := welcomeEmail(name)
	if err != nil {
		return err
	}
	return s.send(email, "Добре дошли в Swing Society! / Welcome to Swing Society!", body)
}

// SendRegistrationNotification notifies the admin of a new registration.
func (s *Service) SendRegistrationNotification(name, email, phone string) error {
	if s.cfg.Email.RecipientEmail == "" {
		return nil
	}
	body, err := registrationNotificationEmail(name, email, phone)
	if err != nil {
		return err
	}
	return s.send(s.cfg.Email.RecipientEmail, "Нова регистрация: "+name, body)
}

// SendContactNotification notifies the admin of a new contact form submission.
func (s *Service) SendContactNotification(name, email, message string) error {
	if s.cfg.Email.RecipientEmail == "" {
		return nil
	}
	body, err := contactNotificationEmail(name, email, message)
	if err != nil {
		return err
	}
	return s.send(s.cfg.Email.RecipientEmail, "Ново запитване от: "+name, body)
}

// SendNewsletterConfirmation sends a subscription confirmation to the subscriber.
func (s *Service) SendNewsletterConfirmation(name, email string) error {
	body, err := newsletterConfirmationEmail(name)
	if err != nil {
		return err
	}
	return s.send(email, "Абонирахте се за Swing Society бюлетин!", body)
}
