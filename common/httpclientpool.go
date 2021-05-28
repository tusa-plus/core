package common

import (
	"net/http"
	"sync"
)

type HttpClientPool struct {
	pool sync.Pool
}

func NewHttpClientPool() HttpClientPool {
	return HttpClientPool{
		pool: sync.Pool{
			New: func() interface{} {
				return &http.Client{}
			},
		},
	}
}

func (httpClientPool *HttpClientPool) Get() *http.Client {
	return httpClientPool.pool.Get().(*http.Client)
}

func (httpClientPool *HttpClientPool) Put(client *http.Client) {
	httpClientPool.pool.Put(client)
}
