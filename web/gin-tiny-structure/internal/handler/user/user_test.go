package user

import (
	"encoding/json"
	"gdemo/internal/middleware"
	"gdemo/internal/models"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestUserSignUp(t *testing.T) {
	router := gin.Default()
	router.POST("/signup", SignUp)

	payload := models.User{
		Name:       "Jeb Blackford",
		Email:      "jeb@tesla.com",
		Age:        rand.Intn(100),
		Password:   "Mars",
		CreditCard: models.CreditCard{Number: "4596964440248549"},
	}
	payloadBytes, err := json.Marshal(payload)
	assert.Nil(t, err)

	r, err := http.NewRequest("POST", "/signup", strings.NewReader(string(payloadBytes)))
	assert.Nil(t, err)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	assert.Equal(t, 200, w.Code)

	var saved UserResponse
	assert.Nil(t, json.Unmarshal(w.Body.Bytes(), &saved))
	assert.Equal(t, payload.Name, saved.Name)
	assert.Equal(t, payload.Email, saved.Email)
	assert.Equal(t, payload.Age, saved.Age)
	assert.Equal(t, payload.CreditCard.Number, saved.CardID)
}

func getToken(t *testing.T) string {
	router := gin.Default()
	router.POST("/signin", SignIn)

	payload := SignInPayload{
		Name:     "Jeb Blackford",
		Password: "Mars",
	}
	payloadBytes, err := json.Marshal(payload)
	assert.Nil(t, err)

	r, err := http.NewRequest("POST", "/signin", strings.NewReader(string(payloadBytes)))
	assert.Nil(t, err)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	assert.Equal(t, 200, w.Code)

	var data map[string]string
	assert.Nil(t, json.Unmarshal(w.Body.Bytes(), &data))
	token, ok := data["token"]
	assert.True(t, ok)

	return token
}

func TestUserSignIn(t *testing.T) {
	getToken(t)
}

func TestCreateUser(t *testing.T) {
	router := gin.Default()
	router.POST("/api/user", middleware.Auth(), CreateUser)

	token := getToken(t)

	payload := models.User{
		Name:       "Davion Ash",
		Email:      "davion@google.com",
		Age:        rand.Intn(100),
		Password:   "Auto Drive",
		CreditCard: models.CreditCard{Number: "4538933620114913"},
	}
	payloadBytes, err := json.Marshal(payload)
	assert.Nil(t, err)

	r, err := http.NewRequest("POST", "/api/user", strings.NewReader(string(payloadBytes)))
	assert.Nil(t, err)
	r.Header.Set("Authorization", "Beare "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	assert.Equal(t, 201, w.Code)
}

func TestGetAllUsers(t *testing.T) {
	router := gin.Default()
	router.GET("/api/users", middleware.Auth(), GetAllUsers)

	token := getToken(t)
	r, err := http.NewRequest("GET", "/api/users", nil)
	assert.Nil(t, err)
	r.Header.Set("Authorization", "Beare "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	assert.Equal(t, 200, w.Code)
	var data map[string][]UserResponse
	assert.Nil(t, json.Unmarshal(w.Body.Bytes(), &data))
	users, ok := data["users"]
	assert.True(t, ok)
	assert.Equal(t, 2, len(users))
}

func TestGetUser(t *testing.T) {
	router := gin.Default()
	router.GET("/api/user/:id", middleware.Auth(), GetUser)

	token := getToken(t)
	r, err := http.NewRequest("GET", "/api/user/1", nil)
	assert.Nil(t, err)
	r.Header.Set("Authorization", "Beare "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	assert.Equal(t, 200, w.Code)
	var data UserResponse
	assert.Nil(t, json.Unmarshal(w.Body.Bytes(), &data))
	assert.Equal(t, uint(1), data.ID)
	assert.Equal(t, "Jeb Blackford", data.Name)
	assert.Equal(t, "4596964440248549", data.CardID)
}

func TestUpdateUser(t *testing.T) {
	router := gin.Default()
	router.PUT("/api/user/:id", middleware.Auth(), UpdateUser)

	payload := models.User{
		Name:       "JD Vance",
		Age:        40,
		Email:      "vance@whitehouse.gov",
		CreditCard: models.CreditCard{Number: "4937977783429292"},
	}
	payloadBytes, err := json.Marshal(payload)
	assert.Nil(t, err)

	r, err := http.NewRequest("PUT", "/api/user/2", strings.NewReader(string(payloadBytes)))
	assert.Nil(t, err)
	r.Header.Set("Authorization", "Beare "+getToken(t))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	assert.Equal(t, 200, w.Code)
	var data UserResponse
	assert.Nil(t, json.Unmarshal(w.Body.Bytes(), &data))
	assert.Equal(t, uint(2), data.ID)
	assert.Equal(t, "JD Vance", data.Name)
	assert.Equal(t, 40, data.Age)
	assert.Equal(t, "4937977783429292", data.CardID)
}

func TestDeleteUser(t *testing.T) {
	router := gin.Default()
	router.DELETE("/api/user/:id", middleware.Auth(), DeleteUser)

	r, err := http.NewRequest("DELETE", "/api/user/2", nil)
	assert.Nil(t, err)
	r.Header.Set("Authorization", "Beare "+getToken(t))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	assert.Equal(t, 204, w.Code)
}

func TestRefreshToken(t *testing.T) {
	router := gin.Default()
	router.POST("/refresh", middleware.Auth(), Refresh)

	token1 := getToken(t)
	time.Sleep(5 * time.Second)

	// valid token
	r, err := http.NewRequest("POST", "/refresh", nil)
	assert.Nil(t, err)
	r.Header.Set("Authorization", "Beare "+token1)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	assert.Equal(t, 304, w.Code)

	time.Sleep(40 * time.Second)

	// refresh token
	r, err = http.NewRequest("POST", "/refresh", nil)
	assert.Nil(t, err)
	r.Header.Set("Authorization", "Beare "+token1)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)
	assert.Equal(t, 200, w.Code)

	var data map[string]string
	assert.Nil(t, json.Unmarshal(w.Body.Bytes(), &data))
	token2, ok := data["token"]
	assert.True(t, ok)
	assert.NotEqual(t, token1, token2)

	// invalid token
	time.Sleep(20 * time.Second)

	r, err = http.NewRequest("POST", "/refresh", nil)
	assert.Nil(t, err)
	r.Header.Set("Authorization", "Beare "+token1)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)
	assert.Equal(t, 401, w.Code)
}
