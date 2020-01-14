// +build dev

package gadmin

import "github.com/gin-gonic/gin"

func loadRes(router *gin.Engine) {
	router.Static("/static", "../static")
	router.LoadHTMLGlob("../template/*")
}
