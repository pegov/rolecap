package router

import (
	"github.com/gin-gonic/gin"
	"github.com/pegov/rolecap/backend/internal/handler"
)

func SetupUserRouter(r *gin.Engine, h handler.AuthHandler) {
	// r.POST("/api/auth/register", h.Register)
	r.POST("/api/auth/login", h.Login)
	r.POST("/api/auth/token", h.Token)
	r.POST("/api/auth/token/refresh", h.RefreshToken)
}
