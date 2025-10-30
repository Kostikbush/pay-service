package ports

import (
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
)

type Service interface {
	InitPay(userId string) string
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (handler *Handler) InitPay(ginContext *gin.Context) {
	userId, ok := ginContext.GetQuery("userId")
	userId = strings.TrimSpace(userId)

	if !ok || userId == "" {
		ginContext.JSON(http.StatusBadRequest, gin.H{"message": "Отсутствуют необходимые параметры", "type": "error"})
		return
	}
	if utf8.RuneCountInString(userId) != 24 {
		ginContext.JSON(http.StatusBadRequest, gin.H{"message": "Переданы невалидные данные", "type": "error"})
		return
	}

	// Бизнес-логика уходит в сервис
	msg := handler.service.InitPay(userId)

	ginContext.JSON(http.StatusOK, gin.H{
		"greeting":  msg,
		"userIdLen": utf8.RuneCountInString(userId),
	})
}
