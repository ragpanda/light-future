package light_future

import "context"

type Future struct{}
type Pool interface {
	Exec(func())
}

func New(ctx context.Context) *Future {
	return nil
}

func (self *Future) Use(pool Pool) *Future {
	return nil
}

func (self *Future) Send() *Future {
	return nil
}

func (self *Future) Await() *FutureResult {
	return nil
}

type FutureResult struct {
}

func (self *FutureResult) Result() (interface{}, error) {
	return nil, nil
}

func (self *FutureResult) ResultWithFill(r interface{}) error {
	return nil
}

func (self *FutureResult) Error() error {
	return nil
}

func Merge(fList ...*Future) *Future {
	return nil
}
