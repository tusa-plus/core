package utils

import (
	"net/http"
	"sync"
)

type HTTPClientPool struct {
	pool sync.Pool
}

func NewHTTPClientPool() HTTPClientPool {
	return HTTPClientPool{
		pool: sync.Pool{
			New: func() interface{} {
				return new(http.Client)
			},
		},
	}
}

func (httpClientPool *HTTPClientPool) Get() *http.Client {
	return httpClientPool.pool.Get().(*http.Client)
}

func (httpClientPool *HTTPClientPool) Put(client *http.Client) {
	httpClientPool.pool.Put(client)
}
