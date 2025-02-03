package api

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(api *API) *gin.Engine {
	router := gin.Default()

	// File
	router.POST("/upload/:path", api.uploadFile)
	router.DELETE("/delete/:path", api.deleteFile)
	router.GET("/read/:path", api.readFile)

	// Directory
	router.POST("/directory/:path", api.createDirectory)
	router.DELETE("/directory/:path", api.deleteDirectory)
	router.GET("/list/:path", api.listFiles)

	return router
}
