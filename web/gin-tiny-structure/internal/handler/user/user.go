package user

import (
	"log/slog"
	"net/http"
	"strconv"

	"gdemo/internal/models"

	"github.com/gin-gonic/gin"
)

type UserResponse struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email" binding:"email"`
	Age    int    `json:"age"`
	CardID string `json:"card_id"`
}

func (r *UserResponse) FillFromModel(u models.User) {
	r.ID = u.ID
	r.Name = u.Name
	r.Email = u.Email
	r.Age = u.Age
	r.CardID = u.CreditCard.Number
}

// GetAllUsers returns all of the persistent users.
//
//	@Description	get all users
//	@Tags			user
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		200	{array}		UserResponse
//	@Failure		401	{object}	gin.H
//	@Failure		500	{object}	gin.H
//	@Router			/api/users [get]
func GetAllUsers(c *gin.Context) {
	var user models.User
	if users, err := user.All(); err == nil {
		resp := make([]UserResponse, len(users))
		for i, u := range users {
			resp[i].FillFromModel(*u)
		}
		c.JSON(http.StatusOK, resp)
	} else {
		slog.Error("GetAllUsers - user.All", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
	}
}

// CreateUser creates a user.
//
//	@Description	create a user
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			payload	body	models.User	true	"Add user"
//	@Security		ApiKeyAuth
//	@Success		201
//	@Failure		400	{object}	gin.H
//	@Failure		401	{object}	gin.H
//	@Failure		500	{object}	gin.H
//	@Router			/api/user [post]
func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := user.Create(); err != nil {
		slog.Error("CreateUser - user.Create", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}

	c.Status(http.StatusCreated)
}

// GetUser gets a specific user by the provided user ID.
//
//	@Summary		Show a user
//	@Description	Get a user
//	@Tags			user
//	@Produce		json
//	@Param			id	path	int	true	"User ID"
//	@Security		ApiKeyAuth
//	@Success		200	{object}	UserResponse
//	@Failure		400	{object}	gin.H
//	@Failure		401	{object}	gin.H
//	@Router			/api/user/{id} [get]
func GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id:" + err.Error()})
		return
	}

	var user models.User
	if err := user.Get(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		var resp UserResponse
		resp.FillFromModel(user)
		c.JSON(http.StatusOK, resp)
	}
}

func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	acid, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id:" + err.Error()})
		return
	}
	if user.ID != 0 && uint(acid) != user.ID {
		slog.Warn("UpdateUser", "payload.ID", user.ID, "path.ID", id)
	}

	if err := user.Update(acid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		user.Get(acid)
		var resp UserResponse
		resp.FillFromModel(user)
		c.JSON(http.StatusOK, resp)
	}
}

// DeleteUser deletes a specific user by the provided user ID.
//
//	@Summary		Delete a user
//	@Description	Delete a user
//	@Tags			user
//	@Produce		json
//	@Param			id	path	int	true	"User ID"
//	@Security		ApiKeyAuth
//	@Success		204
//	@Failure		400	{object}	gin.H
//	@Failure		401	{object}	gin.H
//	@Router			/api/user/{id} [delete]
func DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id:" + err.Error()})
		return
	}

	var user models.User
	if err := user.Delete(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.Status(http.StatusNoContent)
	}
}
