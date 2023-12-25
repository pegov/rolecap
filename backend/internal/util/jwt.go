package util

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pegov/rolecap/backend/internal/model"
)

type JwtBackend interface {
	Encode(payload *model.UserClaims, expiration time.Duration, tokenType string) (string, error)
	Decode(tokenString string, tokenType string) (*model.UserClaims, error)
}

type jwtBackend struct {
	PrivateKey ed25519.PrivateKey
	PublicKey  ed25519.PublicKey
	CurrentKid string
}

func NewJwtBackend(privateKeyBytes []byte, publicKeyBytes []byte, currentKid string) JwtBackend {
	privateBlock, _ := pem.Decode(privateKeyBytes)
	privateParsed, err := x509.ParsePKCS8PrivateKey(privateBlock.Bytes)
	if err != nil {
		log.Fatalln(err)
	}
	privateKey := privateParsed.(ed25519.PrivateKey)

	publicBlock, _ := pem.Decode(publicKeyBytes)
	publicParsed, _ := x509.ParsePKIXPublicKey(publicBlock.Bytes)
	if err != nil {
		log.Fatalln(err)
	}
	publicKey := publicParsed.(ed25519.PublicKey)

	return &jwtBackend{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		CurrentKid: currentKid,
	}
}

func (backend *jwtBackend) Encode(claims *model.UserClaims, expiration time.Duration, tokenType string) (string, error) {
	token := jwt.New(jwt.SigningMethodEdDSA)
	token.Header["kid"] = backend.CurrentKid

	iat := time.Now()
	exp := iat.Add(expiration)
	claims.IssuedAt = jwt.NewNumericDate(iat)
	claims.ExpiresAt = jwt.NewNumericDate(exp)
	claims.Type = tokenType

	token.Claims = claims

	tokenString, err := token.SignedString(backend.PrivateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

var (
	ErrJwtDecodeMissingKid       = errors.New("jwt decode missing kid")
	ErrJwtDecodeInvalidKid       = errors.New("jwt decode invalid kid")
	ErrJwtDecodeInvalidToken     = errors.New("jwt decode invalid token")
	ErrJwtDecodeInvalidTokenType = errors.New("jwt decode invalid token type")
)

func (backend *jwtBackend) Decode(tokenString string, tokenType string) (*model.UserClaims, error) {
	var claims model.UserClaims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (any, error) {
		kid, ok := t.Header["kid"]
		if !ok {
			return nil, ErrJwtDecodeMissingKid
		}

		if kid != backend.CurrentKid {
			return nil, ErrJwtDecodeInvalidKid
		}

		return backend.PublicKey, nil
	})

	if err != nil || !token.Valid {
		return nil, ErrJwtDecodeInvalidToken
	}

	if claims.Type != tokenType {
		return nil, ErrJwtDecodeInvalidTokenType
	}

	return &claims, nil
}
