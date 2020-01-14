package gadmin

import (
	"gadmin/api"
	"gadmin/config"
	"github.com/gin-gonic/gin"
)

func Serve(settingFile string) {
	config.LoadSetting(settingFile)
	api.DBConnect()

	router := gin.Default()

	loadRes(router)

	router.GET("/", api.Index)
	router.GET("/table/:name", api.Table)
	router.GET("/table/:name/data", api.LoadData)

	// Listen and serve on 0.0.0.0:8080
	router.Run(":8080")
}
