package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/openware/kaigara/pkg/vault"
	"github.com/openware/pkg/jwt"
	"github.com/openware/sonic"
)

// VaultConfigMiddleware middleware to set Vault config to gin context
func VaultConfigMiddleware(conf *sonic.VaultConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("VaultConfig", conf)
		c.Next()
	}
}

// OpendaxConfigMiddleware middleware to set kaigara config to gin context
func OpendaxConfigMiddleware(config *sonic.OpendaxConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("OpendaxConfig", config)
		c.Next()
	}
}

// GlobalVaultServiceMiddleware middleware to set global vault service to gin context
func GlobalVaultServiceMiddleware(vaultService *vault.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("GlobalVaultService", vaultService)
		c.Next()
	}
}

// AuthMiddleware middleware to verify bearer token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get bearer token from header
		authHeader := strings.Split(c.GetHeader("Authorization"), "Bearer ")
		if len(authHeader) != 2 {
			c.Abort()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization header not found"})
			return
		}

		// Get global vault service
		vaultService, err := GetGlobalVaultService(c)
		if err != nil {
			c.Abort()
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		key := "sonic_public_key"
		scope := "private"
		vaultService.LoadSecrets(scope)
		result, err := vaultService.GetSecret(key, scope)
		if err != nil {
			c.Abort()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if result == nil {
			c.Abort()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Sonic public key not found"})
			return
		}

		// Save public key to gin context
		c.Set("sonic_public_key", result.(string))

		// Load public key
		keyStore := jwt.KeyStore{}
		keyStore.LoadPublicKeyFromString(result.(string))

		// Parse token
		jwtToken := authHeader[1]
		auth, err := jwt.ParseAndValidate(jwtToken, keyStore.PublicKey)
		if err != nil {
			c.Abort()
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Save auth data to gin context
		c.Set("auth", &auth)

		c.Next()
	}
}

// AdminRoleMiddleware middleware to verity admin role
func AdminRoleMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth, err := GetAuth(c)
		if err != nil {
			c.Abort()
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		if !isAdminRole(auth.Role) {
			c.Abort()
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		c.Next()
	}
}

func isAdminRole(role string) bool {
	var roles = []string{"superadmin", "admin", "accountant", "compliance", "support", "technical"}

	for _, v := range roles {
		if v == role {
			return true
		}
	}

	return false
}
