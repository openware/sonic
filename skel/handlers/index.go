package handlers

import (
	"net/http"

	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/gin"
	"github.com/openware/sonic"
)

// Version variable stores Application Version from main package
var Version string

// Setup set up routes to render view HTML
func Setup(app *sonic.Runtime) {

	router := app.Srv
	// Set up view engine
	router.HTMLRender = ginview.Default()
	Version = app.Version

	// Serve static files
	router.Static("/public", "./public")

	router.GET("/", index)
	router.GET("/page", emptyPage)
	router.GET("/version", version)

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

// Return application version
func version(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"Version": Version})
}
