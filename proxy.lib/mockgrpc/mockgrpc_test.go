package mockgrpc

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	proto "github.com/lancelotXie/proxy/proxy.server/grpc.server/proto"
)

func TestClientServerInit(t *testing.T) {
	size := 1000
	server := newmockServer(nil, size, nil, nil)
	client := newmockClient(server, size)
	server.client = client

	assert.Equal(t, server, client.server)
	assert.Equal(t, size, len(server.lstStream))
	assert.Equal(t, size, len(client.lstconn))
}

func TestGetDomain(t *testing.T) {
	var (
		size   = 1000
		domain = "www.baidu.com"
		ip     = "144.144.144.144"
	)
	lstDomain := make(map[string]string)
	lstDomain[domain] = ip
	server := newmockServer(nil, size, lstDomain, nil)
	client := newmockClient(server, size)
	server.client = client

	req := &proto.DomainReq{
		Domain: domain,
	}
	res, err := client.GetDomain(context.TODO(), req)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, ip, res.IP)
}

func TestDualStream(t *testing.T) {
	var (
		size        = 1000
		ip          = "144.144.144.144"
		port        = 443
		network     = "tcp"
		streamCount = 0
	)
	lstDomain := make(map[string]string)
	server := newmockServer(nil, size, lstDomain, nil)
	client := newmockClient(server, size)
	server.client = client

	cStream, err := client.StreamDual(context.TODO())
	assert.Nil(t, err)

	req, err := proto.NewStreamBytesWithoutData(ip, port, network, "1")
	if !assert.NoError(t, err) {
		return
	}
	err = cStream.Send(req)
	assert.Nil(t, err)
	streamCount++
	assert.Equal(t, streamCount, server.curStream)

	stream2, err := client.StreamDual(context.TODO())
	assert.Nil(t, err)
	err = stream2.Send(req)
	assert.Nil(t, err)
	streamCount++
	assert.Equal(t, streamCount, server.curStream)

	stream3, err := client.StreamDual(context.TODO())
	assert.Nil(t, err)
	err = stream3.Send(req)
	assert.Nil(t, err)
	streamCount++
	assert.Equal(t, streamCount, server.curStream)

	recStream := server.lstStream[streamCount-1]
	res, err := recStream.Recv()
	assert.Equal(t, req.Data, res.Data)
}

func TestDualStreamOnce(t *testing.T) {
	var (
		size        = 1000
		ip          = ""
		port        = 443
		network     = "tcp"
		streamCount = 0
	)
	lstDomain := make(map[string]string)
	server := newmockServer(nil, size, lstDomain, nil)
	client := newmockClient(server, size)
	server.client = client

	cStream, err := client.StreamDual(context.TODO())
	assert.Nil(t, err)

	req, err := proto.NewStreamBytesWithoutData(ip, port, network, "1")
	if !assert.NoError(t, err) {
		return
	}
	err = cStream.Send(req)
	assert.Nil(t, err)
	streamCount++
	assert.Equal(t, streamCount, server.curStream)

	serverStream := server.lstStream[0]
	res, err := serverStream.Recv()
	assert.Nil(t, err)
	assert.Equal(t, req.Data, res.Data)

	cStream.Send(req)
	assert.Equal(t, streamCount, server.curStream)
	stream2, err := client.StreamDual(context.Background())
	stream2.Send(req)
	streamCount++
	assert.Equal(t, streamCount, server.curStream)
}

func TestResolvLocation(t *testing.T) {
	var (
		size     = 1000
		ip       = "144.144.144.144"
		location = "testlocation"
	)
	lstDomain := make(map[string]string)
	lstDomain[ip] = location
	server := newmockServer(nil, size, lstDomain, nil)
	client := newmockClient(server, size)
	server.client = client

	ctx := context.TODO()
	req := proto.NewLocationReq(ip)
	res, err := client.ResolvLocation(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, location, res.Location)

}

func TestStream(t *testing.T) {
	var (
		msgA = "client send"
		msgB = "server write"
	)
	client2server := make(chan proto.StreamBytes, 1)
	server2client := make(chan proto.StreamBytes, 1)
	cStream := newclientstream(nil, nil, client2server, server2client, false)
	sStream := newserverstream(nil, cStream, server2client, client2server)
	cStream.Serverstream = sStream
	cStream.sendDo = nil

	cStream.Send(proto.NewStreamBytes([]byte(msgA)))
	res, err := sStream.Recv()
	assert.Nil(t, err)
	assert.Equal(t, msgA, string(res.Data))

	sStream.Send(proto.NewStreamBytes([]byte(msgB)))
	resB, err := cStream.Recv()
	assert.Nil(t, err)
	assert.Equal(t, msgB, string(resB.Data))
}

func TestWrapClientStreamAndServerStream(t *testing.T) {
	var (
		msgA = "client send"
		msgB = "server write"
	)
	cStream, sStream := NewClientStreamAndServerStream()

	cStream.Send(proto.NewStreamBytes([]byte(msgA)))
	res, err := sStream.Recv()
	assert.Nil(t, err)
	assert.Equal(t, msgA, string(res.Data))

	sStream.Send(proto.NewStreamBytes([]byte(msgB)))
	resB, err := cStream.Recv()
	assert.Nil(t, err)
	assert.Equal(t, msgB, string(resB.Data))

}
