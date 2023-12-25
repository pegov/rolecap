package router

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pegov/rolecap/backend/internal/entity"
	"github.com/pegov/rolecap/backend/internal/handler"
	mock_repo "github.com/pegov/rolecap/backend/internal/mock/repo"
	"github.com/pegov/rolecap/backend/internal/util"
	"github.com/stretchr/testify/assert"
)

var users []entity.User = []entity.User{
	{
		Id:        1,
		Email:     "test@test.ru",
		Username:  "testusername",
		Password:  "hunter2",
		Verified:  true,
		Active:    true,
		CreatedAt: time.Now(),
		LastLogin: time.Now(),
	},
}

var r *gin.Engine = gin.Default()

func TestMain(m *testing.M) {
	repo := mock_repo.NewMockAuthRepo(users)
	privateKeyBytes, err := os.ReadFile("../ed25519_1.key")
	if err != nil {
		panic(err)
	}
	publicKeyBytes, err := os.ReadFile("../ed25519_1.pub")
	if err != nil {
		panic(err)
	}
	jwtBackend := util.NewJwtBackend(privateKeyBytes, publicKeyBytes, "1")
	passwordHasher := util.NewPlainTextPasswordHasher()
	h := handler.NewAuthHandler(repo, jwtBackend, passwordHasher)
	SetupUserRouter(r, h)

	code := m.Run()
	os.Exit(code)
}

func TestLoginHandler(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]string{
		"login":    "testusername",
		"password": "hunter2",
	}
	body, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewReader(body))

	r.ServeHTTP(w, req)

	res := w.Result()
	res.Cookies()
	hasAc, hasRc := false, false
	for _, cookie := range res.Cookies() {
		if cookie.Name == "access_c" {
			hasAc = true
		}
		if cookie.Name == "refresh_c" {
			hasRc = true
		}
	}

	assert.True(t, hasAc)
	assert.True(t, hasRc)

	assert.Equal(t, 200, w.Code)
}

func TestTokenHandler(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]string{
		"login":    "testusername",
		"password": "hunter2",
	}
	body, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewReader(body))

	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	req, _ = http.NewRequest("POST", "/api/auth/token", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestRefreshTokenHandler(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]string{
		"login":    "testusername",
		"password": "hunter2",
	}
	body, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewReader(body))

	r.ServeHTTP(w, req)

	res := w.Result()
	res.Cookies()

	assert.Equal(t, 200, w.Code)

	req, _ = http.NewRequest("POST", "/api/auth/token", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	req, _ = http.NewRequest("POST", "/api/auth/token/refresh", nil)

	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code, "refresh token")
}
