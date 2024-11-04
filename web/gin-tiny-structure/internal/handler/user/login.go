package user

import (
	"log/slog"
	"net/http"
	"time"

	"gdemo/hash"
	"gdemo/internal/models"
	"gdemo/jwthelper"

	"github.com/gin-gonic/gin"
)

type SignInPayload struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// SignIn doc
//
//	@Description	user sign in
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		SignInPayload	true	"User sign in payload"
//	@Success		200		{object}	gin.H			"{token:string}"
//	@Failure		400		{object}	gin.H
//	@Failure		401		{object}	gin.H
//	@Failure		500		{object}	gin.H
//	@Router			/signin [post]
func SignIn(c *gin.Context) {
	var body SignInPayload
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := user.GetByName(body.Name); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !hash.VerifyPassword(user.Password, body.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	if tokenString, err := jwthelper.GenerateToken(); err != nil {
		slog.Error("SignIn - jwthelper.GenerateToken", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
	} else {
		c.JSON(http.StatusOK, gin.H{"token": tokenString})
	}
}

// SignUp godoc
//
//	@Description	user sign up
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		models.User	true	"User sign up payload"
//	@Success		200		{object}	UserResponse
//	@Failure		400		{object}	gin.H
//	@Failure		500		{object}	gin.H
//	@Router			/signup [post]
func SignUp(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := user.Create(); err != nil {
		slog.Error("SignUp - user.Create", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}

	var resp UserResponse
	resp.FillFromModel(user)
	c.JSON(http.StatusOK, resp)
}

// Refresh godoc
//
//	@Description	refresh your token
//	@Tags			user
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		200	{object}	gin.H	"{token:string}"
//	@Success		304
//	@Failure		401	{object}	gin.H
//	@Failure		500	{object}	gin.H
//	@Router			/refresh [post]
func Refresh(c *gin.Context) {
	exp := c.GetFloat64("exp")
	if exp == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "zero expiration"})
		return
	}

	expiration := time.Unix(int64(exp), 0)
	slog.Info("Refresh", "expiration", expiration)
	remain := time.Until(expiration)
	slog.Info("Refresh", "remain", remain)

	if remain >= 30*time.Second {
		c.Status(http.StatusNotModified)
	} else if remain > 0 {
		if tokenString, err := jwthelper.GenerateToken(); err != nil {
			slog.Error("Refresh - jwthelper.GenerateToken", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		} else {
			c.JSON(http.StatusOK, gin.H{"token": tokenString})
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token has invalid claims: token is expired"})
	}
}
