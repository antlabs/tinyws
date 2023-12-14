// Copyright 2021-2024 antlabs. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package quickws

import (
	"net/url"
	"testing"
)

func Test_SecWebSocketAcceptVal(t *testing.T) {
	need := "s3pPLMBiTxaQ9kYGzzhZRbK+xOo="
	got := secWebSocketAcceptVal("dGhlIHNhbXBsZSBub25jZQ==")
	if got != need {
		t.Errorf("need %s, got %s", need, got)
	}
}

func Test_getHttpErrMsg(t *testing.T) {
	t.Run("test 1", func(t *testing.T) {
		err := getHttpErrMsg(111)
		if err == nil {
			t.Errorf("err should not be nil")
		}
	})

	t.Run("test 2", func(t *testing.T) {
		err := getHttpErrMsg(400)
		if err == nil {
			t.Errorf("err should not be nil")
		}
	})
}

type test_getHostName struct {
	data string
	need string
}

func Test_getHostName(t *testing.T) {
	t.Run("test 1", func(t *testing.T) {
		for _, d := range []test_getHostName{
			{
				data: "http://www.baidu.com",
				need: "www.baidu.com:80",
			},
			{
				data: "http://www.baidu.com:333",
				need: "www.baidu.com:333",
			},
			{
				data: "https://www.baidu.com",
				need: "www.baidu.com:443",
			},
		} {

			u, err := url.Parse(d.data)
			if err != nil {
				t.Errorf("err should be nil, got %s", err)
			}
			if getHostName(u) != d.need {
				t.Errorf("need %s, got %s", d.need, getHostName(u))
			}
		}
	})
}
