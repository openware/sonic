package handlers

import (
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/openware/kaigara/pkg/vault"
	"github.com/openware/pkg/jwt"
	"github.com/openware/sonic"
)

type cache struct {
	Mutex sync.RWMutex
	Data  map[string]map[string]interface{}
}

// GetVaultConfig helper returns Vault config from gin context
func GetVaultConfig(ctx *gin.Context) (*sonic.VaultConfig, error) {
	config, ok := ctx.MustGet("VaultConfig").(*sonic.VaultConfig)
	if !ok {
		return nil, fmt.Errorf("Vault config is not found")
	}

	return config, nil
}

// GetOpendaxConfig helper return kaigara config from gin context
func GetOpendaxConfig(ctx *gin.Context) (*sonic.OpendaxConfig, error) {
	config, ok := ctx.MustGet("OpendaxConfig").(*sonic.OpendaxConfig)
	if !ok {
		return nil, fmt.Errorf("Opendax config is not found")
	}

	return config, nil
}

// GetAuth helper return auth from gin context
func GetAuth(ctx *gin.Context) (*jwt.Auth, error) {
	auth, ok := ctx.MustGet("auth").(*jwt.Auth)
	if !ok {
		return nil, fmt.Errorf("Auth is not found")
	}

	return auth, nil
}

// GetGlobalVaultService helper return global vault service from gin context
func GetGlobalVaultService(ctx *gin.Context) (*vault.Service, error) {
	vaultService, ok := ctx.MustGet("GlobalVaultService").(*vault.Service)
	if !ok {
		return nil, fmt.Errorf("Global vault service is not found")
	}

	return vaultService, nil
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
		err = vaultService.LoadSecrets(app, scope)
		if err != nil {
			panic(err)
		}

		if memoryCache.Data[app] == nil {
			memoryCache.Data[app] = make(map[string]interface{})
		}

		if memoryCache.Data[app][scope] == nil {
			memoryCache.Data[app][scope] = make(map[string]interface{})
		}

		current, err := vaultService.GetCurrentVersion(app, scope)
		if err != nil {
			panic(err)
		}

		latest, err := vaultService.GetLatestVersion(app, scope)
		if err != nil {
			panic(err)
		}

		if current != latest || firstRun {
			keys, err := vaultService.ListSecrets(app, scope)
			if err != nil {
				panic(err)
			}

			for _, key := range keys {
				val, err := vaultService.GetSecret(app, key, scope)
				if err != nil {
					panic(err)
				}
				memoryCache.Data[app][scope].(map[string]interface{})[key] = val
			}
		}
	}
}
