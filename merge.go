package light_future

import "context"

func Merge(ctx context.Context, fList ...*Future) *Future {
	return NewFuture(ctx, func(ctx context.Context) (result interface{}, err error) {
		r := make([]interface{}, 0, len(fList))
		for _, f := range fList {
			f.Send()
		}

		mulError := &MultipleFutureError{}
		for _, f := range fList {
			itemResult := f.Await()
			if itemResult.status == statusSuccess {
				r = append(r, itemResult.result)
				mulError.successCount += 1
			} else if itemResult.status == statusError {
				mulError.errList = append(mulError.errList, itemResult.err)
				mulError.errorCount += 1
			}
		}
		if mulError.successCount < len(fList) {
			err = mulError
		}

		result = r
		return

	})

}
