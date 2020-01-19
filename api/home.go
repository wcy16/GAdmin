package api

import (
	"fmt"
	"gadmin/config"
	"github.com/gin-gonic/gin"
	"html/template"
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
			"username": config.Get().Database.Username,
			"tables":   GetTables(),
			"commands": config.Get().Commands,
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
	var rowsAffected int64
	if err != nil {
		rowsAffected = -1
	} else {
		rowsAffected, _ = result.RowsAffected()
	}

	ret := fmt.Sprintf("%d row(s) affected. %v", rowsAffected, err)

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

func GetAddCmd(c *gin.Context) {
	c.HTML(http.StatusOK, "add_cmd.tmpl", nil)
}

func AddCmd(c *gin.Context) {
	cmd := config.Command{}

	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.String(http.StatusBadRequest, err.Error())
	} else {
		cmd.Params = len(cmd.Comments)
		id := config.AddCmd(&cmd)
		// todo save command
		c.String(http.StatusOK, fmt.Sprint(id))
	}
}

func GetCmd(c *gin.Context) {
	sid := c.Param("id")
	id, err := strconv.Atoi(sid)
	if err != nil {
		c.String(http.StatusOK, err.Error())
	} else {
		command := config.Get().Commands[id]
		card := Card{}
		card.Title = command.Name

		for i := 0; i < command.Params; i++ {

		}

		sql := fmt.Sprintf(command.SQL, command.Input...)
		card.Content = template.HTML(sql)

		link := Link{
			Name: "Execute",
			Href: "#",
		}
		card.Link = append(card.Link, link)

		c.HTML(http.StatusOK, "exe_cmd.tmpl", gin.H{
			"submit":      "/cmd/" + sid,
			"description": command.Description,
			"card":        card,
		})
	}
}

func ExeCmd(c *gin.Context) {
	var params []interface{}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.String(http.StatusOK, err.Error())
	} else {
		sid := c.Param("id")
		id, err := strconv.Atoi(sid)
		if err != nil {
			c.String(http.StatusOK, err.Error())
		} else {
			command := config.Get().Commands[id]
			sql := fmt.Sprintf(command.SQL, params...)
			if command.Query {
				cols, rows, err := Query(sql)
				if err == nil {
					c.HTML(http.StatusOK, "datatable.tmpl", gin.H{
						"rows": rows,
						"cols": cols,
					})
				} else {
					c.String(http.StatusOK, err.Error())
				}
			} else {
				result, err := db.Exec(sql)
				var rowsAffected int64
				if err != nil {
					rowsAffected = -1
				} else {
					rowsAffected, _ = result.RowsAffected()
				}

				ret := fmt.Sprintf("%d row(s) affected. %v", rowsAffected, err)

				c.String(http.StatusOK, ret)
			}
		}
	}
}
