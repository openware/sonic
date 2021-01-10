package handlers

import (
	"net/http"

	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/gin"
)

// Setup set up routes to render view HTML
func Setup(router *gin.Engine) {

	// Set up view engine
	router.HTMLRender = ginview.Default()

	// Serve static files
	router.Static("/public", "./public")

	router.GET("/", index)
	router.GET("/page", emptyPage)

	SetPageRoutes(router)
}

// index render with master layer
func index(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index", gin.H{
		"title": "Index title!",
		"add": func(a int, b int) int {
			return a + b
		},
	})
}

// render only file, must full name with extension
func emptyPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "page.html", gin.H{"title": "Page file title!!"})
}

// TODO: Add a version handler which return the value of main.Version
