package api

import (
	"fmt"
	"gadmin/config"
	"github.com/gin-gonic/gin"
	"io/ioutil"
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
	cards := make([]Card, 0)
	cards = append(cards, SystemCard())

	c.HTML(http.StatusOK, "index.html", gin.H{
		"cards": cards,

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

func RawSQL(c *gin.Context) {
	c.HTML(http.StatusOK, "sql.tmpl", nil)
}

func ExeRawSQL(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	sql := string(body)

	result, err := db.Exec(sql)
	var rows_affected int64
	if err != nil {
		rows_affected = -1
	} else {
		rows_affected, _ = result.RowsAffected()
	}

	ret := fmt.Sprintf("%d row(s) affected. %v", rows_affected, err)

	c.String(http.StatusOK, ret)
}
