package deamon

import (
	ginzerolog "github.com/dn365/gin-zerolog"
	"github.com/gin-gonic/gin"
)

// NewEngine returns a new gin instance.
func NewEngine() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	if gin.Mode() != "debug" {
		r.Use(ginzerolog.Logger("gin"))
	} else {
		r.Use(gin.Logger())
	}

	r.GET("/process", GetProcList)
	return r
}
