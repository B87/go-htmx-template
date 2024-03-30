package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Secret key used for signing the JWT. In a real application, ensure this key is kept secret and not hard-coded.
var hmacSampleSecret = []byte("your-256-bit-secret")

var ErrInvalidTokenFormat = fmt.Errorf("invalid token format")
var ErrExpiredToken = fmt.Errorf("expired token")
var ErrInvalidSignature = fmt.Errorf("invalid token signature")
var ErrEmpotyToken = fmt.Errorf("empty token")

type Claims struct {
	UserID   uuid.UUID `json:"user_id"`
	UserName string    `json:"user_name"`
	Role     string    `json:"role"`
	Exp      int64     `json:"exp"`
}

func (c Claims) IsAdmin() bool {
	return c.Role == string(RoleAdmin)
}

func (c Claims) IsUser() bool {
	return c.Role == string(RoleUser)
}

func (c Claims) IsExpired() bool {
	return time.Now().Unix() > c.Exp
}

func CreateJWTToken(user User) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour).Unix()
	claims := Claims{
		UserID:   user.ID,
		UserName: user.Email,
		Role:     string(user.Role),
		Exp:      expirationTime,
	}

	// Create token header
	header := base64.StdEncoding.EncodeToString([]byte(`{"alg": "HS256","typ": "JWT"}`))

	// Create token payload
	payload, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	payloadStr := base64.StdEncoding.EncodeToString(payload)

	// Create signature
	signature := computeHMAC(fmt.Sprintf("%s.%s", header, payloadStr), hmacSampleSecret)

	// Construct the token
	token := fmt.Sprintf("%s.%s.%s", header, payloadStr, signature)
	return token, nil
}

// ValidateJWTToken validates a JWT token and returns the claims if the token is valid.
func ValidateJWTToken(tokenStr string) (*Claims, error) {
	if tokenStr == "" {
		return nil, ErrEmpotyToken
	}
	parts := strings.Split(tokenStr, ".")
	if len(parts) != 3 {
		return nil, ErrInvalidTokenFormat
	}

	// Verify the signature
	signature := computeHMAC(fmt.Sprintf("%s.%s", parts[0], parts[1]), hmacSampleSecret)
	if signature != parts[2] {
		return nil, ErrInvalidSignature
	}

	// Decode the payload
	payload, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}

	var claims Claims
	err = json.Unmarshal(payload, &claims)
	if err != nil {
		return nil, err
	}
	if claims.IsExpired() {
		return nil, ErrExpiredToken
	}
	return &claims, nil
}

func computeHMAC(message string, secret []byte) string {
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
