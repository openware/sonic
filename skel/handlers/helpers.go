package handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	kaigara "github.com/openware/kaigara/pkg/config"
)

// GetKaigaraConfig helper return kaigara co0nfig from gin context
func GetKaigaraConfig(ctx *gin.Context) (*kaigara.KaigaraConfig, error) {
	config, ok := ctx.MustGet("KaigaraConfig").(*kaigara.KaigaraConfig)
	if !ok {
		return nil, fmt.Errorf("Kaigara config is not found")
	}

	return config, nil
}
