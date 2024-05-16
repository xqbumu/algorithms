package main

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
	"testing"
	"time"
)

const fcgiTestPoolSize = 5

var cwd string

func TestMain(m *testing.M) {
	var err error
	cwd, err = os.Getwd()
	if err != nil {
		panic(err)
	}

	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestFcgiClientGet(t *testing.T) {
	client := &FcgiClient{
		network: "unix",
		address: "/run/php/php8.2-fpm.sock",
	}
	err := client.Connect()
	if err != nil {
		t.Fatal("connect err:", err.Error())
	}

	req := NewFcgiRequest()
	params := newParams()
	params["REQUEST_METHOD"] = "GET"
	req.SetParams(params)

	resp, stderr, err := client.Call(req)
	if err != nil {
		t.Fatal("call error:", err.Error())
	}

	if len(stderr) > 0 {
		t.Fatal("stderr:", string(stderr))
	}

	t.Log("resp, status:", resp.StatusCode)
	t.Log("resp, status message:", resp.Status)
	t.Log("resp headers:", resp.Header)

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	t.Log("resp body:", string(data))
}

func TestFcgiClientGetAlive(t *testing.T) {
	client := &FcgiClient{
		network: "unix",
		address: "/run/php/php8.2-fpm.sock",
	}
	client.KeepAlive()
	err := client.Connect()
	if err != nil {
		t.Fatal("connect err:", err.Error())
	}

	for i := 0; i < 10; i++ {
		req := NewFcgiRequest()
		params := newParams()
		params["REQUEST_METHOD"] = "GET"
		req.SetParams(params)

		resp, _, err := client.Call(req)
		if err != nil {
			t.Fatal("do error:", err.Error())
		}

		t.Log("resp, status:", resp.StatusCode)
		t.Log("resp, status message:", resp.Status)
		t.Log("resp headers:", resp.Header)

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()
		t.Log("resp body:", string(data))

		time.Sleep(1 * time.Second)
	}
}

func TestFcgiClientPost(t *testing.T) {
	client := &FcgiClient{
		network: "unix",
		address: "/run/php/php8.2-fpm.sock",
	}
	err := client.Connect()
	if err != nil {
		t.Fatal("connect err:", err.Error())
	}

	req := NewFcgiRequest()
	params := newParams()
	params["REQUEST_METHOD"] = "POST"
	params["CONTENT_TYPE"] = "application/x-www-form-urlencoded"
	req.SetParams(newParams())

	r := bytes.NewReader([]byte("name12=value&hello=world&name13=value&hello4=world"))
	//req.SetParam("CONTENT_LENGTH", fmt.Sprintf("%d", r.Len()))
	req.SetBody(r, uint32(r.Len()))

	resp, _, err := client.Call(req)
	if err != nil {
		t.Fatal("do error:", err.Error())
	}

	t.Log("resp, status:", resp.StatusCode)
	t.Log("resp, status message:", resp.Status)
	t.Log("resp headers:", resp.Header)

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	t.Log("resp body:", string(data))
}

func TestFcgiClientPerformance(t *testing.T) {
	threads := 100
	countRequests := 500
	countSuccess := 0
	countFail := 0
	locker := sync.Mutex{}
	beforeTime := time.Now()
	wg := sync.WaitGroup{}
	wg.Add(threads)

	pool := FcgiSharedPool("unix", "/run/php/php8.2-fpm.sock", fcgiTestPoolSize)

	for i := 0; i < threads; i++ {
		go func(i int) {
			defer wg.Done()

			for j := 0; j < countRequests; j++ {
				client, err := pool.Client()
				if err != nil {
					t.Fatal("connect err:", err.Error())
				}

				req := NewFcgiRequest()
				req.SetTimeout(5 * time.Second)
				params := newParams()
				params["REQUEST_METHOD"] = "GET"
				req.SetParams(params)

				resp, _, err := client.Call(req)
				if err != nil {
					locker.Lock()
					countFail++
					locker.Unlock()
					continue
				}

				if resp.StatusCode == 200 {
					data, err := io.ReadAll(resp.Body)
					if err != nil || !strings.Contains(string(data), "Welcome") {
						locker.Lock()
						countFail++
						locker.Unlock()
					} else {
						locker.Lock()
						countSuccess++
						locker.Unlock()
					}
				} else {
					locker.Lock()
					countFail++
					locker.Unlock()
				}
				resp.Body.Close()
			}
		}(i)
	}

	wg.Wait()

	t.Log("success:", countSuccess, "fail:", countFail, "qps:", int(float64(countSuccess+countFail)/time.Since(beforeTime).Seconds()))
}

func BenchmarkFcgiClient_KeppAlive(b *testing.B) {
	pool := FcgiSharedPool("unix", "/run/php/php8.2-fpm.sock", fcgiTestPoolSize)
	params := newParams()
	params["REQUEST_METHOD"] = "GET"

	for i := 0; i < b.N; i++ {
		client, err := pool.Client()
		if err != nil {
			b.Fatal("connect err:", err.Error())
		}

		req := NewFcgiRequest()
		req.SetTimeout(5 * time.Second)
		req.SetParams(params)

		resp, _, err := client.Call(req)
		if err != nil {
			b.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			b.Fatal("resp status code error")
		}
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
}

func newParams() map[string]string {
	return map[string]string{
		"SCRIPT_FILENAME": path.Join(cwd, "./testdata/welcome.php"),
		"SERVER_SOFTWARE": "ferry/1.0.0",
		"REMOTE_ADDR":     "127.0.0.1",
		"QUERY_STRING":    "name=value&__ACTION__=/@wx",

		"SERVER_NAME":       "ferry.local",
		"SERVER_ADDR":       "127.0.0.1:80",
		"SERVER_PORT":       "80",
		"REQUEST_URI":       "/welcome.php?__ACTION__=/@wx",
		"DOCUMENT_ROOT":     path.Join(cwd, "./testdata"),
		"GATEWAY_INTERFACE": "CGI/1.1",
		"REDIRECT_STATUS":   "200",
		"HTTP_HOST":         "ferry.local",
	}
}
