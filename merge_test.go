package light_future

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMergeNormal(t *testing.T) {
	ctx := context.Background()

	var fList []*Future
	for i := 0; i < 100; i++ {
		a := i
		f := NewFuture(ctx, NewClosureRunnable(func(ctx context.Context) (interface{}, error) {
			return a, nil
		})).Send()
		fList = append(fList, f)
	}

	r, err := Merge(ctx, fList...).Send().Await().Result()
	assert.Nil(t, err)
	assert.Len(t, r.([]interface{}), 100)
	for k, v := range r.([]interface{}) {
		assert.Nil(t, err)
		vInt := v.(int)
		assert.Equal(t, k, vInt)
	}

}

func TestMergeError(t *testing.T) {
	ctx := context.Background()

	var fList []*Future
	for i := 0; i < 100; i++ {
		a := i
		f := NewFuture(ctx, NewClosureRunnable(func(ctx context.Context) (interface{}, error) {
			panic("what?")
			return a, nil
		})).Send()
		fList = append(fList, f)
	}

	r, err := Merge(ctx, fList...).Send().Await().Result()
	assert.Len(t, r.([]interface{}), 0)
	assert.NotNil(t, err)
	assert.True(t, len(err.Error()) > 0)
	for _, v := range r.([]interface{}) {
		assert.Nil(t, v)
	}

}
