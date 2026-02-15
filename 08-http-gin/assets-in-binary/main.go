package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
we need to embed static files in go binary to use them and distribute while running
*/
//lets embed assets/* and templates/*

//go:embed assets/* templates/*
var f embed.FS

func main() {
	//lets get a router
	router := gin.Default()
	templ := template.Must(template.New("").ParseFS(f, "templates/*.tmpl"))
	router.SetHTMLTemplate(templ)

	router.StaticFS("/public", http.FS(f))

	// router.GET("/", func(c *gin.Context) {
	// 	c.HTML(http.StatusOK, "index.tmpl", gin.H{
	// 		"title": "main website",
	// 	})
	// })

	// router.GET("/img", func(c *gin.Context) {
	// 	file, err := f.ReadFile("assets/cowboybepop.png")
	// 	if err != nil {
	// 		c.String(http.StatusInternalServerError, err.Error())
	// 		return
	// 	}
	// 	c.Data(
	// 		http.StatusOK,
	// 		"images/png",
	// 		file,
	// 	)
	// })

	log.Fatal(router.Run(":8080"))
}
