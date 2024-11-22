package xio

import "sync"

var mutex sync.RWMutex

var bufferPools = map[int]*sync.Pool{
	512: {
		New: func() any {
			return make([]byte, 512)
		},
	},
	1024: {
		New: func() any {
			return make([]byte, 512)
		},
	},
	1520: {
		New: func() any {
			return make([]byte, 512)
		},
	},
}

func Get(size int) []byte {
	mutex.RLock()

	if pool, exist := bufferPools[size]; exist {
		mutex.RUnlock()
		return pool.Get().([]byte)
	}

	mutex.RUnlock()

	mutex.Lock()
	defer mutex.Unlock()

	pool, exist := bufferPools[size]
	if exist {
		return pool.Get().([]byte)
	}

	pool = &sync.Pool{
		New: func() any {
			return make([]byte, size)
		},
	}

	bufferPools[size] = pool

	return pool.Get().([]byte)
}

func Put(buf []byte) {
	var size = len(buf)

	mutex.RLock()
	if pool, exist := bufferPools[size]; exist {
		mutex.RUnlock()
		pool.Put(buf)
		return
	}

	mutex.RUnlock()

	mutex.Lock()
	defer mutex.Unlock()

	pool, exist := bufferPools[size]
	if exist {
		pool.Put(buf)
		return
	}

	pool = &sync.Pool{
		New: func() any {
			return make([]byte, size)
		},
	}

	bufferPools[size] = pool

	pool.Put(buf)
}
