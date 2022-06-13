package main

import (
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/gin-gonic/gin"
)

const savePath = "/home/d/public/dav/%s/%s"
const saveP = "/home/d/public/dav/%s"
const visURL = "https://d.evlic.cn/public/dav/%s/%s"

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.POST("/@upload", func(ctx *gin.Context) {
		key := ctx.PostForm("key")
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
		if key == "" {
			key = "def"
		}
		// 参数校验结束

		process(ctx, key, data.File["data"])
	})
	r.Run(":12333")
}

const (
	// AllSuccess 全部成功
	AllSuccess = iota
	// PartSuccess 部分成功
	PartSuccess
	// Fail 保存操作失败
	Fail
)

// PathExists 检查路径是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func process(ctx *gin.Context, key string, data []*multipart.FileHeader) {
	path := fmt.Sprintf(saveP, key)
	if ok, err := PathExists(path); err != nil || !ok {
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			log.Println("mkdir err!!")
			ctx.String(500, "mkdir err!!")
			return
		}
	}

	var (
		status byte
		cnt    int64
		n      = int64(len(data))
		urls   = make([]string, n)

		wg sync.WaitGroup
	)

	wg.Add(len(data))
	for idx, d := range data {
		go func(file *multipart.FileHeader, idx int) {
			err := ctx.SaveUploadedFile(file, fmt.Sprintf(savePath, key, file.Filename))
			if err != nil {
				status = PartSuccess
				go atomic.AddInt64(&cnt, 1)
				go log.Println("save fail")
			} else {
				urls[idx] = fmt.Sprintf(visURL, key, file.Filename)
			}
			wg.Done()
		}(d, idx)
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
