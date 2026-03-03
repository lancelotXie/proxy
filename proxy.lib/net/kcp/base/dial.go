package base

import (
	"crypto/sha1"
	"net"

	"github.com/lancelotXie/proxy/proxy.lib/configuration"
	"github.com/pkg/errors"
	"github.com/xtaci/kcp-go"
	"golang.org/x/crypto/pbkdf2"
)

type kcpSetting struct {
	key         string
	salt        string
	dataShard   int
	parityShard int
	noDelay     int
	interval    int
	resend      int
	nc          int
	sndWnd      int
	rcvWnd      int
	mtu         int
	ackNoDelay  bool
	dscp        int
	readBuffer  int
	writeBuffer int
}

func defaultKCPSetting() kcpSetting {
	return kcpSetting{
		key:         "",
		salt:        "",
		dataShard:   10,
		parityShard: 3,
		noDelay:     1,
		interval:    10,
		resend:      2,
		nc:          1,
		sndWnd:      1024,
		rcvWnd:      1024,
		mtu:         1350,
		ackNoDelay:  true,
		dscp:        46,
		readBuffer:  16777217,
		writeBuffer: 16777217,
	}
}

func loadKCPSetting() (cfg kcpSetting) {
	cfg = defaultKCPSetting()
	cfg.key, _ = configuration.GetString("kcp.key")
	cfg.salt, _ = configuration.GetString("kcp.salt")
	if value, ok := configuration.GetInt("kcp.dataShard"); ok {
		cfg.dataShard = value
	}
	if value, ok := configuration.GetInt("kcp.parityShard"); ok {
		cfg.parityShard = value
	}
	if value, ok := configuration.GetInt("kcp.noDelay"); ok {
		cfg.noDelay = value
	}
	if value, ok := configuration.GetInt("kcp.interval"); ok {
		cfg.interval = value
	}
	if value, ok := configuration.GetInt("kcp.resend"); ok {
		cfg.resend = value
	}
	if value, ok := configuration.GetInt("kcp.nc"); ok {
		cfg.nc = value
	}
	if value, ok := configuration.GetInt("kcp.sndWnd"); ok {
		cfg.sndWnd = value
	}
	if value, ok := configuration.GetInt("kcp.rcvWnd"); ok {
		cfg.rcvWnd = value
	}
	if value, ok := configuration.GetInt("kcp.mtu"); ok {
		cfg.mtu = value
	}
	if value, ok := configuration.GetInt("kcp.ackNoDelay"); ok {
		cfg.ackNoDelay = value != 0
	}
	if value, ok := configuration.GetInt("kcp.dscp"); ok {
		cfg.dscp = value
	}
	if value, ok := configuration.GetInt("kcp.readBuffer"); ok {
		cfg.readBuffer = value
	}
	if value, ok := configuration.GetInt("kcp.writeBuffer"); ok {
		cfg.writeBuffer = value
	}
	return cfg
}

func dial(addr string) (conn net.Conn, err error) {
	cfg := loadKCPSetting()
	var block kcp.BlockCrypt
	if cfg.key != "" && cfg.salt != "" {
		pass := pbkdf2.Key([]byte(cfg.key), []byte(cfg.salt), 4096, 32, sha1.New)
		block, _ = kcp.NewAESBlockCrypt(pass)
	}

	s, err := kcp.DialWithOptions(addr, block, cfg.dataShard, cfg.parityShard)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	defer func() {
		if err != nil {
			s.Close()
		}
	}()

	s.SetStreamMode(true)
	s.SetWriteDelay(false)
	s.SetNoDelay(cfg.noDelay, cfg.interval, cfg.resend, cfg.nc)
	s.SetWindowSize(cfg.sndWnd, cfg.rcvWnd)
	s.SetMtu(cfg.mtu)
	s.SetACKNoDelay(cfg.ackNoDelay)

	if err = s.SetDSCP(cfg.dscp); err != nil {
		return
	}
	if err = s.SetReadBuffer(cfg.readBuffer); err != nil {
		return
	}
	if err = s.SetWriteBuffer(cfg.writeBuffer); err != nil {
		return
	}

	conn = s
	return
}
