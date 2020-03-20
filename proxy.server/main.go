package main

import (
	"context"
	"flag"
	"log"
	"net"

	configuration "github.com/lancelot/proxy/proxy.lib/configuration"
	config "github.com/lancelot/proxy/proxy.lib/configuration/grpc"
	"github.com/lancelot/proxy/proxy.lib/debug"
	mynet "github.com/lancelot/proxy/proxy.lib/net"
	path "github.com/lancelot/proxy/proxy.lib/path"
	"github.com/lancelot/proxy/proxy.server/kcp"
	service "github.com/lancelot/proxy/proxy.server/server"

	"github.com/pkg/errors"
)

// ErrLocalListenAddrNotFound 本地监听地址不存在
var ErrLocalListenAddrNotFound = errors.New("local addr not found")

func main() {
	ctx := context.Background()

	ctrlPort := initSystem()
	defer closeSystem()

	// 开启代理服务
	server, err := initProxyServer()
	if err == nil {
		go func() {
			err = server.Start()
			if err != nil {
				log.Println(errors.WithStack(err))
			}
		}()
	} else {
		log.Println(err)
	}

	// 开启 KCP 转发服务
	err = startKCPRelay(ctx)
	if err != nil {
		log.Println(err)
	}

	// 开启配置管理服务
	err = config.Start("0.0.0.0", ctrlPort)
	if err != nil {
		log.Println(err)
	}
}

func initSystem() (portInt int) {
	dbg := flag.String("debug", "off", "是否开启调试模式")
	port := flag.Int("ctrl-port", 8085, "port for controller")
	flag.Parse()

	debug.TryRun(*dbg)
	path.Init(path.Server)

	err := configuration.Load()
	if err != nil {
		log.Println(errors.WithStack(err))
	}
	portInt = *port
	return
}

func closeSystem() {
	config.Stop()
}

func initProxyServer() (server mynet.Service, err error) {
	addr, ok := configuration.GetString("listen.addr")
	if !ok {
		err = errors.WithStack(ErrLocalListenAddrNotFound)
		return
	}

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	server = service.NewService(lis)
	return
}

// startKCPRelay 开启 kcp 转发
func startKCPRelay(ctx context.Context) (err error) {
	addr, ok := configuration.GetString("listen.addr")
	if !ok {
		err = errors.WithStack(ErrLocalListenAddrNotFound)
		return
	}

	err = kcp.Start(ctx, addr)
	return
}
