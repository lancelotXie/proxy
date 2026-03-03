package flow

import (
	"context"
	"io"
	"log"
	"net"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func sendwrite(leftconn net.Conn, rightconn net.Conn,
	leftwrite string, rightwrite string, lclone net.Conn, rclone net.Conn,
	waitg *sync.WaitGroup, t *testing.T) {
	waitg.Add(4)
	go func() {
		for num := 0; num < 4; num++ {
			cwrite := leftwrite + strconv.Itoa(num)
			wlen, err := leftconn.Write([]byte(cwrite))
			assert.NotZero(t, wlen)
			assert.Nil(t, err)
		}
		waitg.Done()
	}()
	go func() {
		for num := 0; num < 4; num++ {
			cwrite := rightwrite + strconv.Itoa(num)
			wlen, err := rightconn.Write([]byte(cwrite))
			assert.NotZero(t, wlen)
			assert.Nil(t, err)
		}
		waitg.Done()
	}()

	go func() {
		for num := 0; num < 4; num++ {
			cwrite := leftwrite + strconv.Itoa(num)
			crbytes := make([]byte, 1024)
			wlen, err := rightconn.Read(crbytes)
			if err == nil {
				assert.NotZero(t, wlen)
				assert.Equal(t, cwrite, string(crbytes[:wlen]))
			} else if err != io.EOF {
				log.Println("client read err:", err)
			}
		}
		waitg.Done()
	}()

	go func() {
		for num := 0; num < 4; num++ {
			cwrite := rightwrite + strconv.Itoa(num)
			crbytes := make([]byte, 1024)
			wlen, err := leftconn.Read(crbytes)
			if err == nil {
				assert.NotZero(t, wlen)
				assert.Equal(t, cwrite, string(crbytes[:wlen]))
			} else if err != io.EOF {
				log.Println("client read err:", err)
			}
		}
		waitg.Done()
	}()
	waitg.Wait()

	leftconn.Close()
	rightconn.Close()
	lclone.Close()
	rclone.Close()
}

func TestMockconn(t *testing.T) {
	m := &mockconn{}
	var leftconn net.Conn
	var rightconn net.Conn
	var leftwrite = "left write"
	var rightwrite = "right write"
	leftconn, rightconn = m.GetConn()
	waitg := &sync.WaitGroup{}
	waitg.Add(4)
	go func() {
		for num := 0; num < 4; num++ {
			cwrite := leftwrite + strconv.Itoa(num)
			wlen, err := leftconn.Write([]byte(cwrite))
			assert.NotZero(t, wlen)
			assert.Nil(t, err)
		}
		waitg.Done()
	}()
	go func() {
		for num := 0; num < 4; num++ {
			cwrite := rightwrite + strconv.Itoa(num)
			wlen, err := rightconn.Write([]byte(cwrite))
			assert.NotZero(t, wlen)
			assert.Nil(t, err)
		}
		waitg.Done()
	}()

	go func() {
		for num := 0; num < 4; num++ {
			cwrite := leftwrite + strconv.Itoa(num)
			crbytes := make([]byte, 1024)
			wlen, err := rightconn.Read(crbytes)
			if err == nil {
				assert.NotZero(t, wlen)
				assert.Equal(t, cwrite, string(crbytes[:wlen]))
			} else if err != io.EOF {
				log.Println("client read err:", err)
			}
		}
		waitg.Done()
	}()

	go func() {
		for num := 0; num < 4; num++ {
			cwrite := rightwrite + strconv.Itoa(num)
			crbytes := make([]byte, 1024)
			wlen, err := leftconn.Read(crbytes)
			if err == nil {
				assert.NotZero(t, wlen)
				assert.Equal(t, cwrite, string(crbytes[:wlen]))
			} else if err != io.EOF {
				log.Println("client read err:", err)
			}
		}
		waitg.Done()
	}()
	waitg.Wait()

	leftconn.Close()
	rightconn.Close()
}

func TestTrans_DataFlow(t *testing.T) {
	// Use net.Pipe to simulate connections
	lLocal, lRemote := net.Pipe()
	rLocal, rRemote := net.Pipe()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start Trans
	c := &Conn2Conn{}

	// We need to close the connections eventually to stop Trans
	go func() {
		c.Trans(ctx, lRemote, rRemote)
	}()

	// Send data from lLocal -> rLocal
	msg := []byte("hello from left")
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		_, err := lLocal.Write(msg)
		assert.NoError(t, err)
	}()

	buf := make([]byte, 1024)
	n, err := rLocal.Read(buf)
	assert.NoError(t, err)
	assert.Equal(t, string(msg), string(buf[:n]))

	wg.Wait()

	// Send data from rLocal -> lLocal
	msg2 := []byte("hello from right")
	wg.Add(1)

	go func() {
		defer wg.Done()
		_, err := rLocal.Write(msg2)
		assert.NoError(t, err)
	}()

	n, err = lLocal.Read(buf)
	assert.NoError(t, err)
	assert.Equal(t, string(msg2), string(buf[:n]))

	wg.Wait()

	// Clean up
	lLocal.Close()
	rLocal.Close()
}

func TestTrans_Cancel(t *testing.T) {
	_, lRemote := net.Pipe()
	_, rRemote := net.Pipe()

	ctx, cancel := context.WithCancel(context.Background())

	c := &Conn2Conn{}

	done := make(chan struct{})
	go func() {
		c.Trans(ctx, lRemote, rRemote)
		close(done)
	}()

	// Cancel context
	cancel()

	select {
	case <-done:
		// Success
	case <-time.After(1 * time.Second):
		t.Fatal("Trans did not exit after context cancellation")
	}
}
