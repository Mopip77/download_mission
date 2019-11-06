package middleware

import (
	"github.com/gin-gonic/gin"
	"os"
)

func BasicAuth() gin.HandlerFunc {
	return gin.BasicAuth(gin.Accounts{
		os.Getenv("APP_ID"): os.Getenv("APP_PWD"),
	})
}