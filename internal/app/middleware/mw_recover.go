package middleware

import (
	"PeachyTalkService/internal/app/ginx"
	"PeachyTalkService/pkg/errors"
	"PeachyTalkService/pkg/logger"
	"bytes"
	"fmt"
	"io/ioutil"
	"runtime"

	"github.com/gin-gonic/gin"
)

var (
	dunno     = []byte("???")
	centerDot = []byte(".")
	dot       = []byte(".")
	slash     = []byte("/")
)

// RecoveryMiddleware 崩溃恢复中间件
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				stack := stack(3)
				logger.WithContext(c.Request.Context()).WithField(logger.StackKey, string(stack)).Errorf("[panic]:%v", err)
				ginx.ResError(c, errors.ErrInternalServer)
			}
		}()
		c.Next()
	}
}

// stack 返回一个格式良好的堆栈帧，跳过跳过帧。
func stack(skip int) []byte {
	buf := new(bytes.Buffer)
	//	当我们循环时，我们打开文件并读取它们。这些变量记录当前加载的文件。
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ { // 跳过预期的帧数
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		//	至少打印这么多。如果我们找不到来源，它就不会显示。
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		fmt.Fprintf(buf, "\t%s:%s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}

//source 返回第 n 行的空间修剪切片。
func source(lines [][]byte, n int) []byte {
	n-- // 在堆栈跟踪中，行是 1 索引的，但我们的数组是 0 索引的
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

//如果可能，函数返回包含 PC 的函数的名称
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	//该名称包括包的路径名，这是不必要的，
	//因为文件名已经包括在内。另外，它有中心点。
	//也就是说，我们看到
	//runtimedebug.T·ptrmethod
	//并且想要 T.ptrmethod 包路径也可能包含点
	//（例如 code.google.com...），所以先去掉路径前缀
	if lastslash := bytes.LastIndex(name, slash); lastslash >= 0 {
		name = name[lastslash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}
