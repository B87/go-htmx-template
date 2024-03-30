package handlers

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/B87/go-htmx-template/internal/auth"
)

func TestCastClaimsFromContext(t *testing.T) {

	gin.SetMode(gin.TestMode)
	ginCtx := &gin.Context{}
	ginCtx.Set("Claims", "test")

	_, err := castClaimsFromContext(ginCtx)
	assert.Equal(t, "claims could not be cast to *auth.Claims", err.Error())

	ginCtx.Set("Claims", &auth.Claims{UserName: "test"})

	claims, err := castClaimsFromContext(ginCtx)
	assert.Nil(t, err)

	assert.Equal(t, "test", claims.UserName)

}
