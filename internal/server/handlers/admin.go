package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/B87/go-htmx-template/internal/auth"
)

type AdminData struct {
	IsAuthenticated bool
	IsAdmin         bool
	UserEmail       string
}

func AdminHandler() func(c *gin.Context) {
	return GeneralAdminHandler("admin.html")
}

func AdminSettings() func(c *gin.Context) {
	return GeneralAdminHandler("admin_settings.html")
}

func GeneralAdminHandler(templateName string) func(c *gin.Context) {
	return func(c *gin.Context) {
		claims, err := extractClaims(c)
		if err != nil {
			return // extractClaims already handles the response
		}
		if !checkAdminRole(c, claims) {
			return // checkAdminRole already handles the response
		}
		if c.Request.Method == http.MethodGet {
			c.HTML(http.StatusOK, templateName, AdminData{
				IsAuthenticated: true, IsAdmin: true, UserEmail: claims.UserName,
			})
			return
		}
		c.AbortWithStatus(http.StatusMethodNotAllowed)
	}
}

// Extracts claims and checks for errors.
func extractClaims(c *gin.Context) (*auth.Claims, error) {
	claims, err := castClaimsFromContext(c)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return nil, err
	}
	return claims, nil
}

// Checks if the user is an admin and aborts if not.
func checkAdminRole(c *gin.Context, claims *auth.Claims) bool {
	if !claims.IsAdmin() {
		c.Redirect(http.StatusSeeOther, "/")
		return false
	}
	return true
}
