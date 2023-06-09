package quickws

import (
	"bytes"
	"crypto/md5"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	testBinaryMessage64kb = bytes.Repeat([]byte("1"), 65535)
	testTextMessage64kb   = bytes.Repeat([]byte("中"), 65535/len("中"))
	testBinaryMessage10   = bytes.Repeat([]byte("1"), 10)
)

type testMessageHandler struct {
	DefCallback
	t       *testing.T
	need    []byte
	callbed int32
	server  bool
	count   int
	got     chan []byte
}

func (t *testMessageHandler) OnMessage(c *Conn, op Opcode, msg []byte) {
	need := append([]byte(nil), t.need...)
	atomic.StoreInt32(&t.callbed, 1)
	if t.count == 0 {
		return
	}
	t.count--

	message := "#client"
	if t.server {
		message = "#server"
	}
	if len(msg) < 30 {
		assert.Equal(t.t, msg, need, message)
	} else {

		md51 := md5.Sum(need)
		md52 := md5.Sum(msg)
		assert.Equal(t.t, md51, md52, message)
	}
	err := c.WriteMessage(op, msg)
	assert.NoError(t.t, err)
}

func (t *testMessageHandler) OnClose(c *Conn, err error) {
	message := "#client.OnClose"
	if t.server {
		message = "#server.OnClose"
	}

	assert.NoError(t.t, err, message)
}

func newServrEcho(t *testing.T, data []byte) *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := Upgrade(w, r,
			WithServerCallback(&testMessageHandler{t: t, need: data, server: true, count: -1}),
		)
		assert.NoError(t, err)
		if err != nil {
			return
		}

		c.ReadLoop()
	}))

	ts.URL = "ws" + strings.TrimPrefix(ts.URL, "http")
	return ts
}

func Test_ReadMessage10(t *testing.T) {
	ts := newServrEcho(t, testBinaryMessage10)
	client := &testMessageHandler{t: t, need: append([]byte(nil), testBinaryMessage10...), count: 1}
	c, err := Dial(ts.URL, WithClientCallback(client))
	assert.NoError(t, err)
	go c.ReadLoop()

	tmp := append([]byte(nil), testBinaryMessage10...)
	err = c.WriteMessage(Binary, tmp)
	time.Sleep(time.Second / 2)
	assert.NoError(t, err)
	assert.Equal(t, atomic.LoadInt32(&client.callbed), int32(1))
}

func Test_ReadMessage64k(t *testing.T) {
	ts := newServrEcho(t, testBinaryMessage64kb)
	client := &testMessageHandler{t: t, need: append([]byte(nil), testBinaryMessage64kb...), count: 1}
	c, err := Dial(ts.URL, WithClientCallback(client))
	assert.NoError(t, err)
	go c.ReadLoop()

	tmp := append([]byte(nil), testBinaryMessage64kb...)
	err = c.WriteMessage(Binary, tmp)
	time.Sleep(time.Second / 2)
	assert.NoError(t, err)
	assert.Equal(t, atomic.LoadInt32(&client.callbed), int32(1))
}

func Test_ReadMessage64k_Text(t *testing.T) {
	ts := newServrEcho(t, testTextMessage64kb)
	client := &testMessageHandler{t: t, need: append([]byte(nil), testTextMessage64kb...), count: 1}
	c, err := Dial(ts.URL, WithClientCallback(client))
	assert.NoError(t, err)
	go c.ReadLoop()

	tmp := append([]byte(nil), testTextMessage64kb...)
	err = c.WriteMessage(Text, tmp)
	assert.NoError(t, err)
	time.Sleep(time.Second / 2)

	assert.NoError(t, err)
	assert.Equal(t, atomic.LoadInt32(&client.callbed), int32(1))
}