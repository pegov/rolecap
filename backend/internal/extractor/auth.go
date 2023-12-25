package extractor

import (
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/pegov/rolecap/backend/internal/model"
	"github.com/pegov/rolecap/backend/internal/repo"
	"github.com/pegov/rolecap/backend/internal/util"
)

var (
	AccessTokenType  string = "access"
	RefreshTokenType string = "refresh"

	AccessTokenCookieName  string = "access_c"
	RefreshTokenCookieName string = "refresh_c"
)

type AuthExtractor interface {
	GetUser(c *gin.Context) (*model.AccessPayload, error)
	IsAdmin(c *gin.Context) (bool, error)
}

type authExtractor struct {
	authRepo   repo.AuthRepo
	jwtBackend util.JwtBackend
}

func NewAuthExtractor(authRepo repo.AuthRepo, jwtBackend util.JwtBackend) AuthExtractor {
	return &authExtractor{
		authRepo:   authRepo,
		jwtBackend: jwtBackend,
	}
}

func (e *authExtractor) GetUser(c *gin.Context) (*model.AccessPayload, error) {
	tokenString, err := c.Cookie(AccessTokenCookieName)
	if err != nil {
		return nil, err
	}

	claims, err := e.jwtBackend.Decode(tokenString, AccessTokenType)
	if err != nil {
		return nil, err
	}

	return model.AccessPayloadFromUserClaims(claims), nil
}

func (e *authExtractor) IsAdmin(c *gin.Context) (bool, error) {
	tokenString, err := c.Cookie(AccessTokenCookieName)
	if err != nil {
		return false, err
	}

	claims, err := e.jwtBackend.Decode(tokenString, AccessTokenType)
	if err != nil {
		return false, err
	}

	return slices.Contains(claims.Roles, "admin"), nil
}
