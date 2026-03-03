package flow

import (
	"context"
	"io"
	"net"
	"sync"
	"sync/atomic"

	"github.com/lancelotXie/proxy/proxy.lib/ctxtransid"
	"github.com/lancelotXie/proxy/proxy.lib/logger"
)

var bufferPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, 32*1024)
	},
}

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

// 确保 meteredConn 实现了 io.ReaderFrom 接口，以便在 Linux 下利用 splice 零拷贝技术
var _ io.ReaderFrom = (*meteredConn)(nil)

func (m *meteredConn) Read(b []byte) (n int, err error) {
	n, err = m.Conn.Read(b)
	if n > 0 {
		atomic.AddUint64(&transitedBytes, uint64(n))
	}
	return
}

// ReadFrom 实现了 io.ReaderFrom 接口
// 这允许 io.Copy 在底层支持的情况下（如 Linux 上的 TCPConn）使用 splice 系统调用，
// 从而实现零拷贝数据传输，大幅降低 CPU 占用。
func (m *meteredConn) ReadFrom(r io.Reader) (n int64, err error) {
	// 尝试解包 src，获取底层的 net.Conn，以便底层 ReadFrom 能识别出 TCPConn
	var src io.Reader = r
	if mc, ok := r.(*meteredConn); ok {
		src = mc.Conn
	}

	// 如果底层连接支持 ReadFrom (例如 *net.TCPConn)，则直接调用
	// 在 Linux 上，这会尝试使用 splice
	if rf, ok := m.Conn.(io.ReaderFrom); ok {
		n, err = rf.ReadFrom(src)
	} else {
		// 如果不支持 ReadFrom，回退到使用缓冲区的复制
		// 此时我们需要手动管理 buffer，因为 io.CopyBuffer 传入的 buf 在 ReadFrom 优先路径下不会被使用
		buf := bufferPool.Get().([]byte)
		defer bufferPool.Put(buf)
		// 使用 unwrapped src，避免调用 meteredConn.Read 导致重复统计（虽然这里 Read 不会被调用，因为 ReadFrom 控制了流程，但为了清晰）
		// 实际上 io.CopyBuffer 会调用 src.Read 和 m.Conn.Write
		n, err = io.CopyBuffer(m.Conn, src, buf)
	}

	// 统一统计流量
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
		return
	}
	atomic.AddInt64(&activeConns, 1)
	defer atomic.AddInt64(&activeConns, -1)

	logger.Info(ctx, "Tunnel 开始流通")
	id := ctxtransid.GetID(ctx)

	// 使用 meteredConn 包装以统计流量
	l := &meteredConn{Conn: lconn}
	r := &meteredConn{Conn: rconn}

	// 确保连接关闭
	// 使用 Once 确保只关闭一次（虽然 Close 可以多次调用，但为了清晰）
	var closeOnce sync.Once
	closeConns := func() {
		closeOnce.Do(func() {
			lconn.Close()
			rconn.Close()
		})
	}
	defer closeConns()

	// 监听上下文取消
	stop := context.AfterFunc(ctx, closeConns)
	defer stop()

	// 启动一个 goroutine 处理 l <- r
	go func() {
		defer closeConns() // 确保出错或完成后关闭连接，触发另一方退出
		io.Copy(l, r)
	}()

	// 主 goroutine 处理 r <- l
	io.Copy(r, l)

	logger.Info(context.Background(), "id:", id, "Tunnel 停止流通")
}
