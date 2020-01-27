package api

import (
	"gadmin/auth"
	"gadmin/config"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SignIn(c *gin.Context) {
	c.HTML(http.StatusOK, "signin.html", gin.H{
		"prefix": config.Prefix,
	})
}

func SignInCheck(c *gin.Context) {
	type SignIn struct {
		Username string `form:"username"`
		Password string `form:"password"`
		Remember string `form:"remember"`
	}
	signIn := SignIn{}

	if err := c.ShouldBind(&signIn); err != nil {
		c.String(http.StatusNotAcceptable, err.Error())
	} else {
		if config.CheckUser(signIn.Username, signIn.Password) {
			if signIn.Remember != "" {
				auth.SetToken(c, 2)
			} else {
				auth.SetToken(c, 1)
			}
			c.String(http.StatusOK, "111")
		} else {
			c.String(http.StatusNotAcceptable, "error")
		}
	}
}

func SignOut(c *gin.Context) {
	auth.DelToken(c)
	c.Redirect(http.StatusSeeOther, config.Prefix+"/signin")
}
