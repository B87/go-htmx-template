package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HomeData struct {
	IsAuthenticated bool
	IsAdmin         bool
	Message         string
}

// Handler function using the values set by the middleware
func HomeHandler(c *gin.Context) {
	claims, err := castClaimsFromContext(c)
	if err != nil || claims == nil {
		c.HTML(http.StatusOK, "home.html", HomeData{Message: "Hello, Guest!", IsAuthenticated: false, IsAdmin: false})
		return
	}

	c.HTML(http.StatusOK, "home.html", HomeData{
		Message: fmt.Sprintf("Hello, %s!", claims.UserName), IsAuthenticated: true, IsAdmin: claims.IsAdmin()})
}
