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
	router.GET("/raw_sql", api.RawSQL)
	router.POST("/exe_sql", api.ExeRawSQL)
	router.POST("/query_sql", api.QueryRawSQL)
	router.GET("/cmd/:id", api.ExeCmd)
	router.GET("/add_cmd", api.AddCmd)

	// Listen and serve on 0.0.0.0:8080
	router.Run(":8080")
}
