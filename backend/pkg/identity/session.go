package identity

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"
)

const (
	// In OWASP Session management cheatsheet (https://cheatsheetseries.owasp.org/cheatsheets/Session_Management_Cheat_Sheet.html#session-id-entropy),
	// the size recommanded is 16 so it should be fine for a while...
	tokenByteSize = 32

	accessTokenExpirationTime  = time.Minute * 15
	refreshTokenExpirationTime = time.Hour * 24 * 31 * 3
)

type Token struct {
	value     string
	expiredAt time.Time
}

func (t Token) Token() string {
	return t.value
}

func (t Token) ExpiredAt() time.Time {
	return t.expiredAt
}

func generateAccessToken() (Token, error) {
	return generateToken(accessTokenExpirationTime)
}

func generateRefreshToken() (Token, error) {
	return generateToken(refreshTokenExpirationTime)
}

func generateToken(expiration time.Duration) (Token, error) {
	token, err := generateRandomString(tokenByteSize)
	if err != nil {
		return Token{}, fmt.Errorf("failed to generate token: %w", err)
	}

	return Token{
		value:     token,
		expiredAt: time.Now().Add(expiration),
	}, nil
}

func generateRandomString(byteSize int) (string, error) {
	rawToken := make([]byte, byteSize)
	_, err := rand.Read(rawToken)
	if err != nil {
		return "", fmt.Errorf("%s: %s", "failed to read random source", err.Error())
	}

	return base64.StdEncoding.EncodeToString(rawToken), nil
}

type authClient struct {
	id              uint64
	userID          UserID
	clientReference string
}

func generateNewAuthClient(userID UserID) (authClient, error) {
	clientReference, err := generateRandomString(tokenByteSize)
	if err != nil {
		return authClient{}, err
	}

	return authClient{
		clientReference: clientReference,
		userID:          userID,
	}, nil
}
