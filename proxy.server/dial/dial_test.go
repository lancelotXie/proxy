package dial

import (
	"context"
	"testing"
	"time"

	"github.com/lancelotXie/proxy/proxy.lib/address"

	"github.com/stretchr/testify/assert"
)

var req = "GET / HTTP/1.1\r\n" +
	"Host: www.baidu.com\r\n" +
	"Connection: keep-alive\r\n" +
	"DNT: 1\r\n" +
	"Upgrade-Insecure-Requests: 1\r\n" +
	"User-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.108 Safari/537.36\r\n" +
	"Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3" +
	"Accept-Encoding: gzip, deflate\r\n" +
	"Accept-Language: en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7\r\n\r\n"

func TestConn2Web(t *testing.T) {
	var d Dial
	var webip = "110.242.74.102"
	var port int = 80
	d = &WebServer{}
	inaddr := address.New(webip, int(port), "")
	conn, err := d.Conn2Web(context.TODO(), inaddr)
	assert.Nil(t, err)
	assert.NotNil(t, conn)
	defer conn.Close()

	// 设置 2 秒超时，替代原来的 os.Exit(0)
	conn.SetDeadline(time.Now().Add(time.Second * 2))

	conn.Write([]byte(req))
	for {
		rbytes := make([]byte, 1024)
		n, err := conn.Read(rbytes)
		if err != nil {
			// 超时或其他错误时退出循环
			break
		}
		assert.NotZero(t, n)
	}
}
