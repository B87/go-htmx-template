package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/russross/blackfriday/v2"

	"github.com/B87/go-htmx-template/internal/auth"
)

var ErrClaimsNotFound = fmt.Errorf("claims not found in context")

func castClaimsFromContext(c *gin.Context) (*auth.Claims, error) {
	claims, exists := c.Get("Claims")
	if !exists {
		return nil, ErrClaimsNotFound
	}
	parsedClaims, ok := claims.(*auth.Claims)
	if !ok {
		return nil, fmt.Errorf("claims could not be cast to *auth.Claims")
	}
	return parsedClaims, nil
}

func NewExpiredJWTCookie() *http.Cookie {
	return &http.Cookie{
		Name:     "jwt",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   true,
	}
}

func MarkdownToHTML(content string) template.HTML {
	output := blackfriday.Run([]byte(content))
	return template.HTML(output)
}
