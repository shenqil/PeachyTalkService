package schema

type AccessType string

const (
	Subscribe AccessType = "1"
	Publish   AccessType = "2"
)

// IMClient IM 客户端
type IMClient struct {
	ClientID string `json:"clientId"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password"`
}

// IMAcl IM 权限
type IMAcl struct {
	Access     AccessType `json:"access"`
	Username   string     `json:"username" binding:"required"`
	ClientID   string     `json:"clientId"`
	IPAddr     string     `json:"ipAddr"`
	Topic      string     `json:"topic"`
	MountPoint string     `json:"mountPoint"`
}
