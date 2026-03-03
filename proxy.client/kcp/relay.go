package kcp

import (
	"context"
	"net"

	"github.com/lancelotXie/proxy/proxy.lib/configuration"
	"github.com/lancelotXie/proxy/proxy.lib/logger"
	"github.com/lancelotXie/proxy/proxy.lib/net/kcp"
	"github.com/lancelotXie/proxy/proxy.lib/net/relay"

	"github.com/pkg/errors"
)

// Start 开始 KCP 转发服务
func Start(ctx context.Context, listenAddr string) (err error) {
	addr, ok := configuration.GetString("remote.addr")
	if !ok {
		err = errors.New("missing KCP remote address")
		logger.Error(ctx, "Start KCP relay failed: ", err)
		return
	}
	logger.Info(ctx, "KCP 远端地址为：", addr)

	if listenAddr == "" {
		listenAddr = "127.0.0.1:9529"
	}
	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	logger.Info(ctx, "KCP 转发服务成功监听：", listenAddr)

	dialer := func() (conn net.Conn, err error) {
		conn, err = kcp.Dial(addr)
		return
	}

	s := relay.New(ctx, dialer, lis)
	err = s.Start()
	if err != nil {
		err = errors.WithStack(err)
		logger.Error(ctx, "Start KCP relay failed: ", err)
		return
	}
	return
}
