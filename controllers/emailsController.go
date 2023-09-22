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

func (ec EmailsController) SendEmail(ctx *gin.Context) {
	ec.emailsService.SendEmail(ctx)
}
