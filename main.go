package main

import (
	"easy-upload/api/uploader"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.POST("/@upload", uploader.Batch)
	r.POST("/@u2local", uploader.U2LocalFS)
	r.Run(":12333")
}