package logger

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync/atomic"

	"github.com/lancelotXie/proxy/proxy.lib/ctxtransid"
)

type logger struct {
	val int64
}

var infoheader = "Info:"
var warnheader = "Warn:"
var errheader = "Error:"

// Set : 设置日志 打印级别
func (l *logger) Set(level LogLevel) { l.set(level) }

func (l *logger) set(level LogLevel) {
	switch level {
	case INFO:
		atomic.StoreInt64(&l.val, int64(INFO))
	case WARN:
		atomic.StoreInt64(&l.val, int64(WARN))
	case ERROR:
		atomic.StoreInt64(&l.val, int64(ERROR))
	case NOLOG:
		atomic.StoreInt64(&l.val, int64(NOLOG))
	case _LogLevel:
		atomic.StoreInt64(&l.val, int64(_LogLevel))
	default:
		atomic.StoreInt64(&l.val, 4)
	}
}

// 判断是否 打印 该log opts,并根据log头header组合 格式输出
func (l *logger) judgelog(ctx context.Context, num LogLevel, header string, opts ...interface{}) {
	level := atomic.LoadInt64(&l.val)
	if int64(num) >= level {
		ctxinfo := ctxtransid.GetID(ctx)
		res := trans2string(ctx, header, opts...)
		log.Printf("%s\tid:%s\t%s\n", header, ctxinfo, res)
	}
}

// 将需要打印的内容转换为string
func trans2string(ctx context.Context, header string, opts ...interface{}) (res string) {
	lstparam := dealInterface(opts)
	for _, v := range lstparam {
		curstr := fmt.Sprint(v)
		res += curstr
	}
	return
}

// 拆解 不定长 入参
func dealInterface(lstparam []interface{}) (res []interface{}) {
	for _, v := range lstparam {
		if _v, ok := v.([]interface{}); ok {
			res = append(res, dealInterface(_v)...)
		} else {
			res = append(res, v)
		}
	}
	return
}

func (l *logger) Info(ctx context.Context, m ...interface{}) { l.judgelog(ctx, INFO, infoheader, m...) }

func (l *logger) Warning(ctx context.Context, m ...interface{}) {
	l.judgelog(ctx, WARN, warnheader, m...)
}

func (l *logger) Error(ctx context.Context, m ...interface{}) { l.judgelog(ctx, ERROR, errheader, m...) }

func (l *logger) SetOutput(w io.Writer) { log.SetOutput(w) }
