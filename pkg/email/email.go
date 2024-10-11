package email

import (
	"bytes"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"path/filepath"
	"regexp"
	"stori_challenge/pkg/models"
	"strings"
	"text/template"

	"github.com/joho/godotenv"
)

// SendEmail sends an email using the SMTP protocol with the given EmailData.
func SendEmail(data models.EmailData) error {
	// Path to the .env file in the root of the project
	envPath := filepath.Join(".", ".env")

	// Load the .env file
	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// SMTP server configuration
	smtpServer := os.Getenv("SMTP_SERVER")     // SMTP server address
	smtpPort := os.Getenv("SMTP_PORT")         // SMTP server port
	senderEmail := os.Getenv("SMTP_SENDER")    // Sender's email address
	senderPassword := os.Getenv("SMTP_PASSWD") // Sender's email password

	// Validate environment variables for SMTP configuration
	if smtpServer == "" || smtpPort == "" || senderEmail == "" || senderPassword == "" {
		return fmt.Errorf("missing SMTP configuration in environment variables")
	}

	// Authenticate with the SMTP server
	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpServer)

	// External HTML template file for the email body
	templateFile := filepath.Join("web", "template", "email_template.html")

	// Parse the HTML template file
	t, err := template.ParseFiles(templateFile)
	if err != nil {
		return fmt.Errorf("failed to parse email template: %w", err)
	}

	// Create a buffer to hold the output of the template
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		return fmt.Errorf("failed to execute email template: %w", err)
	}

	// Email message in HTML format
	htmlMessage := tpl.String()

	// List of email recipients
	to := []string{os.Getenv("SMTP_SENDER")} // Add sender's email to recipients
	to = append(to, data.EmailTo)            // Add recipient's email

	// CC (carbon copy) recipients
	cc := os.Getenv("SMTP_CC")
	ccEmails := []string{}
	if cc != "" {
		ccEmails = append(ccEmails, cc) // Add CC emails if provided
	}

	// Subject of the email
	subject := os.Getenv("SMTP_SUBJECT")

	// Create the email body with headers
	body := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\nCc: %s\r\nMIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n%s",
		strings.Join(to, ","),       // Combine all To recipients
		subject,                     // Subject line
		strings.Join(ccEmails, ","), // Combine all CC recipients
		htmlMessage))                // HTML message body

	// Combine To and CC recipients for sending
	recipients := append(to, ccEmails...)

	// Send the email
	if err = smtp.SendMail(smtpServer+":"+smtpPort, auth, senderEmail, recipients, body); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	// Log the successful email sending
	log.Printf("Email successfully sent to: %s", strings.Join(recipients, ", "))

	return nil
}

// IsValidEmail validates the format of an email address using a regular expression.
func IsValidEmail(email string) bool {
	// Regular expression to validate the email format
	re := regexp.MustCompile(`^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`)
	return re.MatchString(email) // Return true if the email matches the pattern
}
