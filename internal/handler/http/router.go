package httphandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// RegisterRoutes wires all URL shortener routes onto the given Gin engine.
func RegisterRoutes(r *gin.Engine, h *URLHandler) {
	r.POST("/shorten", h.Shorten)
	r.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/:hash", h.Redirect)
}
