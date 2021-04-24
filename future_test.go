package light_future

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/go-errors/errors"
	"github.com/stretchr/testify/assert"
)

func TestFutureNormal(t *testing.T) {
	ctx := context.Background()

	var fList []*Future
	for i := 0; i < 100; i++ {
		a := i
		f := NewFuture(ctx, NewClosureRunnable(func(ctx context.Context) (interface{}, error) {
			return a, nil
		})).Send()
		fList = append(fList, f)
	}

	assert.Len(t, fList, 100)
	for k, f := range fList {
		r, err := f.Await().Result()
		assert.Nil(t, err)
		rInt := r.(int)
		assert.Equal(t, k, rInt)
	}

}

func TestFutureNormalConcurrencyGetResult(t *testing.T) {

	ctx := context.Background()

	var fList []*Future
	for i := 0; i < 100; i++ {
		a := i
		f := NewFuture(ctx, NewClosureRunnable(func(ctx context.Context) (interface{}, error) {
			b := a
			return b, nil
		})).Use(&GoroutineInfanitePool{}).Send()
		fList = append(fList, f)
	}

	assert.Len(t, fList, 100)
	wg := sync.WaitGroup{}
	for k, f := range fList {
		for i := 0; i < 20; i++ {
			wg.Add(1)
			go func(k int, f *Future) {
				defer wg.Done()
				r, err := f.Await().Result()
				assert.Nil(t, err)
				rInt := r.(int)
				assert.Equal(t, k, rInt)
			}(k, f)
		}
	}
	wg.Wait()

}

func TestFutureNormalGetResultWithFill(t *testing.T) {
	ctx := context.Background()

	type Tmp struct {
		A int
	}

	var fList []*Future
	for i := 0; i < 100; i++ {
		a := i
		f := NewFuture(ctx, NewClosureRunnable(func(ctx context.Context) (interface{}, error) {

			return &Tmp{
				A: a,
			}, nil
		})).Send()
		fList = append(fList, f)
	}

	assert.Len(t, fList, 100)
	for k, f := range fList {
		tmp := &Tmp{}
		err := f.Await().ResultWithFill(tmp)
		assert.Nil(t, err)
		assert.Equal(t, k, tmp.A)
	}
}

func TestFutureNormalCancel(t *testing.T) {
	ctx := context.Background()

	type Tmp struct {
		A int
	}

	var fList []*Future
	for i := 0; i < 100; i++ {
		a := i
		f := NewFuture(ctx, NewClosureRunnable(func(ctx context.Context) (interface{}, error) {
			var channel chan int
			select {
			case <-channel:
				t.Logf("channel done")
			case <-ctx.Done():
				t.Logf("context done")
				return nil, errors.Errorf("ctx cancelled")
			}
			return &Tmp{
				A: a,
			}, nil
		})).Send().Cancel()
		fList = append(fList, f)
	}

	assert.Len(t, fList, 100)
	for _, f := range fList {
		tmp := &Tmp{}
		err := f.Await().ResultWithFill(tmp)
		assert.NotNil(t, err)
		assert.Equal(t, "ctx cancelled", err.Error())
	}
}

func TestFutureNormalTimeout(t *testing.T) {
	ctx := context.Background()

	type Tmp struct {
		A int
	}

	var fList []*Future
	for i := 0; i < 100; i++ {
		a := i
		f := NewFuture(ctx, NewClosureRunnable(func(ctx context.Context) (interface{}, error) {
			var channel chan int
			select {
			case <-channel:
				t.Logf("channel done")
			case <-ctx.Done():
				t.Logf("context done")
				return nil, errors.Errorf("ctx cancelled")
			}
			return &Tmp{
				A: a,
			}, nil
		})).Timeout(1 * time.Second).Send()
		fList = append(fList, f)
	}

	assert.Len(t, fList, 100)
	for _, f := range fList {
		tmp := &Tmp{}
		err := f.Await().ResultWithFill(tmp)
		assert.NotNil(t, err)
		assert.Equal(t, "ctx cancelled", err.Error())
	}
}

func TestFutureReturnError(t *testing.T) {
	ctx := context.Background()

	var fList []*Future
	for i := 0; i < 100; i++ {
		f := NewFuture(ctx, NewClosureRunnable(func(ctx context.Context) (interface{}, error) {
			return nil, errors.Errorf("dead")
		})).Send()
		fList = append(fList, f)
	}

	assert.Len(t, fList, 100)
	for _, f := range fList {
		r, err := f.Await().Result()
		assert.NotNil(t, err)
		assert.Nil(t, r)
		assert.Equal(t, "dead", err.Error())
	}

}

func TestFuturePanic(t *testing.T) {
	ctx := context.Background()

	var fList []*Future
	for i := 0; i < 100; i++ {
		f := NewFuture(ctx, NewClosureRunnable(func(ctx context.Context) (interface{}, error) {
			panic("dead")
			return 1, nil
		})).Send()
		fList = append(fList, f)
	}

	assert.Len(t, fList, 100)
	for _, f := range fList {
		r, err := f.Await().Result()
		assert.NotNil(t, err)
		assert.Nil(t, r)
		assert.NotNil(t, f.Await().Error())
	}

}

func TestFutureSendBefore(t *testing.T) {
	ctx := context.Background()
	assert.Panics(t, func() {
		NewFuture(ctx, NewClosureRunnable(func(ctx context.Context) (interface{}, error) {
			panic("dead")
			return 1, nil
		})).Await()
	}, "Await before send")
	assert.Panics(t, func() {
		NewFuture(ctx, NewClosureRunnable(func(ctx context.Context) (interface{}, error) {
			panic("dead")
			return 1, nil
		})).Cancel()
	}, "Cancel before send")

}
