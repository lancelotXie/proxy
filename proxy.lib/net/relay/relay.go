package relay

import (
	"context"
	"errors"
	"net"

	"github.com/eaglexiang/go/tunnel"
	"github.com/lancelotXie/proxy/proxy.lib/logger"
)

// DialFunc 拨号函数
type DialFunc func() (conn net.Conn, err error)

type relay struct {
	ctx    context.Context
	dial   DialFunc
	src    net.Listener
	cancel context.CancelFunc
}

func newRelay(ctx context.Context, dialer DialFunc, src net.Listener) (r *relay) {
	r = new(relay)
	r.ctx = ctx
	r.dial = dialer
	r.src = src
	return
}

// 开始服务
func (r *relay) Start() (err error) {
	ctx, cancel := context.WithCancel(r.ctx)
	r.cancel = cancel

	go r.listen(ctx)

	return
}

// listen 监听转发请求
func (r *relay) listen(ctx context.Context) {
	logger.Info(ctx, "开始监听转发请求")
	for {
		conn, err := r.src.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				logger.Info(ctx, "停止监听转发请求")
				return
			default:
			}
			if errors.Is(err, net.ErrClosed) {
				logger.Info(ctx, "停止监听转发请求")
				return
			}
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				logger.Warning(ctx, "监听转发请求遇到超时错误：", err)
				continue
			}
			logger.Error(ctx, "监听转发请求发生错误：", err)
			break
		}

		go r.handle(ctx, conn)
	}
	logger.Info(ctx, "停止监听转发请求")
}

// handle 获取面向目的地址的连接，并建立与 src 之间的双向流动
func (r *relay) handle(ctx context.Context, src net.Conn) {
	dst, err := r.dial()
	if err != nil {
		logger.Warning(ctx, err)
		src.Close()
		return
	}

	tunnel.Flow(src, dst)
}

// Close 关闭服务
func (r *relay) Close() (err error) {
	err = r.src.Close()
	if r.cancel != nil {
		r.cancel()
	}
	return
}
