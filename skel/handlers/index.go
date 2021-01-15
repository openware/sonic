package handlers

import (
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
	fullPath, err := os.Getwd()
	if err != nil {
		log.Error("getwd", "Can't return path: "+err.Error())
	}

	cssFiles, err := WalkMatch(fullPath+"/public/assets", "*.*.css")
	if err != nil {
		log.Error("walkMatch", "Can't take list of paths for js files: "+err.Error())
	}

	jsFiles, err := WalkMatch(fullPath+"/public/assets", "*.*.js")
	if err != nil {
		log.Error("walkMatch", "Can't take list of paths for js files: "+err.Error())
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

func WalkMatch(root, pattern string) ([]string, error) {
	fullPath, err := os.Getwd()
	if err != nil {
		log.Error("getwd", "Can not return path: "+err.Error())
	}
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, strings.Replace(path, fullPath, "", -1))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}

// TODO: Add a version handler which return the value of main.Version
