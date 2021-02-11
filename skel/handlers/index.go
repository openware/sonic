package handlers

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/gin"
	"github.com/openware/kaigara/pkg/vault"
	"github.com/openware/sonic"
)

// Version variable stores Application Version from main package
var (
	Version     string
	memoryCache map[string]map[string]interface{} = make(map[string]map[string]interface{})
)

// Setup set up routes to render view HTML
func Setup(app *sonic.Runtime) {

	router := app.Srv
	// Set up view engine
	router.HTMLRender = ginview.Default()
	Version = app.Version
	kaigaraConfig := app.Conf.KaigaraConfig

	// Serve static files
	router.Static("/public", "./public")

	router.GET("/", index)
	router.GET("/page", emptyPage)
	router.GET("/version", version)

	SetPageRoutes(router)

	vaultAPI := router.Group("/api/v2/admin")
	vaultAPI.Use(KaigaraConfigMiddleware(&kaigaraConfig))
	vaultAPI.GET("/secrets", GetSecrets)

	vaultAPI.PUT(":component/secret", SetSecret)

	vaultPublicAPI := router.Group("/api/v2/public")
	vaultPublicAPI.Use(KaigaraConfigMiddleware(&kaigaraConfig))

	vaultPublicAPI.GET("/config", GetPublicConfigs)

	// Initialize Vault Service
	vaultService := vault.NewService(kaigaraConfig.VaultAddr, kaigaraConfig.VaultToken, "global", kaigaraConfig.DeploymentID)

	go StartConfigCaching(vaultService)
}

// StartConfigCaching will fetch latest data from vault every 30 seconds
func StartConfigCaching(vaultService *vault.Service) {
	for {
		<-time.After(2 * time.Second)
		go WriteCache(vaultService, "public")
	}
}

// WriteCache read latest vault version and fetch keys values from vault
func WriteCache(vaultService *vault.Service, scope string) {
	appNames, err := vaultService.ListAppNames()
	if err != nil {
		panic(err)
	}

	for _, app := range appNames {
		vaultService.SetAppName(app)
		err = vaultService.LoadSecrets(scope)
		if err != nil {
			panic(err)
		}

		if memoryCache[app] == nil {
			memoryCache[app] = make(map[string]interface{})
		}

		if memoryCache[app][scope] == nil {
			memoryCache[app][scope] = make(map[string]interface{})
		}

		current, err := vaultService.GetCurrentVersion(scope)
		if err != nil {
			panic(err)
		}

		latest, err := vaultService.GetLatestVersion(scope)
		if err != nil {
			panic(err)
		}

		if current != latest {
			keys, err := vaultService.ListSecrets(scope)
			if err != nil {
				panic(err)
			}

			for _, key := range keys {
				val, err := vaultService.GetSecret(key, scope)
				if err != nil {
					panic(err)
				}
				memoryCache[app][scope].(map[string]interface{})[key] = val
			}
		}
	}
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
