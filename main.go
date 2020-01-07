package gadmin

import (
	"gadmin/api"
	"gadmin/config"
	"gadmin/static"
	gtemp "gadmin/template"
	"github.com/gin-gonic/gin"
	"html/template"
)

func Serve(settingFile string) {
	config.LoadSetting(settingFile)
	api.DBConnect()

	router := gin.Default()
	//router.LoadHTMLGlob("../template/*")
	html := template.Must(loadTemplate())

	router.SetHTMLTemplate(html)

	router.GET("/", api.Index)
	router.GET("/table/:name", api.Table)
	router.GET("/table/:name/data", api.LoadData)

	router.StaticFS("/static", static.AssetFile())

	// Listen and serve on 0.0.0.0:8080
	router.Run(":8080")
}

func loadTemplate() (*template.Template, error) {
	t := template.New("")
	for _, name := range gtemp.AssetNames() {
		asset, err := gtemp.Asset(name)
		if err != nil {
			return nil, err
		}
		t, err = t.New(name).Parse(string(asset))
	}
	return t, nil
}
