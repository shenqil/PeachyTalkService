package hook

import (
	"fmt"
	"github.com/LyricTian/queue"
	"github.com/sirupsen/logrus"
	"os"
)

//ExecCloser 将 logrus 条目写入 store 并关闭 store
type ExecCloser interface {
	Exec(entry *logrus.Entry) error
	Close() error
}

// FilterHandle 过滤器处理函数
type FilterHandle func(entry *logrus.Entry) *logrus.Entry

type options struct {
	maxQueues  int
	maxWorkers int
	extra      map[string]interface{}
	filter     FilterHandle
	levels     []logrus.Level
}

// Option 一个钩子参数选项
type Option func(*options)

var defaultOptions = options{
	maxQueues:  512,
	maxWorkers: 1,
	levels: []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
		logrus.TraceLevel,
	},
}

// SetMaxQueues 设置缓冲区的数量
func SetMaxQueues(maxQueues int) Option {
	return func(o *options) {
		o.maxQueues = maxQueues
	}
}

// SetMaxWorkers 设置工作线程数
func SetMaxWorkers(maxWorkers int) Option {
	return func(o *options) {
		o.maxWorkers = maxWorkers
	}
}

// SetExtra 设置扩展参数
func SetExtra(extra map[string]interface{}) Option {
	return func(o *options) {
		o.extra = extra
	}
}

// SetFilter 设置入口过滤器
func SetFilter(filter FilterHandle) Option {
	return func(o *options) {
		o.filter = filter
	}
}

// SetLevels 设置可用的日志级别
func SetLevels(levels ...logrus.Level) Option {
	return func(o *options) {
		o.levels = levels
	}
}

// New 创建要添加到记录器实例的钩子
func New(exec ExecCloser, opt ...Option) *Hook {
	opts := defaultOptions
	for _, o := range opt {
		o(&opts)
	}

	q := queue.NewQueue(opts.maxQueues, opts.maxWorkers)
	q.Run()

	return &Hook{
		opts: opts,
		q:    q,
		e:    exec,
	}
}

// 挂钩将日志发送到 mongo 数据库
type Hook struct {
	opts options
	q    *queue.Queue
	e    ExecCloser
}

// Levels 返回可用的日志记录级别
func (h *Hook) Levels() []logrus.Level {
	return h.opts.levels
}

// 触发日志事件时调用 Fire
func (h *Hook) Fire(entry *logrus.Entry) error {
	entry = h.copyEntry(entry)
	h.q.Push(queue.NewJob(entry, func(v interface{}) {
		h.exec(v.(*logrus.Entry))
	}))
	return nil
}
func (h *Hook) copyEntry(e *logrus.Entry) *logrus.Entry {
	entry := logrus.NewEntry(e.Logger)
	entry.Data = make(logrus.Fields)
	entry.Time = e.Time
	entry.Level = e.Level
	entry.Message = e.Message
	for k, v := range e.Data {
		entry.Data[k] = v
	}
	return entry
}
func (h *Hook) exec(entry *logrus.Entry) {
	for k, v := range h.opts.extra {
		if _, ok := entry.Data[k]; !ok {
			entry.Data[k] = v
		}
	}

	if filter := h.opts.filter; filter != nil {
		entry = filter(entry)
	}

	err := h.e.Exec(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[logrus-hook] execution error: %s", err.Error())
	}
}

// Flush 等待日志队列为空
func (h *Hook) Flush() {
	h.q.Terminate()
	h.e.Close()
}
