package email

import (
	"bytes"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"path/filepath"
	"stori_challenge/pkg/models"
	"strings"
	"text/template"
)

func SendEmail(data models.EmailData) error {

	// Configuración del servidor SMTP
	smtpServer := os.Getenv("SMTP_SERVER")
	smtpPort := os.Getenv("SMTP_PORT")
	senderEmail := os.Getenv("SMTP_SENDER")
	senderPassword := os.Getenv("SMTP_PASSWD")

	// Validación de variables de entorno
	if smtpServer == "" || smtpPort == "" || senderEmail == "" || senderPassword == "" {
		return fmt.Errorf("missing SMTP configuration in environment variables")
	}

	// Autenticación con el servidor SMTP
	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpServer)

	// Plantilla HTML externa
	//templateFile := "email_template.html"
	templateFile := filepath.Join("web", "template", "email_template.html")

	// Parseamos la plantilla HTML
	t, err := template.ParseFiles(templateFile)
	if err != nil {
		return fmt.Errorf("failed to parse email template: %w", err)
	}

	// Creamos un buffer para almacenar la salida de la plantilla
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		return fmt.Errorf("failed to execute email template: %w", err)
	}

	// Mensaje del correo electrónico
	htmlMessage := tpl.String()

	// Destinatarios
	to := []string{os.Getenv("SMTP_SENDER")}
	to = append(to, data.EmailTo)

	// Destinatarios en copia (CC)
	cc := os.Getenv("SMTP_CC")
	ccEmails := []string{}
	if cc != "" {
		ccEmails = append(ccEmails, cc)
	}

	subject := os.Getenv("SMTP_SUBJECT")

	// Cuerpo del correo
	body := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\nCc: %s\r\nMIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n%s",
		strings.Join(to, ","),
		subject,
		strings.Join(ccEmails, ","),
		htmlMessage))

	// Combinar destinatarios y copia
	recipients := append(to, ccEmails...)

	// Enviar el correo
	if err = smtp.SendMail(smtpServer+":"+smtpPort, auth, senderEmail, recipients, body); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	// Log del correo enviado
	log.Printf("Email enviado con éxito a: %s", strings.Join(recipients, ", "))

	return nil
}
