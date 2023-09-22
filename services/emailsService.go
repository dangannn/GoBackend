package services

import (
	"GoBackend/repositories"
	"crypto/tls"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
	"net/http"
)

type EmailService struct {
	emailsRepository *repositories.EmailRepository
}

func NewEmailService(emailsRepository *repositories.EmailRepository) *EmailService {
	return &EmailService{
		emailsRepository: emailsRepository,
	}
}

func (es EmailService) SendEmail(ctx *gin.Context) {
	fmt.Println("зашло")
	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", "nikonorovdan14@gmail.com")

	// Set E-Mail receivers
	m.SetHeader("To", "nikonorovdan14@gmail.com")

	// Set E-Mail subject
	m.SetHeader("Subject", "Gomail test subject")

	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/plain", "This is Gomail test body")

	// Settings for SMTP server
	d := gomail.NewDialer("smtp.gmail.com", 587, "nikonorovdan14@gmail.com", "oeyl hhlc vsev idet")

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		panic(err)
	}
	ctx.JSON(http.StatusOK, "email sent")
	return
}
