package extractor

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pegov/rolecap/backend/internal/entity"
	mock_repo "github.com/pegov/rolecap/backend/internal/mock/repo"
	"github.com/pegov/rolecap/backend/internal/model"
	"github.com/pegov/rolecap/backend/internal/util"
	"github.com/stretchr/testify/assert"
)

var (
	TestId1       = 1
	TestUsername1 = "testusername"
)

func TestGetUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	repo := mock_repo.NewMockAuthRepo([]entity.User{})
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest("POST", "/", nil)
	c.Request = req
	pri, err := os.ReadFile("../ed25519_1.key")
	if !assert.Nil(t, err) {
		return
	}
	pub, err := os.ReadFile("../ed25519_1.pub")
	if !assert.Nil(t, err) {
		return
	}

	j := util.NewJwtBackend(pri, pub, "1")
	claims := &model.UserClaims{
		Id:       TestId1,
		Username: TestUsername1,
		Roles:    []string{},
	}
	tokenType := "access"
	token, err := j.Encode(claims, 6*time.Hour, tokenType)
	// c.SetCookie("access_c", token, 60, "/", "127.0.0.1", false, true)
	req.AddCookie(&http.Cookie{Name: "access_c", Value: token})

	extractor := NewAuthExtractor(repo, j)
	accessPayload, err := extractor.GetUser(c)

	assert.NotNil(t, accessPayload)
	assert.Equal(t, TestId1, accessPayload.Id)
	assert.Equal(t, TestUsername1, accessPayload.Username)
}
