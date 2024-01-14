package schema

import "io"

// File 上传文件信息
type File struct {
	Name     string    // 文件名称
	Size     int64     // 文件大小
	Type     string    // 文件类型
	IsAvatar bool      // 是否是头像
	Reader   io.Reader // 文件流
}
