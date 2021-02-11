package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/openware/kaigara/pkg/vault"
)

func SetSecret(ctx *gin.Context) {
	kaigaraConfig, err := GetKaigaraConfig(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	appName := ctx.Param("component")
	vaultService := vault.NewService(kaigaraConfig.VaultAddr, kaigaraConfig.VaultToken, appName, kaigaraConfig.DeploymentID)

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
