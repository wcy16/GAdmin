// +build dev

package gadmin

import "github.com/gin-gonic/gin"

func loadRes(router *gin.Engine, prefix string) {
	router.Static(prefix+"/static", "../static")
	router.LoadHTMLGlob("../template/*")
}
