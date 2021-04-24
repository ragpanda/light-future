package light_future

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Tmp struct {
	A int
	B bool
	C string
}

func BenchmarkFuture(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ctx := context.Background()

		var fList []*Future
		for i := 0; i < 100; i++ {
			a := i
			f := NewFuture(ctx, func(ctx context.Context) (interface{}, error) {
				return &Tmp{
					A: a + 10,
					C: fmt.Sprintf("%d", a*100000),
				}, nil
			}).Send()
			fList = append(fList, f)
		}

		for k, f := range fList {
			r, _ := f.Await().Result()
			rt := r.(*Tmp)
			assert.Equal(b, k, rt.A-10)
		}

	}
}

func BenchmarkGoroutineUsingChannelReturn(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ctx := context.Background()

		var fList []chan *Tmp
		for i := 0; i < 100; i++ {
			chanInt := make(chan *Tmp, 1)
			go func(ctx context.Context, i int) {
				b := &Tmp{
					A: i + 10,
					C: fmt.Sprintf("%d", i),
				}
				chanInt <- b
			}(ctx, i)
			fList = append(fList, chanInt)

		}

		for k, f := range fList {
			r := <-f
			assert.Equal(b, k, r.A-10)
		}

	}
}
