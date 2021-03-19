package main

import (
	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	"github.com/stnc/pongo2gin"
	"github.com/telemac/goutils/files"
	"log"
	"net/http"
	"path"
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
	uri := c.Param("uri")

	fileName := path.Join("templates/", uri)
	fileExists, _ := files.FileExists(fileName)
	if !fileExists {
		c.HTML(404, "404.html", pongo2.Context{
			"file": fileName,
		})
		return
	}

	c.HTML(http.StatusOK, uri,
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
	r.GET("templates/*uri", GetPongoTemplates)
	log.Fatal(r.Run(":8888"))
}
