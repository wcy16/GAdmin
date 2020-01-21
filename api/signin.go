package api

import (
	"gadmin/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SignIn(c *gin.Context) {
	c.HTML(http.StatusOK, "signin.html", nil)
}

func SignInCheck(c *gin.Context) {
	// todo validate
	auth.SetToken(c, 1)
	c.String(http.StatusOK, "111")
}
