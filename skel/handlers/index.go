package handlers

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/gin"
	"github.com/openware/kaigara/pkg/vault"
	"github.com/openware/pkg/utils"
	"github.com/openware/sonic"
)

// Version variable stores Application Version from main package
var (
	Version      string
	DeploymentID string
	memoryCache  = cache{
		Data:  make(map[string]interface{}),
		Mutex: sync.RWMutex{},
	}
	SonicPublicKey  string
	PeatioPublicKey string
	BarongPublicKey string
)

// Initialize scope which goroutine will fetch every 30 seconds
const scope = "public"

// Setup set up routes to render view HTML
func Setup(app *sonic.Runtime) {
	// Get config and env
	Version = app.Version
	DeploymentID = app.Conf.DeploymentID
	SonicPublicKey = utils.GetEnv("SONIC_PUBLIC_KEY", "")
	PeatioPublicKey = utils.GetEnv("PEATIO_PUBLIC_KEY", "")
	BarongPublicKey = utils.GetEnv("BARONG_PUBLIC_KEY", "")
	vaultConfig := app.Conf.Vault
	opendaxConfig := app.Conf.Opendax

	log.Println("DeploymentID in config:", app.Conf.DeploymentID)

	// Get app router
	router := app.Srv

	// Set up view engine
	router.HTMLRender = ginview.Default()

	// Serve static files
	router.Static("/public", "./public")

	router.GET("/", index)
	router.GET("/page", emptyPage)
	router.GET("/version", version)

	router.NoRoute(notFound)

	SetPageRoutes(router)

	// Initialize Vault Service
	vaultService := vault.NewService(vaultConfig.Addr, vaultConfig.Token, DeploymentID)

	adminAPI := router.Group("/api/v2/admin")
	adminAPI.Use(VaultServiceMiddleware(vaultService))
	adminAPI.Use(OpendaxConfigMiddleware(&opendaxConfig))
	adminAPI.Use(AuthMiddleware())
	adminAPI.Use(RBACMiddleware([]string{"superadmin"}))
	adminAPI.GET("/secrets", GetSecrets)
	adminAPI.PUT(":component/secret", SetSecret)
	adminAPI.POST("/platforms/new", CreatePlatform)

	publicAPI := router.Group("/api/v2/public")
	publicAPI.Use(VaultServiceMiddleware(vaultService))

	publicAPI.GET("/config", GetPublicConfigs)

	// Define all public env on first system start
	WriteCache(vaultService, scope, true)
	go StartConfigCaching(vaultService, scope)
}

// StartConfigCaching will fetch latest data from vault every 30 seconds
func StartConfigCaching(vaultService *vault.Service, scope string) {
	for {
		<-time.After(20 * time.Second)

		memoryCache.Mutex.Lock()
		WriteCache(vaultService, scope, false)
		memoryCache.Mutex.Unlock()
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

func notFound(ctx *gin.Context) {
	ctx.Redirect(http.StatusMovedPermanently, "/public/index.html")
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
