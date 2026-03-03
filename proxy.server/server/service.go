package service

import (
	"context"
	"log"

	"github.com/lancelotXie/proxy/proxy.lib/logger"
	dial "github.com/lancelotXie/proxy/proxy.server/dial"
	flow "github.com/lancelotXie/proxy/proxy.server/flow"
	grpc "github.com/lancelotXie/proxy/proxy.server/grpc.server"
)

// Server : 服务端。
type server struct {
	base   grpc.Service  // 底层监听器
	Dnet   dial.Dial     // 请求网络连接
	Tunnel flow.Conn2Net // web 和 自定义协议服务端 的连接管道
}

// Start : 开始监听 请求 服务.循环直到 出错.
func (s *server) Start() (err error) {
	s.serve(s.base)
	return nil
}

// Close :	关闭服务
func (s *server) Close() error { return s.base.Close() }

// serve : 循环.从自定义协议服务端 接受到 连接:rightconn，根据rightconn的远端地址 去请求
// web 端的连接:leftconn，之后连通 leftconn 和 rightconn
func (s *server) serve(lis grpc.Service) {
	for {
		ctx, rightconn, err := lis.Accept()
		if err != nil {
			logger.Error(context.Background(), err)
			log.Println("err:accept failed", err)
			return
		}
		logger.Info(ctx, "服务端收到 grpc server 连接")
		go func() {
			leftconn, err := s.Dnet.Conn2Web(ctx, rightconn.RemoteAddr())
			if err != nil {
				logger.Error(ctx, "和web连接 出错:", err)
				rightconn.Close()
			} else {
				s.Tunnel.Trans(ctx, leftconn, rightconn)
			}
		}()
	}
}
