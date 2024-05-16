// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestHTTPDialCancel(t *testing.T) {
	withServer := func(tcpDelay, tlsDelay, httpDelay time.Duration, fn func(srv *httptest.Server)) (int, int, int) {
		var (
			tcpDialCount      int32
			tlsHandshakeCount int32
			httpHandleCount   int32
		)

		srv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			atomic.AddInt32(&httpHandleCount, 1)
			time.Sleep(httpDelay)
		}))

		srv.EnableHTTP2 = false

		srv.Config.ErrorLog = log.New(io.Discard, "", 0)

		// Pause for a moment during the handshake so we can see what happens when
		// we cancel the Context of a completed HTTP Request.
		srv.TLS = &tls.Config{}
		srv.TLS.GetConfigForClient = func(chi *tls.ClientHelloInfo) (*tls.Config, error) {
			atomic.AddInt32(&tlsHandshakeCount, 1)
			time.Sleep(tlsDelay)
			return nil, nil
		}

		srv.StartTLS()
		defer srv.Close()

		// Before making any requests, add a delay to the TCP Dialer so we can see
		// what happens when we cancel the Context of a completed HTTP Request.
		srv.Client().Transport.(*http.Transport).DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			atomic.AddInt32(&tcpDialCount, 1)
			time.Sleep(tcpDelay)
			return (&net.Dialer{}).DialContext(ctx, network, addr)
		}
		// Allow a large connection pool
		srv.Client().Transport.(*http.Transport).MaxIdleConnsPerHost = 100

		fn(srv)

		return int(atomic.LoadInt32(&tcpDialCount)), int(atomic.LoadInt32(&tlsHandshakeCount)), int(atomic.LoadInt32(&httpHandleCount))
	}

	doRequest := func(ctx context.Context, srv *httptest.Server, timeout time.Duration) error {
		if timeout > 0 {
			// BUG: canceling the context associated with an already-complete
			// HTTP request leads to an increase in TLS handshake count.
			sub, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()
			ctx = sub
		}

		req, err := http.NewRequestWithContext(ctx, "GET", srv.URL, nil)
		if err != nil {
			return fmt.Errorf("NewRequestWithContext: %w", err)
		}

		resp, err := srv.Client().Do(req)
		if err != nil {
			return fmt.Errorf("Do Request: %w", err)
		}
		defer resp.Body.Close()

		_, err = io.Copy(io.Discard, resp.Body)
		if err != nil {
			return fmt.Errorf("Discard Body: %w", err)
		}

		return nil
	}

	callWithDelays := func(t *testing.T, tcpDelay, tlsDelay, httpDelay time.Duration, delays []time.Duration, timeout time.Duration) (int, int, int, time.Duration) {
		var total time.Duration
		tcp, tls, http := withServer(tcpDelay, tlsDelay, httpDelay, func(srv *httptest.Server) {
			ctx := context.Background()
			var wg sync.WaitGroup
			for _, delay := range delays {
				time.Sleep(delay)
				wg.Add(1)
				go func() {
					defer wg.Done()
					start := time.Now()
					err := doRequest(ctx, srv, timeout)
					if err != nil {
						t.Errorf("HTTP request failed: %v", err)
					}
					atomic.AddInt64((*int64)(&total), int64(time.Now().Sub(start)))
				}()
			}
			wg.Wait()
		})
		return tcp, tls, http, total / time.Duration(len(delays))
	}

	varyCancel := func(poolSize int, tcpDelay, tlsDelay, httpDelay time.Duration, delays []time.Duration) func(t *testing.T) {
		return func(t *testing.T) {
			fn := func(timeout time.Duration) func(t *testing.T) {
				return func(t *testing.T) {
					tcp, tls, http, avg := callWithDelays(t, tcpDelay, tlsDelay, httpDelay, delays, timeout)
					if tcp > poolSize {
						t.Errorf("TCP handshake count; %d > %d", tcp, poolSize)
					}
					if tls > poolSize {
						t.Errorf("TLS handshake count; %d > %d", tls, poolSize)
					}
					if t.Failed() {
						t.Logf("timeout=%s tcp=%d tls=%d http=%d", timeout, tcp, tls, http)
					}
					t.Logf("average duration %s", avg)
				}
			}
			// No timeout, so no context.WithCancel / WithTimeout call
			t.Run("no timeout/cancel", fn(0))
			// Huge timeout, key change is the presence of "defer cancel()" once
			// the outbound request is complete
			t.Run("with timeout/cancel", fn(20*time.Minute))
		}
	}

	// Go's HTTP client connection pool has discarded useful progress on
	// outbound TCP handshakes for several releases.
	//
	// NOTE: this test doesn't work well on the Playground.
	t.Run("minimal pool with slow TCP", func(t *testing.T) {
		t.Logf("NOTE: expect failure here for all recent Go versions")
		varyCancel(
			2,                    // ideal pool size
			200*time.Millisecond, // delay in TCP
			0,                    // delay in TLS
			50*time.Millisecond,  // delay in HTTP
			[]time.Duration{
				0,                      // t=0 ms     create connection 1
				400 * time.Millisecond, // t=400 ms   use warm connection 1 until t=450
				20 * time.Millisecond,  // t=420 ms   trigger new connection 2 (BUG: work may be discarded!)
				380 * time.Millisecond, // t=800 ms   observe pool size (use 1)
				0,                      // t=800 ms   observe pool size (use 2, or dial 3)
			})(t)
	})

	// New in Go 1.17 via https://golang.org/issue/32406, Go's HTTP client
	// connection pool now also discards useful progress on outbound TLS
	// handshakes.
	//
	// At t=0ms, the first request triggers a new connection
	// At t=0ms, the TCP handshake for the first connection is complete and the TLS handshake begins
	// At t=200ms, the TLS handshake completes and the first HTTP request begins
	// At t=250ms, the first HTTP request completes and the first connection enters the idle pool
	// At t=400ms, the second request removes the first connection from the idle pool
	// At t=420ms, the third request finds an empty pool and dials a second connection
	// At t=420ms, the second connection TLS handshake begins
	// At t=450ms, the second HTTP request completes and hands its connection to the pool
	// At t=450ms, the third request intercepts the first connection before it enters the pool
	// At t=500ms, the third HTTP request completes and returns to the application code
	// At t=500ms, the application code has the full HTTP response, so cancels its Context
	// At t=500ms, Go 1.17's call to tls.Conn.HandshakeContext aborts
	// At t=620ms, Go 1.16's call to tls.Conn.Handshake completes, and goes into the idle pool
	// At t=800ms, the fourth request removes the first connection from the idle pool
	// At t=800ms, the fifth request uses an idle connection (Go 1.16) or dials fresh (Go 1.17+)
	t.Run("minimal pool with slow TLS", varyCancel(
		2,                    // ideal pool size
		0,                    // delay in TCP
		200*time.Millisecond, // delay in TLS
		50*time.Millisecond,  // delay in HTTP
		[]time.Duration{
			0,                      // t=0 ms     create connection 1
			400 * time.Millisecond, // t=400 ms   use warm connection 1 until t=450
			20 * time.Millisecond,  // t=420 ms   trigger new connection 2 (BUG: work may be discarded!)
			380 * time.Millisecond, // t=800 ms   observe pool size (use 1)
			0,                      // t=800 ms   observe pool size (use 2, or dial 3)
		}))

	// The impact of discarding useful progress on TLS handshakes is unbounded:
	// A client running Go 1.17 or newer, which creates a context for each
	// request which it cancels when the request is complete, may steadily churn
	// through new TLS connections. It can do this even when its maximum
	// outbound concurrency is below the MaxIdleConnsPerHost limit.
	t.Run("large pool with slow TLS", varyCancel(
		8,                    // ideal pool size
		0,                    // delay in TCP
		200*time.Millisecond, // delay in TLS
		50*time.Millisecond,  // delay in HTTP
		[]time.Duration{
			0,                      // t=0 ms     create connection 1
			0,                      // t=0 ms     create connection 2
			0,                      // t=0 ms     create connection 3
			0,                      // t=0 ms     create connection 4
			400 * time.Millisecond, // t=400 ms   use warm connection 1 until t=450
			0,                      // t=400 ms   use warm connection 2 until t=450
			0,                      // t=400 ms   use warm connection 3 until t=450
			0,                      // t=400 ms   use warm connection 4 until t=450
			20 * time.Millisecond,  // t=420 ms   trigger new connection 5 (BUG: work may be discarded!)
			0,                      // t=420 ms   trigger new connection 6 (BUG: work may be discarded!)
			0,                      // t=420 ms   trigger new connection 7 (BUG: work may be discarded!)
			0,                      // t=420 ms   trigger new connection 8 (BUG: work may be discarded!)
			380 * time.Millisecond, // t=800 ms   use warm connection 1 until t=850
			0,                      // t=800 ms   use warm connection 2 until t=850
			0,                      // t=800 ms   use warm connection 3 until t=850
			0,                      // t=800 ms   use warm connection 4 until t=850
			20 * time.Millisecond,  // t=820 ms   use warm connection 5, or trigger new connection 9 (BUG: work may be discarded!)
			0,                      // t=820 ms   use warm connection 6, or trigger new connection 10 (BUG: work may be discarded!)
			0,                      // t=820 ms   use warm connection 7, or trigger new connection 11 (BUG: work may be discarded!)
			0,                      // t=820 ms   use warm connection 8, or trigger new connection 12 (BUG: work may be discarded!)
			380 * time.Millisecond, // t=1200 ms  use warm connection 1 until t=1250
			0,                      // t=1200 ms  use warm connection 2 until t=1250
			0,                      // t=1200 ms  use warm connection 3 until t=1250
			0,                      // t=1200 ms  use warm connection 4 until t=1250
			20 * time.Millisecond,  // t=1220 ms  use warm connection 5, or trigger new connection 13 (BUG: work may be discarded!)
			0,                      // t=1220 ms  use warm connection 6, or trigger new connection 14 (BUG: work may be discarded!)
			0,                      // t=1220 ms  use warm connection 7, or trigger new connection 15 (BUG: work may be discarded!)
			0,                      // t=1220 ms  use warm connection 8, or trigger new connection 16 (BUG: work may be discarded!)
			380 * time.Millisecond, // t=1600 ms  use warm connection 1 until t=1650
			0,                      // t=1600 ms  use warm connection 2 until t=1650
			0,                      // t=1600 ms  use warm connection 3 until t=1650
			0,                      // t=1600 ms  use warm connection 4 until t=1650
			20 * time.Millisecond,  // t=1620 ms  use warm connection 5, or trigger new connection 17 (BUG: work may be discarded!)
			0,                      // t=1620 ms  use warm connection 6, or trigger new connection 18 (BUG: work may be discarded!)
			0,                      // t=1620 ms  use warm connection 7, or trigger new connection 19 (BUG: work may be discarded!)
			0,                      // t=1620 ms  use warm connection 8, or trigger new connection 20 (BUG: work may be discarded!)
			380 * time.Millisecond, // t=2020 ms  observe pool size (use 1)
			0,                      // t=2020 ms  observe pool size (use 2)
			0,                      // t=2020 ms  observe pool size (use 3)
			0,                      // t=2020 ms  observe pool size (use 4)
			0,                      // t=2020 ms  observe pool size (use 5, or dial 21)
			0,                      // t=2020 ms  observe pool size (use 6, or dial 22)
			0,                      // t=2020 ms  observe pool size (use 7, or dial 23)
			0,                      // t=2020 ms  observe pool size (use 8, or dial 24)
		}))
}
