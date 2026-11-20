package middleware

import (
	"database/sql"
	"my-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, err := c.Cookie("session_id")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Сессия не найдена"})
			return
		}
		findSession, ok := models.GetSession(sessionID)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Вы не авторизованы"})
			return
		}
		user, err := models.GetUserByID(findSession.UserID)
		if err == sql.ErrNoRows {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не найден"})
			return
		}
		if user.Is_blocked {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Вы заблокированы"})
			return
		}
		c.Set("userID", user.Id_user)
		c.Next()
	}
}
