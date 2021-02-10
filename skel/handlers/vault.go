package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/openware/kaigara/pkg/vault"
)

// GetSecrets handles GET '/secrets'
func GetSecrets(ctx *gin.Context) {
	kaigaraConfig, err := GetKaigaraConfig(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Initialize the VaultService without an appName since we'll use all of them
	vaultService := vault.NewService(kaigaraConfig.VaultAddr, kaigaraConfig.VaultToken, "global", kaigaraConfig.DeploymentID)
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

func parseScopes(scopes string) []string {
	return strings.Split(scopes, ",")
}
