package routes

import (
	"github.com/gin-gonic/gin"
	"testproject/controllers"
)

func InitializeRoutes() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", controllers.Ping)
	r.GET("/compressPdf", controllers.CompressPdf)
	r.POST("/uploadAndCompressPDF", controllers.UploadAndCompressPDF)
	return r
}
