package api

import (
	"fmt"
	"html/template"
	"runtime"
)

type Card struct {
	Title   string
	Content template.HTML
	Link    []Link
}

type Link struct {
	Name string
	Href string
}

func SystemCard() Card {
	card := Card{}

	card.Title = "System Info"

	var dbPing string

	if err := db.Ping(); err != nil {
		dbPing = err.Error()
	} else {
		dbPing = "connected"
	}

	card.Content = template.HTML(fmt.Sprintf("<b>Database:</b> %s<br><b>System:</b> %s",
		dbPing, runtime.GOOS))

	return card
}
