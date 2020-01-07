package api

import (
	"gadmin/config"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type DataTable struct {
	Draw            int        `json:"draw"`
	RecordsTotal    int        `json:"recordsTotal"`
	RecordsFiltered int        `json:"recordsFiltered"`
	Data            [][]string `json:"data"`
}

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"sidebar": gin.H{
			"prefix":   "table",
			"username": config.SETTING.Database.Username,
			"tables":   GetTables(),
		},
	})
}

func Table(c *gin.Context) {
	name := c.Param("name")
	cols := GetCols(name)
	c.HTML(http.StatusOK, "table.tmpl", gin.H{
		"title": name,
		"cols":  cols,
	})
}

func LoadData(c *gin.Context) {
	name := c.Param("name")
	draw, _ := strconv.Atoi(c.Query("draw"))
	start, _ := strconv.Atoi(c.Query("start"))
	length, _ := strconv.Atoi(c.Query("length"))

	_, rows := GetRows(name, uint(length), uint(start))

	count := GetCount(name)
	data := DataTable{
		Draw:            draw,
		RecordsTotal:    count,
		RecordsFiltered: count,
		Data:            rows,
	}

	c.JSON(http.StatusOK, data)
}
