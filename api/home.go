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
			"commands": config.SETTING.Commands,
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

	_, rows, err := GetRows(name, uint(length), uint(start))

	if err != nil {
		c.String(http.StatusOK, err.Error())
	} else {
		count := GetCount(name)
		data := DataTable{
			Draw:            draw,
			RecordsTotal:    count,
			RecordsFiltered: count,
			Data:            rows,
		}

		c.JSON(http.StatusOK, data)
	}
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

func QueryRawSQL(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	sql := string(body)

	cols, rows, err := Query(sql)

	if err == nil {
		c.HTML(http.StatusOK, "datatable.tmpl", gin.H{
			"rows": rows,
			"cols": cols,
		})
	} else {
		c.String(http.StatusOK, err.Error())
	}

}

// todo
func AddCmd(c *gin.Context) {
	c.String(http.StatusOK, "add cmd to do")
}

func ExeCmd(c *gin.Context) {
	sid := c.Param("id")
	id, err := strconv.Atoi(sid)
	if err != nil {
		c.String(http.StatusOK, err.Error())
	} else {
		c.String(http.StatusOK, config.SETTING.Commands[id].SQL)
	}
}
