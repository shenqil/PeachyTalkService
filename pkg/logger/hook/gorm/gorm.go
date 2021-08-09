package gorm

import (
	"encoding/json"
	"ginAdmin/pkg/logger"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"time"
)

var tableName string

// Config 配置参数
type Config struct {
	DBType       string
	DSN          string
	MaxLifetime  int
	MaxOpenConns int
	MaxIdleConns int
	TableName    string
}

// Hook gorm日志钩子
type Hook struct {
	db *gorm.DB
}

// LogItem 存储日志项
type LogItem struct {
	ID         uint      `gorm:"column:id;primary_key;auto_increment;"` // id
	Level      string    `gorm:"column:level;size:20;index;"`           // 日志级别
	TraceID    string    `gorm:"column:trace_id;size:128;index;"`       // 跟踪ID
	UserID     string    `gorm:"column:user_id;size:36;index;"`         // 用户ID
	Tag        string    `gorm:"column:tag;size:128;index;"`            // Tag
	Version    string    `gorm:"column:version;index;size:64;"`         // 版本号
	Message    string    `gorm:"column:message;size:1024;"`             // 消息
	Data       string    `gorm:"column:data;type:text;"`                // 日志数据(json)
	ErrorStack string    `gorm:"column:error_stack;type:text;"`         // Error Stack
	CreatedAt  time.Time `gorm:"column:created_at;index"`               // 创建时间
}

// New 创建基于gorm的钩子实例(需要指定表名)
func New(c *Config) *Hook {
	tableName = c.TableName

	var gormDB *gorm.DB
	var err error
	switch c.DBType {
	case "mysql":
		gormDB, err = gorm.Open(mysql.Open(c.DSN), &gorm.Config{})
	case "postgres":
		gormDB, err = gorm.Open(postgres.Open(c.DSN), &gorm.Config{})
	case "sqlite3":
		gormDB, err = gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	}

	if err != nil {
		panic(err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		panic(err)
	}

	sqlDB.SetMaxIdleConns(c.MaxIdleConns)
	sqlDB.SetMaxOpenConns(c.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(c.MaxLifetime) * time.Second)

	gormDB.AutoMigrate(new(LogItem))
	return &Hook{
		db: gormDB,
	}
}

// Exec 执行日志写入
func (h *Hook) Exec(entry *logrus.Entry) error {
	item := &LogItem{
		Level:     entry.Level.String(),
		Message:   entry.Message,
		CreatedAt: entry.Time,
	}

	data := entry.Data
	if v, ok := data[logger.TraceIDKey]; ok {
		item.TraceID, _ = v.(string)
		delete(data, logger.TraceIDKey)
	}
	if v, ok := data[logger.UserIDKey]; ok {
		item.UserID, _ = v.(string)
		delete(data, logger.UserIDKey)
	}
	if v, ok := data[logger.TagKey]; ok {
		item.Tag, _ = v.(string)
		delete(data, logger.TagKey)
	}
	if v, ok := data[logger.StackKey]; ok {
		item.ErrorStack, _ = v.(string)
		delete(data, logger.StackKey)
	}

	if v, ok := data[logger.VersionKey]; ok {
		item.Version, _ = v.(string)
		delete(data, logger.VersionKey)
	}

	if len(data) > 0 {
		b, _ := json.Marshal(data)
		item.Data = string(b)
	}

	return h.db.Create(item).Error
}

// Close 关闭钩子
func (h *Hook) Close() error {
	sqlDB, err := h.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// TableName 表名
func (LogItem) TableName() string {
	return tableName
}
