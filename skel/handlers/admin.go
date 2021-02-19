package handlers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/openware/kaigara/pkg/vault"
)

const (
	// RequestTimeout default value to 30 seconds
	RequestTimeout = time.Duration(30 * time.Second)
)

// SetSecret handles PUT '/api/v2/admin/secret'
func SetSecret(ctx *gin.Context) {
	vaultConfig, err := GetVaultConfig(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	appName := ctx.Param("component")
	vaultService := vault.NewService(vaultConfig.Addr, vaultConfig.Token, appName, DeploymentID)

	key := ctx.PostForm("key")
	value := ctx.PostForm("value")
	scope := ctx.PostForm("scope")

	if key == "" || value == "" || scope == "" {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": "param missing (key, value or scope)"})
		return
	}

	vaultService.LoadSecrets(scope)
	err = vaultService.SetSecret(key, value, scope)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = vaultService.SaveSecrets(scope)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result, err := vaultService.GetSecret(key, scope)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

// GetSecrets handles GET '/api/v2/admin/secrets'
func GetSecrets(ctx *gin.Context) {
	vaultConfig, err := GetVaultConfig(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Initialize the VaultService without an appName since we'll use all of them
	vaultService := vault.NewService(vaultConfig.Addr, vaultConfig.Token, "global", DeploymentID)
	scopes := []string{"public", "private", "secret"}

	appNames, err := vaultService.ListAppNames()
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	result := make(map[string]map[string]interface{})

	for _, app := range appNames {
		vaultService.SetAppName(app)

		result[app] = make(map[string]interface{})

		for _, scope := range scopes {
			if err := vaultService.LoadSecrets(scope); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			result[app][scope] = make(map[string]interface{})

			if scope == "secret" {
				secretsKeys, err := vaultService.ListSecrets(scope)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				for _, key := range secretsKeys {
					result[app][scope].(map[string]interface{})[key] = "******"
				}
			} else {
				secrets, err := vaultService.GetSecrets(scope)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				result[app][scope] = secrets
			}

		}
	}

	ctx.JSON(http.StatusOK, result)
}

// CreatePlatformParams from request parameter
type CreatePlatformParams struct {
	Name        string `json:"name" binding:"requied"`
	PlatformURL string `json:"platform_url" binding:"requied"`
}

// CreatePlatform to handler '/api/v2/admin/platforms/new'
func CreatePlatform(ctx *gin.Context) {
	auth, err := GetAuth(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Allow only "superadmin" to create new platform
	if auth.Role != "superadmin" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var params CreatePlatformParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	// TODO: Get opendax-cloud url from config
	url, err := url.Parse("http://www.test.com")
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	// Request payload
	payload := map[string]interface{}{
		"pid":          params.Name,
		"uid":          auth.UID,
		"email":        auth.Email,
		"public_key":   ctx.GetHeader("PublicKey"),
		"platform_url": params.PlatformURL,
	}

	// Convert payload to json string
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	// Create new HTTP request
	req, err := http.NewRequest(http.MethodPost, url.String(), bytes.NewBuffer(jsonPayload))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Add request header
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	// Call HTTP request
	httpClient := &http.Client{Timeout: RequestTimeout}
	res, err := httpClient.Do(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer res.Body.Close()

	// Convert response body to []byte
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check for API error
	if res.StatusCode != http.StatusCreated {
		ctx.JSON(http.StatusBadRequest, resBody)
		return
	}

	ctx.JSON(http.StatusCreated, resBody)
}
