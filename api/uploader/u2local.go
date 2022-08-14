package uploader

import (
	"easy-upload/pkg/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"mime/multipart"
	"os"
	"strings"
	"time"
)

func U2LocalFS(ctx *gin.Context) {
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
				ctx.String(400, "input err! -> has '..'")
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

	process2Local(ctx, key, data.File["data"])
}

func process2Local(ctx *gin.Context, key string, data []*multipart.FileHeader) {

	if len(data) == 0 {
		ctx.String(400, "input err ! -> server can't get 'data'!")
	}
	path := fmt.Sprintf(savePathWhitDate, key, time.Now().Format(fileDateLayout))
	if ok, err := util.PathExists(path); err != nil || !ok {
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			log.Println("mkdir err!!")
			ctx.String(500, "mkdir err!!")
			return
		}
	}

	var (
		file          = data[0]
		fileSave2Path = fmt.Sprintf(save2Path, path, file.Filename)
		status        int
	)

	err := ctx.SaveUploadedFile(file, fileSave2Path)
	// 只要相对路径
	dealPath := func(path string) string {
		if strings.HasPrefix(path, savePathPrefix) {
			return path[len(savePathPrefix):]
		}
		return path
	}

	if err != nil {
		if len(data) > 1 {
			status = PartSuccess
		} else {
			status = AllSuccess
		}
	}

	respMsg := map[string]any{
		"status": status,
		"path":   dealPath(fileSave2Path),
	}

	ctx.JSON(200, respMsg)
}