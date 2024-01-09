package handler

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pegov/rolecap/backend/internal/extractor"
	"github.com/pegov/rolecap/backend/internal/model"
	"github.com/pegov/rolecap/backend/internal/repo"
	"github.com/pegov/rolecap/backend/internal/util"
)

type AuthHandler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	Logout(c *gin.Context)
	Token(c *gin.Context)
	RefreshToken(c *gin.Context)
}

type authHandler struct {
	Repo           repo.AuthRepo
	JwtBackend     util.JwtBackend
	PasswordHasher util.PasswordHasher
	Logger         *slog.Logger
}

func NewAuthHandler(
	authRepo repo.AuthRepo,
	jwtBackend util.JwtBackend,
	passwordHasher util.PasswordHasher,
	logger *slog.Logger,
) AuthHandler {
	return &authHandler{authRepo, jwtBackend, passwordHasher, logger}
}

func (s *authHandler) Register(c *gin.Context) {
	var data model.UserRegisterRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(422, gin.H{"detail": "Unprocessable entity"})
		return
	}

	if err := data.Validate(); err != nil {
		c.JSON(400, gin.H{"detail": err.Error()})
		return
	}

	existingEmail, err := s.Repo.GetByEmail(data.Email)
	if err != nil {
		c.JSON(500, gin.H{"detail": "Internal server error"})
		return
	}
	if existingEmail != nil {
		c.JSON(400, gin.H{"detail": "email already exists"})
		return
	}

	existingUsername, err := s.Repo.GetByUsername(data.Username)
	if err != nil {
		c.JSON(500, gin.H{"detail": "Internal server error"})
		return
	}
	if existingUsername != nil {
		c.JSON(400, gin.H{"detail": "username already exists"})
		return
	}
	hashedPasswordBytes, err := s.PasswordHasher.Hash([]byte(data.Password1))
	if err != nil {
		c.JSON(500, gin.H{"detail": "Internal server error"})
		return
	}

	userCreate := model.UserCreate{
		Email:    data.Email,
		Username: data.Username,
		Password: string(hashedPasswordBytes),
	}

	entity, err := s.Repo.Create(userCreate)
	if err != nil {
		c.JSON(500, gin.H{"detail": "Internal server error"})
		return
	}

	claims := &model.UserClaims{
		Id:       entity.Id,
		Username: entity.Username,
		Roles:    []string{},
	}

	claims.Type = extractor.AccessTokenType
	accessToken, err := s.JwtBackend.Encode(claims, 6*time.Hour, extractor.AccessTokenType)
	if err != nil {
		c.JSON(500, gin.H{"detail": "Internal server error"})
		return
	}
	claims.Type = extractor.RefreshTokenType
	refreshToken, err := s.JwtBackend.Encode(claims, 31*24*time.Hour, extractor.RefreshTokenType)
	if err != nil {
		c.JSON(500, gin.H{"detail": "Internal server error"})
		return
	}

	c.SetCookie(extractor.AccessTokenCookieName, accessToken, 6*60*60, "/", "127.0.0.1", false, true)
	c.SetCookie(extractor.RefreshTokenCookieName, refreshToken, 31*24*60*60, "/", "127.0.0.1", false, true)
}

func (s *authHandler) Login(c *gin.Context) {
	var data model.UserLoginRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(422, gin.H{"detail": "Unprocessable entity"})
		return
	}

	entity, err := s.Repo.GetByLogin(data.Login)
	if err != nil {
		c.JSON(500, gin.H{"detail": "Internal server error"})
		return
	}

	if entity == nil {
		c.JSON(404, gin.H{"detail": "Not found"})
		return
	}

	err = s.PasswordHasher.Compare([]byte(entity.Password), []byte(data.Password))
	if err != nil {
		c.JSON(401, gin.H{"detail": "Not authenticated"})
		return
	}

	claims := &model.UserClaims{
		Id:       entity.Id,
		Username: entity.Username,
		Roles:    []string{},
	}

	claims.Type = extractor.AccessTokenType
	accessToken, err := s.JwtBackend.Encode(claims, 6*time.Hour, extractor.AccessTokenType)
	if err != nil {
		c.JSON(500, gin.H{"detail": "Internal server error"})
		return
	}
	claims.Type = extractor.RefreshTokenType
	refreshToken, err := s.JwtBackend.Encode(claims, 31*24*time.Hour, extractor.RefreshTokenType)
	if err != nil {
		c.JSON(500, gin.H{"detail": "Internal server error"})
		return
	}

	c.SetCookie(extractor.AccessTokenCookieName, accessToken, 6*60*60, "/", "127.0.0.1", false, true)
	c.SetCookie(extractor.RefreshTokenCookieName, refreshToken, 31*24*60*60, "/", "127.0.0.1", false, true)
}

func (s *authHandler) Logout(c *gin.Context) {
	c.SetCookie(extractor.AccessTokenCookieName, "", -1, "/", "127.0.0.1", false, true)
	c.SetCookie(extractor.RefreshTokenCookieName, "", -1, "/", "127.0.0.1", false, true)
}

func (s *authHandler) Token(c *gin.Context) {
	tokenString, err := c.Cookie(extractor.AccessTokenCookieName)
	if err != nil {
		c.JSON(401, gin.H{"detail": "Not authenticated"})
		return
	}

	claims, err := s.JwtBackend.Decode(tokenString, extractor.AccessTokenType)
	if err != nil {
		c.JSON(401, gin.H{"detail": "Not authenticated"})
		return
	}

	c.JSON(200, claims)
}

func (s *authHandler) RefreshToken(c *gin.Context) {
	tokenString, err := c.Cookie(extractor.RefreshTokenCookieName)
	if err != nil {
		c.JSON(401, gin.H{"detail": "Not authenticated"})
		return
	}

	claims, err := s.JwtBackend.Decode(tokenString, extractor.RefreshTokenType)
	if err != nil {
		c.JSON(401, gin.H{"detail": "Not authenticated"})
		return
	}

	entity, err := s.Repo.GetById(claims.Id)
	if err != nil {
		c.JSON(500, gin.H{"detail": "Internal server error"})
		return
	}

	if entity == nil {
		c.JSON(404, gin.H{"detail": "Not found"})
		return
	}

	newClaims := &model.UserClaims{
		Id:       entity.Id,
		Username: entity.Username,
		Roles:    []string{},
	}
	newClaims.Type = extractor.AccessTokenType

	newAccessToken, err := s.JwtBackend.Encode(newClaims, 6*time.Hour, extractor.AccessTokenType)
	if err != nil {
		c.JSON(500, gin.H{"detail": "Internal server error"})
		return
	}

	c.SetCookie(extractor.AccessTokenCookieName, newAccessToken, 6*60*60, "/", "127.0.0.1", false, true)
}
