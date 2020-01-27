package gadmin

import (
	"gadmin/api"
	"gadmin/auth"
	"gadmin/config"
	"github.com/gin-gonic/gin"
)

func Serve(settingFile string, r *gin.Engine, prefix string) {
	config.LoadSetting(settingFile)
	api.DBConnect()
	config.Prefix = prefix

	loadRes(r, prefix)

	router := r.Group(prefix)

	router.GET("/signin", api.SignIn)
	router.POST("/signin", api.SignInCheck)
	router.GET("/signout", api.SignOut)

	router.Use(auth.CookieCheck())
	router.Use(auth.CookieUpdate())

	router.GET("/", api.Index)

	router.GET("/table/:name", api.Table)
	router.POST("/table/:name", api.TableInsert)
	router.PUT("/table/:name", api.TableEdit)
	router.DELETE("/table/:name", api.TableDel)
	router.GET("/table/:name/data", api.LoadData)

	router.GET("/raw_sql", api.RawSQL)
	router.POST("/exe_sql", api.ExeRawSQL)
	router.POST("/query_sql", api.QueryRawSQL)
	router.GET("/cmd/:id", api.GetCmd)
	router.POST("/cmd/:id", api.ExeCmd)
	router.GET("/cmd", api.GetAddCmd)
	router.POST("/cmd", api.AddCmd)

	// Listen and serve on 0.0.0.0:8080
	r.Run(":8080")
}
