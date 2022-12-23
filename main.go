package main

import (
	"github.com/gin-gonic/gin"
	"upload-cdn-service/cloudbucket"
)

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Greeting from GCS Service API"})
	})

	r.POST("/upload-file", cloudbucket.HandleFileUploadToBucket)
	r.GET("/list-file", cloudbucket.GetListFile)
	r.GET("/list-folder", cloudbucket.GetListFolder)

	r.Run()
}
