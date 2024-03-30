package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/B87/go-htmx-template/internal/auth"
)

func LoginHandler(userRepo auth.UserRepository) func(c *gin.Context) {
	return func(c *gin.Context) {
		// Handling GET request
		if c.Request.Method == http.MethodGet {
			claims, err := castClaimsFromContext(c)
			if err == nil && !claims.IsExpired() {
				// If the user is already authenticated, redirect to the home page
				c.Redirect(http.StatusSeeOther, "/")
				return
			}

			// Render the login page
			c.HTML(http.StatusOK, "login.html", nil)
			return
		}

		// Handling POST request
		if c.Request.Method == http.MethodPost {
			// Parse form data
			email := c.PostForm("email")
			password := c.PostForm("password")

			// Authenticate the user
			user, err := userRepo.Authenticate(email, password)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
				return
			}

			// Create a JWT token
			token, err := auth.CreateJWTToken(user)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
				return
			}

			// Set the token in the response as a HTTP-only cookie
			c.SetCookie("jwt", token, 72*3600, "/", "", true, true)
			c.Request.Method = http.MethodGet
			c.Redirect(http.StatusSeeOther, "/")
			return
		}
	}
}

func LogoutHandler(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		// Clear the JWT cookie
		http.SetCookie(c.Writer, NewExpiredJWTCookie())
		c.Redirect(http.StatusSeeOther, "/")
		return
	}
}

func SignupHandler(userRepo auth.UserRepository) func(c *gin.Context) {
	return func(c *gin.Context) {
		// Handling GET request
		if c.Request.Method == http.MethodGet {
			c.HTML(http.StatusOK, "signup.html", nil)
			return
		}

		// Handling POST request
		if c.Request.Method == http.MethodPost {
			// Extract the information from the form
			email := c.PostForm("email")
			password := c.PostForm("password")

			// Attempt to register the user
			err := userRepo.Create(email, password, auth.RoleUser)

			if err == auth.ErrUserAlreadyExists {
				c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
				return
			}
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
				return
			}

			// On successful registration, redirect to the login page
			c.Request.Method = http.MethodGet
			c.Redirect(http.StatusSeeOther, "/login")
			return
		}
	}
}
