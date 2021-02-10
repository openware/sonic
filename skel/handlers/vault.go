package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/openware/kaigara/pkg/vault"
)

// GetSecrets to handle '/secrets'
func GetSecrets(ctx *gin.Context) {
	kaigaraConfig, err := GetKaigaraConfig(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	vaultService := vault.NewService(kaigaraConfig.VaultAddr, kaigaraConfig.VaultToken, kaigaraConfig.AppName, kaigaraConfig.DeploymentID)
	scopes := parseScopes(kaigaraConfig.Scopes)
	result := map[string]map[string]interface{}{}

	for _, scope := range scopes {
		if err := vaultService.LoadSecrets(scope); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		secrets, err := vaultService.GetSecrets(scope)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		result[scope] = secrets
	}

	ctx.JSON(http.StatusOK, result)
}

func parseScopes(scopes string) []string {
	return strings.Split(scopes, ",")
}
