// Copyright 2021-2023 antlabs. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package quickws

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// 测试客户端和服务端都有的配置项
func Test_CommonOption(t *testing.T) {
	t.Run("2.server.local: WithServerTCPDelay", func(t *testing.T) {
		run := int32(0)
		done := make(chan bool, 1)
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c, err := Upgrade(w, r, WithServerTCPDelay(), WithServerOnMessageFunc(func(c *Conn, mt Opcode, payload []byte) {
				atomic.AddInt32(&run, int32(1))
				done <- true
			}))
			if err != nil {
				t.Error(err)
			}
			c.StartReadLoop()
		}))

		defer ts.Close()

		url := strings.ReplaceAll(ts.URL, "http", "ws")
		con, err := Dial(url, WithClientCallback(&testDefaultCallback{}))
		if err != nil {
			t.Error(err)
		}
		defer con.Close()

		con.WriteMessage(Binary, []byte("hello"))
		select {
		case <-done:
		case <-time.After(1000 * time.Millisecond):
		}
		if atomic.LoadInt32(&run) != 1 {
			t.Error("not run server:method fail")
		}
	})

	t.Run("2.server.global: WithServerTCPDelay", func(t *testing.T) {
		run := int32(0)
		done := make(chan bool, 1)
		upgrade := NewUpgrade(WithServerTCPDelay(), WithServerOnMessageFunc(func(c *Conn, mt Opcode, payload []byte) {
			atomic.AddInt32(&run, int32(1))
			done <- true
		}))

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c, err := upgrade.Upgrade(w, r)
			if err != nil {
				t.Error(err)
			}
			c.StartReadLoop()
		}))

		defer ts.Close()

		url := strings.ReplaceAll(ts.URL, "http", "ws")
		con, err := Dial(url, WithClientCallback(&testDefaultCallback{}))
		if err != nil {
			t.Error(err)
		}

		con.WriteMessage(Binary, []byte("hello"))
		select {
		case <-done:
		case <-time.After(100 * time.Millisecond):
		}
		if atomic.LoadInt32(&run) != 1 {
			t.Error("not run server:method fail")
		}
	})

	t.Run("2.client: WithClientTCPDelay ", func(t *testing.T) {
		run := int32(0)
		done := make(chan bool, 1)
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c, err := Upgrade(w, r, WithServerOnMessageFunc(func(c *Conn, mt Opcode, payload []byte) {
				c.WriteMessage(mt, payload)
			}))
			if err != nil {
				t.Error(err)
			}
			c.StartReadLoop()
		}))

		defer ts.Close()

		url := strings.ReplaceAll(ts.URL, "http", "ws")
		con, err := Dial(url, WithClientTCPDelay(), WithClientOnMessageFunc(func(c *Conn, mt Opcode, payload []byte) {
			atomic.AddInt32(&run, int32(1))
			done <- true
		}))
		if err != nil {
			t.Error(err)
		}

		con.StartReadLoop()
		con.WriteMessage(Binary, []byte("hello"))
		select {
		case <-done:
		case <-time.After(100 * time.Millisecond):
		}
		if atomic.LoadInt32(&run) != 1 {
			t.Error("not run client callback:method fail")
		}
	})

	t.Run("4.server.local: WithServerOnMessageFunc", func(t *testing.T) {
		run := int32(0)
		done := make(chan bool, 1)
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c, err := Upgrade(w, r, WithServerOnMessageFunc(func(c *Conn, mt Opcode, payload []byte) {
				atomic.AddInt32(&run, int32(1))
				done <- true
			}))
			if err != nil {
				t.Error(err)
			}
			c.StartReadLoop()
		}))

		defer ts.Close()

		url := strings.ReplaceAll(ts.URL, "http", "ws")
		con, err := Dial(url, WithClientCallback(&testDefaultCallback{}))
		if err != nil {
			t.Error(err)
		}
		defer con.Close()

		con.WriteMessage(Binary, []byte("hello"))
		select {
		case <-done:
		case <-time.After(1000 * time.Millisecond):
		}
		if atomic.LoadInt32(&run) != 1 {
			t.Error("not run server:method fail")
		}
	})

	t.Run("4.server.global: WithServerOnMessageFunc", func(t *testing.T) {
		run := int32(0)
		done := make(chan bool, 1)
		upgrade := NewUpgrade(WithServerOnMessageFunc(func(c *Conn, mt Opcode, payload []byte) {
			atomic.AddInt32(&run, int32(1))
			done <- true
		}))

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c, err := upgrade.Upgrade(w, r)
			if err != nil {
				t.Error(err)
			}
			c.StartReadLoop()
		}))

		defer ts.Close()

		url := strings.ReplaceAll(ts.URL, "http", "ws")
		con, err := Dial(url, WithClientCallback(&testDefaultCallback{}))
		if err != nil {
			t.Error(err)
		}

		con.WriteMessage(Binary, []byte("hello"))
		select {
		case <-done:
		case <-time.After(100 * time.Millisecond):
		}
		if atomic.LoadInt32(&run) != 1 {
			t.Error("not run server:method fail")
		}
	})

	t.Run("4.client: WithClientOnMessageFunc", func(t *testing.T) {
		run := int32(0)
		done := make(chan bool, 1)
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c, err := Upgrade(w, r, WithServerOnMessageFunc(func(c *Conn, mt Opcode, payload []byte) {
				c.WriteMessage(mt, payload)
			}))
			if err != nil {
				t.Error(err)
			}
			c.StartReadLoop()
		}))

		defer ts.Close()

		url := strings.ReplaceAll(ts.URL, "http", "ws")
		con, err := Dial(url, WithClientOnMessageFunc(func(c *Conn, mt Opcode, payload []byte) {
			atomic.AddInt32(&run, int32(1))
			done <- true
		}))
		if err != nil {
			t.Error(err)
		}

		con.StartReadLoop()
		con.WriteMessage(Binary, []byte("hello"))
		select {
		case <-done:
		case <-time.After(100 * time.Millisecond):
		}
		if atomic.LoadInt32(&run) != 1 {
			t.Error("not run client callback:method fail")
		}
	})

	t.Run("9.server.local: WithServerBufioParseMode", func(t *testing.T) {
		run := int32(0)
		data := make(chan string, 1)
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c, err := Upgrade(w, r, WithServerBufioParseMode(), WithServerOnMessageFunc(func(c *Conn, op Opcode, payload []byte) {
				c.WriteMessage(op, payload)
			}))
			if err != nil {
				t.Error(err)
			}
			c.StartReadLoop()
		}))

		defer ts.Close()

		url := strings.ReplaceAll(ts.URL, "http", "ws")
		con, err := Dial(url, WithClientOnMessageFunc(func(c *Conn, mt Opcode, payload []byte) {
			atomic.AddInt32(&run, int32(1))
			data <- string(payload)
		}))
		if err != nil {
			t.Error(err)
		}
		defer con.Close()

		con.WriteMessage(Binary, []byte("hello"))
		con.StartReadLoop()
		select {
		case d := <-data:
			if d != "hello" {
				t.Errorf("write message or read message fail:got:%s, need:hello\n", d)
			}
		case <-time.After(1000 * time.Millisecond):
		}
		if atomic.LoadInt32(&run) != 1 {
			t.Error("not run server:method fail")
		}
	})

	t.Run("9.server.global: WithServerBufioParseMode", func(t *testing.T) {
		run := int32(0)
		data := make(chan string, 1)
		upgrade := NewUpgrade(WithServerBufioParseMode(), WithServerOnMessageFunc(func(c *Conn, op Opcode, payload []byte) {
			c.WriteMessage(op, payload)
		}))
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c, err := upgrade.Upgrade(w, r)
			if err != nil {
				t.Error(err)
			}
			c.StartReadLoop()
		}))

		defer ts.Close()

		url := strings.ReplaceAll(ts.URL, "http", "ws")
		con, err := Dial(url, WithClientOnMessageFunc(func(c *Conn, mt Opcode, payload []byte) {
			atomic.AddInt32(&run, int32(1))
			data <- string(payload)
		}))
		if err != nil {
			t.Error(err)
		}
		defer con.Close()

		con.WriteMessage(Binary, []byte("hello"))
		con.StartReadLoop()
		select {
		case d := <-data:
			if d != "hello" {
				t.Errorf("write message or read message fail:got:%s, need:hello\n", d)
			}
		case <-time.After(1000 * time.Millisecond):
		}
		if atomic.LoadInt32(&run) != 1 {
			t.Error("not run server:method fail")
		}
	})

	t.Run("9.client: WithClientOnMessageFunc", func(t *testing.T) {
		run := int32(0)
		data := make(chan string, 1)
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c, err := Upgrade(w, r, WithServerBufioParseMode(), WithServerOnMessageFunc(func(c *Conn, op Opcode, payload []byte) {
				c.WriteMessage(op, payload)
			}))
			if err != nil {
				t.Error(err)
			}
			c.StartReadLoop()
		}))

		defer ts.Close()

		url := strings.ReplaceAll(ts.URL, "http", "ws")
		con, err := Dial(url, WithClientBufioParseMode(), WithClientOnMessageFunc(func(c *Conn, mt Opcode, payload []byte) {
			atomic.AddInt32(&run, int32(1))
			data <- string(payload)
		}))
		if err != nil {
			t.Error(err)
		}
		defer con.Close()

		con.WriteMessage(Binary, []byte("hello"))
		con.StartReadLoop()
		select {
		case d := <-data:
			if d != "hello" {
				t.Errorf("write message or read message fail:got:%s, need:hello\n", d)
			}
		case <-time.After(1000 * time.Millisecond):
		}
		if atomic.LoadInt32(&run) != 1 {
			t.Error("not run server:method fail")
		}
	})
}
