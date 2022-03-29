package logger

import (
	"context"
	"fmt"
	"ginAdmin/internal/app/config"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

// Define key
const (
	TraceIDKey = "trace_id"
	UserIDKey  = "user_id"
	TagKey     = "tag"
	VersionKey = "version"
	StackKey   = "stack"
)

var (
	version string
)

// Logger Logrus
type Logger = logrus.Logger

// Entry Logrus entry
type Entry = logrus.Entry

// StandardLogger 获取标准日志
func StandardLogger() *Logger {
	return logrus.StandardLogger()
}

type (
	traceIDKey struct{}
	userIDKey  struct{}
	tagKey     struct{}
	stackKey   struct{}
)

// InitLogger 初始化日志模块
func InitLogger() (func(), error) {
	c := config.C.Log

	// 日志级别
	logrus.SetLevel(logrus.Level(c.Level))

	// 日志格式
	if c.Format == "json" {
		logrus.SetFormatter(new(logrus.JSONFormatter))
	} else {
		logrus.SetFormatter(new(logrus.TextFormatter))
	}

	// 设定日志输出
	var file *os.File
	if c.Output != "" {
		switch c.Output {
		case "stdout":
			logrus.SetOutput(os.Stdout)
		case "stderr":
			logrus.SetOutput(os.Stderr)
		case "file":
			if name := c.OutputFile; name != "" {
				_ = os.MkdirAll(filepath.Dir(name), 0777)

				f, err := os.OpenFile(name, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
				if err != nil {
					return nil, err
				}
				logrus.SetOutput(f)
				file = f
			}
		}
	}

	return func() {
		if file != nil {
			file.Close()
		}
	}, nil
}

// 设置版本号
func SetVersion(v string) {
	version = v
}

// NewTraceIDContext 创建跟踪ID上下文
func NewTraceIDContext(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey{}, traceID)
}

// FromTraceIDContext 从上下文中获取跟踪ID
func FromTraceIDContext(ctx context.Context) string {
	v := ctx.Value(traceIDKey{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// NewUserIDContext 创建用户ID上下文
func NewUserIDContext(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey{}, userID)
}

// FromUserIDContext 从上下文中获取用户ID
func FromUserIDContext(ctx context.Context) string {
	v := ctx.Value(userIDKey{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// NewTagContext 创建Tag上下文
func NewTagContext(ctx context.Context, tag string) context.Context {
	return context.WithValue(ctx, tagKey{}, tag)
}

// FromTagContext 从上下文中获取Tag
func FromTagContext(ctx context.Context) string {
	v := ctx.Value(tagKey{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// NewStackContext 创建Stack上下文
func NewStackContext(ctx context.Context, stack error) context.Context {
	return context.WithValue(ctx, stackKey{}, stack)
}

// FromStackContext 从上下文中获取Stack
func FromStackContext(ctx context.Context) error {
	v := ctx.Value(stackKey{})
	if v != nil {
		if s, ok := v.(error); ok {
			return s
		}
	}
	return nil
}

// WithContext Use context create entry
func WithContext(ctx context.Context) *Entry {
	if ctx == nil {
		ctx = context.Background()
	}

	fields := map[string]interface{}{
		VersionKey: version,
	}

	if v := FromTraceIDContext(ctx); v != "" {
		fields[TraceIDKey] = v
	}

	if v := FromUserIDContext(ctx); v != "" {
		fields[UserIDKey] = v
	}

	if v := FromTagContext(ctx); v != "" {
		fields[TagKey] = v
	}

	if v := FromStackContext(ctx); v != nil {
		fields[StackKey] = fmt.Sprintf("%+v", v)
	}

	return logrus.WithContext(ctx).WithFields(fields)
}

// Define logrus alias
var (
	Tracef = logrus.Tracef
	Debugf = logrus.Debugf
	Infof  = logrus.Infof
	Warnf  = logrus.Warnf
	Errorf = logrus.Errorf
	Fatalf = logrus.Fatalf
	Panicf = logrus.Panicf
	Printf = logrus.Printf
)
