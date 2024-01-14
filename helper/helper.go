package helper

import (
	"errors"
	"fmt"
	"mime/multipart"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/bzzim/glame/models"
	"github.com/bzzim/glame/pkg/utils"
	"github.com/gin-gonic/gin"
)

func LoadConfig(fileName string) (*models.Config, error) {
	data, err := utils.ReadJson[models.Config](fileName)
	if err != nil {
		return nil, err
	}
	return data, err
}

func UploadFile(file *multipart.FileHeader, ctx *gin.Context) (fileName string, err error) {
	fileType := strings.Split(file.Header.Get("Content-Type"), "/")
	if len(fileType) < 2 {
		return "", errors.New("unknown file type")
	}

	allowedTypes := []string{"jpg", "jpeg", "png", "svg", "svg+xml", "x-icon"}
	if !slices.Contains(allowedTypes, fileType[1]) {
		return "", errors.New("not allowed file type")
	}

	path := "data/uploads"
	if err := utils.MakeDir(path); err != nil {
		return "", errors.New("can't upload file")
	}

	ts := strconv.FormatInt(time.Now().Unix(), 10)
	name := fmt.Sprintf("%s--%s", ts, file.Filename)
	if err := ctx.SaveUploadedFile(file, fmt.Sprintf("%s/%s", path, name)); err != nil {
		return "", errors.New("can't save file")
	}

	return name, nil
}
