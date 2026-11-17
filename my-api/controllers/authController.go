package controllers

import (
	"database/sql"
	"my-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func LoginHandler(c *gin.Context) {
	var loginReq LoginRequest

	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	findUser, err := models.GetUserByEmail(loginReq.Email)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверные учётные данные"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := models.CheckPassword(findUser, loginReq.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверные учётные данные"})
		return
	}

	if findUser.Is_blocked {
		c.JSON(http.StatusForbidden, gin.H{"error": "Доступ заблокирован"})
		return
	}

	sessionID, err := models.NewSession(findUser.Id_user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("session_id", sessionID, 86400, "/", "", false, true)

	c.JSON(http.StatusOK, findUser)
}

func LogoutHandler(c *gin.Context) {
	sessionID, err := c.Cookie("session_id")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	models.DeleteSession(sessionID)
	c.SetCookie("session_id", sessionID, -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Выход выполнен"})
}

func MeHandler(c *gin.Context) {
	sessionID, err := c.Cookie("session_id")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	findSession, ok := models.GetSession(sessionID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Вы не авторизованы"})
		return
	}
	user, err := models.GetUserByID(findSession.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}
