package handlers

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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
	cssFiles, err := FilesPaths("/public/assets/*.css")
	if err != nil {
		log.Println("filePaths:", "Can't take list of paths for css files: "+err.Error())
	}

	jsFiles, err := FilesPaths("/public/assets/*.js")
	if err != nil {
		log.Println("filePaths", "Can't take list of paths for js files in public folder: "+err.Error())
	}

	ctx.HTML(http.StatusOK, "index", gin.H{
		"title":    "Index title!",
		"cssFiles": cssFiles,
		"jsFiles":  jsFiles,
		"rootID":   "root",
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

func FilesPaths(pattern string) ([]string, error) {
	var matches []string

	fullPath, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	matches, err = filepath.Glob(fullPath + pattern)
	if err != nil {
		return nil, err
	}

	for i, _ := range matches {
		matches[i] = strings.Replace(matches[i], fullPath, "", -1)
	}

	return matches, nil
}

// TODO: Add a version handler which return the value of main.Version
