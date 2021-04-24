package light_future

import (
	"context"
	"reflect"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/go-errors/errors"
)

type Future struct {
	ctx context.Context

	pool     Pool
	runnable FutureFunc
	execOnce sync.Once
	timeout  *time.Duration
	cancel   func()

	done   chan struct{}
	result unsafe.Pointer
	err    error
}

type status int8

const (
	statusSuccess = iota
	statusError
)

type syncResult struct {
	status status
	result interface{}
	err    error
}

type FutureFunc func(context.Context) (interface{}, error)

func NewFuture(ctx context.Context, runnable FutureFunc) *Future {

	f := &Future{
		ctx:      ctx,
		pool:     &GoroutineInfanitePool{},
		runnable: runnable,
		execOnce: sync.Once{},
	}
	return f
}

const DefaultTimeout = 5 * time.Second

func (self *Future) Use(pool Pool) *Future {
	self.pool = pool
	return self
}

func (self *Future) Timeout(t time.Duration) *Future {
	self.timeout = &t
	return self
}

func (self *Future) Send() *Future {
	self.execOnce.Do(func() {
		var ctx context.Context
		if self.timeout != nil {
			ctx, self.cancel = context.WithTimeout(self.ctx, *self.timeout)
		} else {
			ctx, self.cancel = context.WithCancel(self.ctx)
		}

		self.done = make(chan struct{}, 1)
		self.pool.Exec(func() {
			defer func() {
				if err := recover(); err != nil {
					atomic.StorePointer(&self.result, unsafe.Pointer(&syncResult{
						status: statusError,
						result: nil,
						err:    errors.New(errors.Wrap(err, 2).ErrorStack()),
					}))
					self.done <- struct{}{}
				}

				close(self.done)
			}()

			result, err := self.runnable(ctx)
			atomic.StorePointer(&self.result, unsafe.Pointer(&syncResult{
				status: statusSuccess,
				result: result,
				err:    err,
			}))
			self.done <- struct{}{}

		})
	})

	return self
}

func (self *Future) Cancel() *Future {
	if self.done == nil {
		panic("Cancel before send")
	}
	if self.cancel != nil {
		self.cancel()
	}

	return self
}

func (self *Future) Await() *FutureResult {
	if self.done == nil {
		panic("Await before send")
	}

	<-self.done

	r := (*syncResult)(atomic.LoadPointer(&self.result))
	return &FutureResult{
		syncResult: r,
	}
}

type FutureResult struct {
	*syncResult
}

func (self *FutureResult) Result() (interface{}, error) {
	return self.result, self.err
}

// r is &struct
func (self *FutureResult) ResultWithFill(r interface{}) error {
	if self.result != nil {
		reflect.ValueOf(r).Elem().Set(reflect.ValueOf(self.result).Elem())
	}
	return self.err
}

func (self *FutureResult) Error() error {
	return self.err
}
