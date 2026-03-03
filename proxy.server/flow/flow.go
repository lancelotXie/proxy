package flow

import (
	"context"
	"net"
	"sync/atomic"

	"github.com/lancelotXie/proxy/proxy.lib/ctxtransid"
	"github.com/lancelotXie/proxy/proxy.lib/logger"
	tunnel "github.com/eaglexiang/go/tunnel"
)

// Conn2Net : 连接2个net.Conn
type Conn2Net interface {
	Close() (err error)
	Trans(context.Context, net.Conn, net.Conn)
}

// Conn2Conn : 连通2个net.Conn
type Conn2Conn struct{}

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

// Close :
func (c *Conn2Conn) Close() error { return nil }

// Trans : 连通2个net.Conn
func (c *Conn2Conn) Trans(ctx context.Context, lconn net.Conn, rconn net.Conn) {
	if lconn == nil {
		logger.Error(ctx, "左边连接为nil")
		return
	}
	if rconn == nil {
		logger.Error(ctx, "右边连接为nil")
	}
	atomic.AddInt64(&activeConns, 1)
	defer atomic.AddInt64(&activeConns, -1)
	t := tunnel.GetTunnel()
	t.SetLeft(&meteredConn{Conn: lconn})
	t.SetRight(&meteredConn{Conn: rconn})
	logger.Info(ctx, "Tunnel 开始流通")
	id := ctxtransid.GetID(ctx)
	t.Flow()
	t.Close()
	tunnel.PutTunnel(t)
	logger.Info(context.Background(), "id:", id, "Tunnel 停止流通")
}
