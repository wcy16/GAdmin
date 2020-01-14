// +build !dev

package gadmin

import (
	"gadmin/static"
	gtemp "gadmin/template"
	"github.com/gin-gonic/gin"
	"html/template"
)

func loadRes(router *gin.Engine) {
	html := template.Must(loadTemplate())

	router.SetHTMLTemplate(html)
	router.StaticFS("/static", static.AssetFile())
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
