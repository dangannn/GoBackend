package controllers

import (
	"GoBackend/services"
	"github.com/gin-gonic/gin"
)

type EmailsController struct {
	emailsService *services.EmailService
}

func NewEmailController(emailsService *services.EmailService) *EmailsController {
	return &EmailsController{
		emailsService: emailsService,
	}
}

func (ec EmailsController) SendEmail(c *gin.Context) {
	ec.emailsService.SendEmail(c)
}

func (ec EmailsController) AddView(c *gin.Context) {
	ec.emailsService.AddView()
}

func (ec EmailsController) AddNewComment(c *gin.Context) {
	ec.emailsService.AddNewComment()
}
