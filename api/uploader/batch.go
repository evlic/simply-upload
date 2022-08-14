package uploader

import (
	"easy-upload/pkg/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"mime/multipart"
	"os"
	"strings"
	"sync"
	"sync/atomic"
)

func Batch(ctx *gin.Context) {
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

	batchProcess(ctx, key, data.File["data"])
}

func batchProcess(ctx *gin.Context, key string, data []*multipart.FileHeader) {
	path := fmt.Sprintf(savePath, key)
	if ok, err := util.PathExists(path); err != nil || !ok {
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
			err := ctx.SaveUploadedFile(file, fmt.Sprintf(save2Path, path, file.Filename))
			if err != nil {
				status = PartSuccess
				atomic.AddInt64(&cnt, 1)
				log.Println("save fail")
			} else {
				urls[idx] = fmt.Sprintf(visURL, key, file.Filename)
			}
			wg.Done()
		}(d, idx)
	}
	wg.Wait()

	if cnt == n {
		status = Fail
	}

	ctx.JSON(200, map[string]any{
		"status": status,
		"cnt":    cnt,
		"urls":   urls,
	})
}