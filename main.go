package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"mime/multipart"
	"strings"
	"sync"
	"sync/atomic"
)

const savePath = "/home/d/public/%s/%s"
const visUrl = "https://d.evlic.cn/public/%s/%s"

func main() {
	r := gin.Default()
	r.POST("/upload", func(ctx *gin.Context) {
		key := ctx.Query("key")
		data, err := ctx.MultipartForm()
		if err != nil {
			log.Println("get file error!", err)
			ctx.String(400, "input err!")
			return
		}

		if !strings.HasPrefix(key, "..") {
			for _, d := range data.File["data"] {
				if strings.HasPrefix(d.Filename, "..") {
					ctx.String(400, "input err!")
					return
				}
			}
		} else {
			ctx.String(400, "input err!")
			return
		}
		// 参数校验结束
	})
}

const (
	AllSuccess = iota
	PartSuccess
	Fail
)

func process(ctx *gin.Context, key string, data []*multipart.FileHeader) {
	var (
		status byte
		cnt    int64
		n      = int64(len(data))
		urls   = make([]string, 0, n)

		wg sync.WaitGroup
	)

	wg.Add(len(data))
	for _, d := range data {
		go func(file *multipart.FileHeader) {
			err := ctx.SaveUploadedFile(file, fmt.Sprintf(savePath, key, file.Filename))
			if err != nil {
				status = PartSuccess
				go atomic.AddInt64(&cnt, 1)
				go log.Println("save fail")
			} else {
				urls = append(urls, fmt.Sprintf(visUrl, key, file.Filename))
			}
			wg.Done()
		}(d)
	}
	wg.Wait()

	if cnt == int64(n) {
		status = Fail
	}

	ctx.JSON(200, map[string]any{
		"status": status,
		"cnt":    cnt,
		"urls":   urls,
	})
}
