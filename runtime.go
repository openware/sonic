package sonic

import (
	"log"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/openware/pkg/database"
	"gorm.io/gorm"
)

// Runtime configuration of the application
type Runtime struct {
	Conf Config

	Srv     *gin.Engine
	Version string

	dbOnce sync.Once
	db     *gorm.DB
}

func (r *Runtime) GetDB() *gorm.DB {
	r.dbOnce.Do(func() {
		var err error
		r.db, err = database.Connect(&r.Conf.Database)
		if err != nil {
			log.Fatal("could not open connection to database", err.Error())
		}
	})

	return r.db
}
