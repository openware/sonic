package sonic

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Runtime configuration of the application
type Runtime struct {
	Conf    Config
	DB      *gorm.DB
	Srv     *gin.Engine
	Version string
}
