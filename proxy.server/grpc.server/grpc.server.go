package server

import (
	"context"
	"log"
	"net"
	"strconv"
	"sync"

	address "github.com/lancelotXie/proxy/proxy.lib/address"
	contransid "github.com/lancelotXie/proxy/proxy.lib/ctxtransid"
	"github.com/lancelotXie/proxy/proxy.lib/dns"
	"github.com/lancelotXie/proxy/proxy.lib/location"
	"github.com/lancelotXie/proxy/proxy.lib/logger"
	proto "github.com/lancelotXie/proxy/proxy.server/grpc.server/proto"

	errors "github.com/pkg/errors"
	"google.golang.org/grpc"
)

// grpcServer :
type grpcServer struct {
	server      *grpc.Server
	resolver    dns.Resolver
	listen      net.Listener
	serveWG     *sync.WaitGroup
	closedone   *sync.Once
	grpcerr     chan error
	err         chan error
	connections chan connWithCtx
}

type connWithCtx struct {
	conn net.Conn
	ctx  context.Context
}

func newgrpcServer(lis net.Listener) *grpcServer {
	return &grpcServer{
		server:      grpc.NewServer(),
		resolver:    dns.New(dns.Resolv),
		listen:      lis,
		serveWG:     &sync.WaitGroup{},
		closedone:   &sync.Once{},
		grpcerr:     make(chan error, 1),
		err:         make(chan error, 1),
		connections: make(chan connWithCtx, 100),
	}
}

func (s *grpcServer) start() {
	log.Println("info:监听:", s.listen.Addr().String())
	if err := s.server.Serve(s.listen); err != nil {
		s.grpcerr <- err
		close(s.grpcerr)
	}
}

// Addr :
func (s *grpcServer) Addr() net.Addr {
	return s.listen.Addr()
}

// Close :
func (s *grpcServer) Close() (err error) {
	s.closedone.Do(func() {
		err = errors.New("err:grpc server was closed")
		s.err <- err
		close(s.err)
		err = s.listen.Close()
	})
	return
}

// Accept :
func (s *grpcServer) Accept() (ctx context.Context, conn net.Conn, err error) {
	var ok bool
	select {
	case err, ok = <-s.grpcerr:
		if !ok || err == nil {
			err = errors.New("err: grpc serve loop was closed")
		}
		logger.Error(context.Background(), err)
		return
	case err, ok = <-s.err:
		if !ok || err == nil {
			err = errors.New("err: grpc server was closed")
			logger.Error(context.Background(), err)
		}
		return
	case t, ok := <-s.connections:
		if !ok || t.conn == nil {
			err = errors.New("err: grpc connection channel was closed")
			logger.Error(context.Background(), err)
			return
		}
		conn = t.conn
		ctx = t.ctx
		logger.Info(t.ctx, "accept 服务端收到grpc连接")
		log.Println("info:建立连接", conn.RemoteAddr())
		return
	}
}

/***************************
****实现StreamService接口****
****begin:*******************
****************************/

// StreamDual : 当客户端得到stream，服务端就将进入这里。
// 客户端同时得到stream和error消息，但这个error并不是 下面这个函数的返回值err。
func (s *grpcServer) StreamDual(serverstream proto.StreamService_StreamDualServer) (err error) {
	recv, err := serverstream.Recv()
	if err != nil {
		return
	}
	ip, port, network, idstring, err := proto.UnpackInfo(*recv)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	id, err := strconv.Atoi(idstring)
	if err != nil {
		if idstring != "" {
			logger.Error(context.Background(), "grpc服务端转换id出错", err)
			return
		}
		logger.Error(context.Background(), "grpc客户端请求中 缺少 id")
		return
	}
	ctx := contransid.WithID(context.Background(), id)
	ctx, cancel := context.WithCancel(ctx)
	logger.Info(ctx, "grpc服务端响应 连接请求")

	remote := address.New(ip, int(port), network)
	sconn := newstreamConn(ctx, serverstream, cancel, remote, nil)

	s.connections <- connWithCtx{conn: sconn, ctx: ctx}
	logger.Info(ctx, "连接发送至 Accept信道")

	<-ctx.Done()
	logger.Info(context.Background(), "grpc服务端请求响应完成id:", id)
	return
}

// GetDomain : 解析域名
func (s *grpcServer) GetDomain(ctx context.Context, req *proto.DomainReq) (res *proto.IPRespose, err error) {
	domain := req.Domain
	ip, err := s.resolver.Resolv(domain)
	if err != nil {
		return
	}
	res = &proto.IPRespose{IP: ip}
	log.Println("info:解析", domain, ip)
	return
}

// ResolvLocation :解析 ip 所在地,解析后 给 grpc client回复
func (s *grpcServer) ResolvLocation(ctx context.Context, req *proto.LocationReq) (res *proto.LocationResponse, err error) {
	ip := req.IP

	location, err := location.Resolv(ip)
	if err == nil {
		res = proto.NewLocationResponse(location.String())
	}
	return
}

/***************************
****实现StreamService接口****
****end:*******************
****************************/
