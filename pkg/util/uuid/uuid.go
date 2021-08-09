package uuid

import "github.com/google/uuid"

// UUID 定义别名
type UUID = uuid.UUID

// NewUUID 创建 uuid
func NewUUID() (UUID, error) {
	return uuid.NewRandom()
}

// MustUUID 创建 uuid（如果出现问题则引发恐慌）
func MustUUID() UUID {
	v, err := NewUUID()
	if err != nil {
		panic(err)
	}
	return v
}

// MustString 创建uuid
func MustString() string {
	return MustUUID().String()
}
