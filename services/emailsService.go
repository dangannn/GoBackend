package services

import (
	"GoBackend/repositories"
	"crypto/tls"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
	"gopkg.in/gomail.v2"
	"log"
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
func (es EmailService) TaskScheduling() {
	c := cron.New()

	//0 - нулевая секунда (то есть начало каждой минуты).
	//0 - нулевая минута (то есть начало каждого часа).
	//9 - девятый час (то есть 9:00 утра).
	//* - любой день месяца (то есть не учитывается день месяца).
	//* - любой месяц (то есть не учитывается месяц).
	//* - любой день недели (то есть не учитывается день недели).
	from := "nikonorovdan14@gmail.com"
	// Запланировать задачу для отправки статистики каждый день в 9:00 утра
	err := c.AddFunc("0 34 23 * * *", func() {
		es.sendStatisticsEmail(from, from)
	})
	if err != nil {
		log.Fatalf("Failed to schedule task: %v", err)
	}

	c.Start()

}

func (es EmailService) SendEmail(c *gin.Context) {
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
	c.JSON(http.StatusOK, "email sent")
	return
}

func (es EmailService) AddView() {
	err := es.emailsRepository.AddView()
	if err != nil {
		log.Println("Ошибка добавления просмотра в статистику")

	}
}
func (es EmailService) AddNewComment() {
	err := es.emailsRepository.AddNewComment()
	if err != nil {
		log.Println("Ошибка добавления нового комментария в статистику")

	}
}

func (es EmailService) sendStatisticsEmail(from, to string) {

	stats, err := es.emailsRepository.GetDailyStats()
	if err != nil {
		log.Println("Ошибка получения статистики")
	}

	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", from)

	// Set E-Mail receivers
	m.SetHeader("To", to)

	// Set E-Mail subject
	m.SetHeader("Subject", "Daily stats")

	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/plain", fmt.Sprintf("Новых комментариев за день: %d \nПросмотров на комментариях: %d", stats.NewComments, stats.Views))

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
	log.Println("email sent")
	err = es.emailsRepository.ResetDailyStats()
	if err != nil {
		log.Println("Ошибка сброса статистики")
	}
	return
}
