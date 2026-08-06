package gruenthttputil

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

const (
	TokenIssuer = "gruent"

	AccessTokenExpiry = time.Minute * 15

	RefreshTokenLength = 43
	RefreshTokenExpiry = time.Hour * 24 * 7 // 7 days

	InvitationTokenExpiry = 7 * 24 * time.Hour // 7 days
)

func GenerateAccessToken(userId string, jwtSecret string) (*string, error) {
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userId,
		"iss": TokenIssuer,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(AccessTokenExpiry).Unix(),
	}).SignedString([]byte(jwtSecret))

	if err != nil {
		return nil, err
	}

	return &token, nil
}

func GenerateRefreshToken() (*string, error) {
	refreshToken, err := gonanoid.New(RefreshTokenLength)
	if err != nil {
		return nil, err
	}

	return &refreshToken, nil
}

func GenerateInvitationToken(orgId, invitedById, email string, jwtSecret string) (*string, error) {
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":       TokenIssuer,
		"iat":       time.Now().Unix(),
		"exp":       time.Now().Add(InvitationTokenExpiry).Unix(),
		"orgId":     orgId,
		"invitedBy": invitedById,
		"email":     email,
	}).SignedString([]byte(jwtSecret))

	if err != nil {
		return nil, err
	}

	return &token, nil
}
