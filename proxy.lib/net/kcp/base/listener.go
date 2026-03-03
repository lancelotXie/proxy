package base

import (
	"crypto/sha1"
	"net"

	"github.com/pkg/errors"
	"golang.org/x/crypto/pbkdf2"

	"github.com/xtaci/kcp-go"
)

type listener struct {
	base *kcp.Listener
}

func listen(addr string) (lis *listener, err error) {
	cfg := loadKCPSetting()
	var block kcp.BlockCrypt
	if cfg.key != "" && cfg.salt != "" {
		pass := pbkdf2.Key([]byte(cfg.key), []byte(cfg.salt), 4096, 32, sha1.New)
		block, _ = kcp.NewAESBlockCrypt(pass)
	}

	base, err := kcp.ListenWithOptions(addr, block, cfg.dataShard, cfg.parityShard)
	if err != nil {
		err = errors.WithStack(err)
	}
	defer func() {
		if err != nil {
			base.Close()
		}
	}()

	if err = base.SetDSCP(cfg.dscp); err != nil {
		return
	}
	if err = base.SetReadBuffer(cfg.readBuffer); err != nil {
		return
	}
	if err = base.SetWriteBuffer(cfg.writeBuffer); err != nil {
		return
	}

	lis = new(listener)
	lis.base = base
	return
}

// Accept 接收新的连接
func (l *listener) Accept() (conn net.Conn, err error) {
	cfg := loadKCPSetting()
	s, err := l.base.AcceptKCP()
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	s.SetStreamMode(true)
	s.SetWriteDelay(false)
	s.SetNoDelay(cfg.noDelay, cfg.interval, cfg.resend, cfg.nc)
	s.SetMtu(cfg.mtu)
	s.SetWindowSize(cfg.sndWnd, cfg.rcvWnd)
	s.SetACKNoDelay(cfg.ackNoDelay)

	conn = s
	return
}

func (l *listener) Close() (err error) {
	err = l.base.Close()
	return
}

func (l *listener) Addr() (addr net.Addr) {
	addr = l.base.Addr()
	return
}
