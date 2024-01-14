package api

import (
	"PeachyTalkService/internal/app/ginx"
	"PeachyTalkService/internal/app/schema"
	"PeachyTalkService/internal/app/service"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"io"
	"strconv"
)

var FileSet = wire.NewSet(wire.Struct(new(File), "*"))

// File 文件
type File struct {
	FileSrv *service.File
}

// Upload 上传文件
// @Tags File
// @Security ApiKeyAuth
// @Summary 上传文件
// @Accept multipart/form-data
// @Param name formData string false "文件名称"
// @Param size formData int false "文件大小"
// @Param avatar formData bool false "是否为上传头像"
// @Param file formData file true "file"
// @Success 200 {object} schema.IDResult
// @Failure 400 {object} schema.ErrorResult "{error:{code:0,message:无效的请求参数}}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/file [post]
func (a *File) Upload(c *gin.Context) {
	ctx := c.Request.Context()
	// Get MultipartReader.
	reader, err := c.Request.MultipartReader()
	if err != nil {
		ginx.ResError(c, err)
		return
	}

	item := schema.File{
		Name:     "",
		Size:     -1,
		Type:     "",
		IsAvatar: false,
		Reader:   nil,
	}

	// Traverse form fields and file streams.
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			ginx.ResError(c, err)
			return
		}

		formName := part.FormName()

		// 解析from
		if formName != "file" {
			// Get form field value
			value, err := io.ReadAll(part)
			if err != nil {
				ginx.ResError(c, err)
				return
			}
			switch formName {
			case "name":
				item.Name = string(value)
			case "size":
				size, err := strconv.ParseInt(string(value), 10, 64)
				if err == nil {
					item.Size = size
				}

			case "avatar":
				boolVal, err := strconv.ParseBool(string(value))
				if err == nil {
					item.IsAvatar = boolVal
				}
			}

		} else { // Handle file streams.
			// Get file name
			if item.Name == "" {
				item.Name = part.FileName()
			}
			item.Type = part.Header.Get("Content-Type")
			item.Reader = part
			break // 获取到文件流后直接退出，不接受后面内容
		}
	}

	result, err := a.FileSrv.Upload(ctx, item)
	if err != nil {
		ginx.ResError(c, err)
		return
	}

	ginx.ResSuccess(c, result)
}
