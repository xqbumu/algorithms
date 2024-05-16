package main

import (
	"errors"
	"sync"
	"time"

	"github.com/phuslu/log"
)

var fcgiPools = map[string]*FcgiPool{} // fullAddress => *Pool
var fcgiPoolsLocker = sync.Mutex{}

type FcgiPool struct {
	size    uint16
	timeout time.Duration
	clients []*FcgiClient
	locker  sync.Mutex
}

func FcgiSharedPool(network string, address string, size uint16) *FcgiPool {
	fcgiPoolsLocker.Lock()
	defer fcgiPoolsLocker.Unlock()

	fullAddress := network + "//" + address
	pool, found := fcgiPools[fullAddress]
	if found {
		return pool
	}

	if size == 0 {
		size = 8
	}

	pool = &FcgiPool{
		size: size,
	}

	for i := uint16(0); i < size; i++ {
		client := NewFcgiClient(network, address)
		client.id = i
		client.KeepAlive()

		// prepare one for first request, and left for async request
		if i == 0 {
			if err := client.Connect(); err != nil {
				log.Error().Err(err).Msg("fcgi shared pool connect")
			}
		} else {
			go func() {
				if err := client.Connect(); err != nil {
					log.Error().Err(err).Msg("fcgi shared pool connect")
				}
			}()
		}
		pool.clients = append(pool.clients, client)
	}

	// watch connections
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		for range ticker.C {
			for _, client := range pool.clients {
				if client.isAvailable {
					continue
				}
				_ = client.Connect()
			}
		}
	}()

	fcgiPools[fullAddress] = pool

	return pool
}

func (p *FcgiPool) Client() (*FcgiClient, error) {
	p.locker.Lock()
	defer p.locker.Unlock()

	if len(p.clients) == 0 {
		return nil, errors.New("no available clients to use")
	}

	// find a free one
	for _, client := range p.clients {
		if client.isAvailable && client.isFree {
			return client, nil
		}
	}

	// find available on
	for _, client := range p.clients {
		if client.isAvailable {
			return client, nil
		}
	}

	// use first one
	if err := p.clients[0].Connect(); err == nil {
		return p.clients[0], nil
	}

	return nil, errors.New("no available clients to use")
}
