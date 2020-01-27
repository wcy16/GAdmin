package api

import (
	"fmt"
	"gadmin/config"
	"github.com/gin-gonic/gin"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type DataTable struct {
	Draw            int        `json:"draw"`
	RecordsTotal    int        `json:"recordsTotal"`
	RecordsFiltered int        `json:"recordsFiltered"`
	Data            [][]string `json:"data"`
}

// index page
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

// get table structure
func Table(c *gin.Context) {
	name := c.Param("name")
	cols := GetCols(name)
	c.HTML(http.StatusOK, "table.tmpl", gin.H{
		"title": name,
		"cols":  cols,
		"pks":   GetPKs(name),
	})
}

// insert
func TableInsert(c *gin.Context) {
	if err := c.Request.ParseForm(); err != nil {
		c.String(http.StatusBadRequest, err.Error())
	} else {
		keys := make([]string, 0)
		vals := make([]string, 0)

		for k, v := range c.Request.PostForm {
			if v[0] == "" {
				continue
			}
			keys = append(keys, k)
			vals = append(vals, v[0]) // no array data in post
		}

		tableName := c.Param("name")
		sql := fmt.Sprintf("INSERT INTO %s (`%s`) VALUES ('%s')", tableName,
			strings.Join(keys, "`,`"),
			strings.Join(vals, "','"))

		_, err := db.Exec(sql)

		if err != nil {
			c.String(http.StatusNotAcceptable, err.Error())
		} else {
			c.Status(http.StatusOK)
		}
	}

}

// delete
func TableDel(c *gin.Context) {
	type Request struct {
		Rows [][]string
		Cols []string
	}
	req := Request{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
	} else {
		wheres := make([]string, 0)
		for _, row := range req.Rows {
			s := ""
			for idr, item := range row {
				if idr == 0 {
					s += fmt.Sprintf("(`%s` = '%s'", req.Cols[idr], item)
				} else {
					s += fmt.Sprintf(" AND `%s` = '%s'", req.Cols[idr], item)
				}
			}
			s += ")"
			wheres = append(wheres, s)
		}

		tableName := c.Param("name")
		sql := fmt.Sprintf("DELETE FROM %s WHERE (%s)", tableName, strings.Join(wheres, ") OR ("))

		_, err := db.Exec(sql)

		if err != nil {
			c.String(http.StatusNotAcceptable, err.Error())
		} else {
			c.Status(http.StatusOK)
		}
	}
}

// edit data in a table
func TableEdit(c *gin.Context) {
	type Request struct {
		Rows []string
		Cols []string
		Col  string
		Data string
	}
	req := Request{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
	} else {
		tableName := c.Param("name")
		sql := fmt.Sprintf("UPDATE %s SET `%s` = '%s' WHERE", tableName, req.Col, req.Data)

		for i := 0; i != len(req.Rows); i++ {
			if req.Rows[i] == "" {
				continue
			}
			if i == 0 {
				sql += fmt.Sprintf(" `%s` = '%s'", req.Cols[i], req.Rows[i])
			} else {
				sql += fmt.Sprintf(" AND `%s` = '%s'", req.Cols[i], req.Rows[i])
			}
		}
		_, err := db.Exec(sql)

		if err != nil {
			c.String(http.StatusNotAcceptable, err.Error())
		} else {
			c.String(http.StatusOK, req.Data)
		}
	}
}

// load data in a table
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

// page for raw sql execute
func RawSQL(c *gin.Context) {
	c.HTML(http.StatusOK, "sql.tmpl", nil)
}

// execute sql
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

// query sql
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

// page for add cmd
func GetAddCmd(c *gin.Context) {
	c.HTML(http.StatusOK, "add_cmd.tmpl", nil)
}

// add cmd
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

// page for execute cmd
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

// execute cmd
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
