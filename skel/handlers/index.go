package handlers

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/gin"
	"github.com/openware/kaigara/pkg/vault"
	"github.com/openware/pkg/mngapi/peatio"
	"github.com/openware/pkg/sonic/config"
	"github.com/openware/pkg/sonic/handlers"
	"github.com/openware/pkg/utils"
	"github.com/openware/sonic/skel/daemons"
	"github.com/openware/sonic/skel/models"
)

// Version variable stores Application Version from main package
var (
	Version      string
	DeploymentID string
)

// SonicContext stores requires client services used in handlers
type SonicContext struct {
	PeatioClient *peatio.Client
}

// Initialize scope which goroutine will fetch every 30 seconds
const scope = "public"

// Setup set up routes to render view HTML
func Setup(app *config.Runtime) {
	// Get config and env
	Version = app.Version
	DeploymentID = app.Conf.DeploymentID
	handlers.SonicPublicKey = utils.GetEnv("SONIC_PUBLIC_KEY", "")
	handlers.PeatioPublicKey = utils.GetEnv("PEATIO_PUBLIC_KEY", "")
	handlers.BarongPublicKey = utils.GetEnv("BARONG_PUBLIC_KEY", "")
	vaultConfig := app.Conf.Vault
	opendaxConfig := app.Conf.Opendax
	mngapiConfig := app.Conf.MngAPI

	peatioClient, err := peatio.New(mngapiConfig.PeatioURL, mngapiConfig.JWTIssuer, mngapiConfig.JWTAlgo, mngapiConfig.JWTPrivateKey)
	if err != nil {
		log.Printf("Can't create peatio client: " + err.Error())
		return
	}

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

	handlers.SetPageRoutes(router, &models.Page{})

	// Initialize Vault Service
	vaultService := vault.NewService(vaultConfig.Addr, vaultConfig.Token, DeploymentID)

	adminAPI := router.Group("/api/v2/admin")
	adminAPI.Use(handlers.VaultServiceMiddleware(vaultService))
	adminAPI.Use(handlers.OpendaxConfigMiddleware(&opendaxConfig))
	adminAPI.Use(handlers.AuthMiddleware())
	adminAPI.Use(handlers.RBACMiddleware([]string{"superadmin"}))
	adminAPI.Use(handlers.SonicContextMiddleware(&handlers.SonicContext{
		PeatioClient: peatioClient,
	}))

	adminAPI.GET("/secrets", handlers.GetSecrets)
	adminAPI.PUT(":component/secret", handlers.SetSecret)
	adminAPI.POST("/platforms/new", func(ctx *gin.Context) {
		handlers.CreatePlatform(ctx, daemons.CreateNewLicense, daemons.FetchConfiguration)
		return
	})

	publicAPI := router.Group("/api/v2/public")
	publicAPI.Use(handlers.VaultServiceMiddleware(vaultService))

	publicAPI.GET("/config", handlers.GetPublicConfigs)

	// Define all public env on first system start
	handlers.WriteCache(vaultService, scope, true)
	go handlers.StartConfigCaching(vaultService, scope)

	// Run LicenseRenewal
	go daemons.LicenseRenewal("finex", app, vaultService)

	// Fetch currencies and markets from the main platform periodically
	enabled, err := daemons.GetXLNEnabledFromVault(vaultService)
	if err != nil {
		log.Printf("cannot determine whether XLN is enabled: " + err.Error())
	}
	if enabled {
		go daemons.FetchConfigurationPeriodic(peatioClient, vaultService, opendaxConfig.Addr)
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
	// Any file path other than .html will be invalid.
	invalidPath := regexp.MustCompile(`^\/?((?:\w+\/)*(\w*\.[^\.html]+))`)
	if invalidPath.MatchString(ctx.Request.RequestURI) {
		ctx.Status(http.StatusNotFound)
		return
	}
	log.Printf("Path %s not found, defaulting to index.html\n", ctx.Request.URL.Path)
	index(ctx)
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

	for i := range matches {
		matches[i] = strings.Replace(matches[i], fullPath, "", -1)
	}

	return matches, nil
}
