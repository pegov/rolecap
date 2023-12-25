package model

import (
	"errors"
	"slices"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type UserRegisterRequest struct {
	Email     string `json:"email"`
	Username  string `json:"username"`
	Password1 string `json:"password1"`
	Password2 string `json:"password2"`
}

var (
	ErrUserRegisterPasswordMismatch = errors.New("password mismatch")
)

func (data *UserRegisterRequest) Vaidate() error {
	// TODO: validate
	data.Email = strings.TrimSpace(data.Email)
	data.Username = strings.TrimSpace(data.Email)
	if data.Password1 != data.Password2 {
		return ErrUserRegisterPasswordMismatch
	}

	return nil
}

type UserCreate struct {
	Email    string `db:"email"`
	Username string `db:"username"`
	Password string `db:"password"`
}

type UserLoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserClaims struct {
	Id       int      `json:"id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	Type     string   `json:"type"`
	jwt.RegisteredClaims
}

type AccessPayload struct {
	Id       int
	Username string
	Roles    []string
}

func (data *AccessPayload) IsAdmin() bool {
	return slices.Contains(data.Roles, "admin")
}

func AccessPayloadFromUserClaims(claims *UserClaims) *AccessPayload {
	return &AccessPayload{
		Id:       claims.Id,
		Username: claims.Username,
		Roles:    claims.Roles,
	}
}
