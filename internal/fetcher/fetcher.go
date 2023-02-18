package fetcher

import (
	"context"
	"errors"
	"sync"

	"github.com/jackc/puddle/v2"
)

var ErrMethodNotRegister = errors.New("method requested not registered")

type FetchMethod string

type Fetcher interface {
	Fetch(url string) (string, error)
}

type FetcherFactory interface {
	Create() (Fetcher, error)
	Destroy(Fetcher)
}

type FetcherPool struct {
	sync.RWMutex
	pools map[FetchMethod]*puddle.Pool[Fetcher]
}

func NewPool() *FetcherPool {
	return &FetcherPool{
		pools: make(map[FetchMethod]*puddle.Pool[Fetcher]),
	}
}

func (f *FetcherPool) Register(method FetchMethod, factory FetcherFactory, maxSize int) error {
	f.Lock()
	defer f.Unlock()

	cfg := &puddle.Config[Fetcher]{
		Constructor: func(ctx context.Context) (res Fetcher, err error) {
			return factory.Create()
		},
		Destructor: func(res Fetcher) {
			factory.Destroy(res)
		},
		MaxSize: int32(maxSize),
	}

	p, err := puddle.NewPool(cfg)
	if err != nil {
		return err
	}

	f.pools[method] = p

	return nil
}

func (f *FetcherPool) Get(ctx context.Context, method FetchMethod) (*puddle.Resource[Fetcher], error) {
	f.RLock()
	defer f.RUnlock()

	pool, ok := f.pools[method]
	if !ok {
		return nil, ErrMethodNotRegister
	}

	return pool.Acquire(ctx)
}

func (f *FetcherPool) Methods() []FetchMethod {
	f.RLock()
	defer f.RUnlock()

	methods := make([]FetchMethod, 0)
	for method := range f.pools {
		methods = append(methods, method)
	}

	return methods
}
