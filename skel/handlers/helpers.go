package handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	kaigara "github.com/openware/kaigara/pkg/config"
	"github.com/openware/kaigara/pkg/vault"
)

// GetKaigaraConfig helper return kaigara co0nfig from gin context
func GetKaigaraConfig(ctx *gin.Context) (*kaigara.KaigaraConfig, error) {
	config, ok := ctx.MustGet("KaigaraConfig").(*kaigara.KaigaraConfig)
	if !ok {
		return nil, fmt.Errorf("Kaigara config is not found")
	}

	return config, nil
}

// WriteCache read latest vault version and fetch keys values from vault
// 'firstRun' variable will help to run writing to cache on first system start
// as on the start latest and current versions are the same
func WriteCache(vaultService *vault.Service, scope string, firstRun bool) {
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

		if current != latest || firstRun {
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
