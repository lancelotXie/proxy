/*
 * @Author: EagleXiang
 * @LastEditors: EagleXiang
 * @Email: eagle.xiang@outlook.com
 * @Github: https://github.com/eaglexiang
 * @Date: 2019-12-06 21:24:36
 * @LastEditTime: 2019-12-07 17:09:06
 */

package door

import (
	"context"
	"log"
	"net"
	"sync/atomic"

	"github.com/lancelotXie/proxy/proxy.client/parsers"

	"github.com/eaglexiang/go/tunnel"
)

// Door 流量进入的大门
type Door struct {
	lis net.Listener
}

var activeConns int64
var transitedBytes uint64

type meteredConn struct {
	net.Conn
}

func (m *meteredConn) Read(b []byte) (n int, err error) {
	n, err = m.Conn.Read(b)
	if n > 0 {
		atomic.AddUint64(&transitedBytes, uint64(n))
	}
	return
}

func SnapshotMetrics() (active int64, bytes uint64) {
	active = atomic.LoadInt64(&activeConns)
	bytes = atomic.LoadUint64(&transitedBytes)
	return
}

// NewDoor 构造新的 Door
func NewDoor(lis net.Listener) *Door {
	return &Door{
		lis: lis,
	}
}

// Serve 开始服务
func (d *Door) Serve() {
	log.Println("开始服务")

	for {
		conn, err := d.lis.Accept()
		if err != nil {
			log.Println(err)
			break
		}

		ctx := context.Background()

		go handleReq(ctx, conn)
	}
}

// handleReq 处理单个请求
func handleReq(ctx context.Context, left net.Conn) {
	err := _handleReq(ctx, left)
	if err != nil {
		log.Println(err)
	}
}

func _handleReq(ctx context.Context, left net.Conn) (err error) {
	p, renew, err := parsers.GetParser(left)
	if err != nil {
		return
	}
	defer renew.Close()

	right, err := p.Parse(ctx, renew)
	if err != nil {
		return
	}

	atomic.AddInt64(&activeConns, 1)
	defer atomic.AddInt64(&activeConns, -1)
	t := createTunnel(&meteredConn{Conn: renew}, &meteredConn{Conn: right})
	log.Println("开始流动：", t.Left().RemoteAddr(), "<>", t.Right().RemoteAddr())
	t.Flow()
	log.Println("结束流动：", t.Left().RemoteAddr(), "<>", t.Right().RemoteAddr())
	t.Close()
	tunnel.PutTunnel(t)
	return
}

// 增加函数的粒度

func createTunnel(left, right net.Conn) *tunnel.Tunnel {
	t := tunnel.GetTunnel()
	t.SetLeft(left)
	t.SetRight(right)
	return t
}
