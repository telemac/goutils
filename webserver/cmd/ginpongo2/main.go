package main

import (
	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	"github.com/stnc/pongo2gin"
	"log"
	"net/http"
)

//GetPongoTemplates all list
func GetPongoTemplates(c *gin.Context) {
	posts := []string{
		"Larry Ellison",
		"Carlos Slim Helu",
		"Mark Zuckerberg",
		"Amancio Ortega ",
		"Jeff Bezos",
		" Warren Buffet ",
		"Bill Gates",
		"selman tunç",
	}
	// Call the HTML method of the Context to render a template
	path := c.Param("path")
	c.HTML(http.StatusOK, path,
		pongo2.Context{
			"title": "hello pongo",
			"posts": posts,
		},
	)
}

func main() {

	gin.SetMode(gin.DebugMode)
	r := gin.Default()
	r.Use(gin.Recovery())
	r.HTMLRender = pongo2gin.TemplatePath("templates")
	r.GET("templates/*path", GetPongoTemplates)
	log.Fatal(r.Run(":8888"))
}
