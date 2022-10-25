package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"priority-task-manager/shared/pkg/services/account"
)

type ExtractAccount struct {
	service account.Service
}

func MakeExtractAccount(service account.Service) ExtractAccount {
	return ExtractAccount{
		service: service,
	}
}

func (ea ExtractAccount) Extract() gin.HandlerFunc {
	return func(context *gin.Context) {
		acc, err := ea.service.GetAccount(context.GetHeader("X-Account-Hash"))
		if err != nil {
			context.String(http.StatusUnauthorized, "Account token header missed")
		} else {
			context.Set("account", acc)
			context.Next()
		}
	}
}
