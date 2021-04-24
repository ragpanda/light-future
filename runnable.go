package light_future

import "context"

type Runnable interface {
	Run(context.Context) (interface{}, error)
}

// Params pass by closure
type ClosureRunnable struct {
	execFunc func(ctx context.Context) (result interface{}, err error)
}

func NewClosureRunnable(execFunc func(ctx context.Context) (result interface{}, err error)) *ClosureRunnable {
	return &ClosureRunnable{
		execFunc: execFunc,
	}
}
func (self *ClosureRunnable) Run(ctx context.Context) (interface{}, error) {
	return self.execFunc(ctx)
}
