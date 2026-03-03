/*
 * @Author: EagleXiang
 * @LastEditors  : EagleXiang
 * @Email: eagle.xiang@outlook.com
 * @Github: https://github.com/eaglexiang
 * @Date: 2019-12-05 21:14:44
 * @LastEditTime : 2019-12-19 23:26:26
 */

package main

import (
	"context"
	"flag"
	"log"
	"net"
	_ "net/http/pprof"
	"strconv"
	"time"

	"github.com/lancelotXie/proxy/proxy.client/door"
	"github.com/lancelotXie/proxy/proxy.client/kcp"
	"github.com/lancelotXie/proxy/proxy.lib/configuration"
	config "github.com/lancelotXie/proxy/proxy.lib/configuration/grpc"
	"github.com/lancelotXie/proxy/proxy.lib/debug"
	"github.com/lancelotXie/proxy/proxy.lib/logger"
	"github.com/lancelotXie/proxy/proxy.lib/path"
	client "github.com/lancelotXie/proxy/proxy.server/grpc.client"
	kcpgo "github.com/xtaci/kcp-go"

	_ "github.com/lancelotXie/proxy/proxy.client/parsers/http"
	_ "github.com/lancelotXie/proxy/proxy.client/parsers/https"
)

const defaultProxyListenAddr = "0.0.0.0:9527"
const defaultGRPCAddr = "127.0.0.1:9529"

func main() {
	ctx := context.Background()

	dbg := flag.String("debug", "off", "是否开启调试模式")
	port := flag.Int("ctrl-port", 8085, "port for controller")
	flag.Parse()

	// 初始化系统
	initSystem(*dbg)
	defer closeSystem()
	startMetrics()

	// 开启 KCP 转发服务
	grpcAddr := resolveGRPCAddr()
	err := kcp.Start(ctx, grpcAddr)
	if err != nil {
		log.Println(err)
	}

	// 开启本地代理服务
	err = startProxy()
	if err != nil {
		log.Println(err)
	}

	err = config.Start("127.0.0.1", *port)
	if err != nil {
		log.Println(err)
	}
}

func startMetrics() {
	const interval = time.Second * 60
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		var lastBytes uint64
		for range ticker.C {
			active, bytes := door.SnapshotMetrics()
			snmp := kcpgo.DefaultSnmp.Copy()
			bps := uint64(0)
			if bytes >= lastBytes {
				bps = (bytes - lastBytes) / uint64(interval/time.Second)
			}
			lastBytes = bytes
			logger.Info(
				context.Background(),
				"metrics role=client active=", active,
				" bps=", bps,
				" kcp.retrans=", snmp.RetransSegs,
				" kcp.lost=", snmp.LostSegs,
			)
		}
	}()
}

func initSystem(dbg string) {
	debug.TryRun(dbg)

	path.Init(path.Client)

	err := configuration.Load()
	if err != nil {
		log.Println(err)
	}

	err = initGRPCClient()
	if err != nil {
		log.Println(err)
	}
}

func closeSystem() {
	config.Stop()
}

func initGRPCClient() (err error) {
	addr := resolveGRPCAddr()
	host, portStr, splitErr := net.SplitHostPort(addr)
	if splitErr != nil {
		log.Println(splitErr)
		host = "127.0.0.1"
		portStr = "9529"
	}
	port, convErr := strconv.Atoi(portStr)
	if convErr != nil {
		log.Println(convErr)
		port = 9529
	}
	log.Println("gRPC 目的地址被初始化为：", host+":"+strconv.Itoa(port))
	err = client.Init(host, port)
	return
}

// startProxy 开启本地代理服务
func startProxy() (err error) {
	lis, err := net.Listen("tcp", resolveProxyListenAddr())
	if err != nil {
		return
	}
	log.Println("本地代理服务成功监听：", lis.Addr().String())

	d := door.NewDoor(lis)
	go d.Serve()
	return
}

func resolveGRPCAddr() string {
	addr, ok := configuration.GetString("grpc.addr")
	if ok && addr != "" {
		return addr
	}
	return defaultGRPCAddr
}

func resolveProxyListenAddr() string {
	addr, ok := configuration.GetString("local.addr")
	if ok && addr != "" {
		return addr
	}
	return defaultProxyListenAddr
}
