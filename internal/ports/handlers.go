package ports

import (
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
)

type Service interface {
	InitPay(userId string) (string, error)
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
		ginContext.JSON(http.StatusBadRequest, gin.H{"message": ErrInvalidParams.Error(), "type": "error"})
		return
	}

	if utf8.RuneCountInString(userId) != 24 {
		ginContext.JSON(http.StatusBadRequest, gin.H{"message": ErrInvalidParams.Error(), "type": "error"})
		return
	}

	msg, err := handler.service.InitPay(userId)

	if(err != nil) {
		ginContext.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "type": "error"})
	}

	ginContext.JSON(http.StatusOK, gin.H{"message": msg, "type": "success"})
}
