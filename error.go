package light_future

import (
	"fmt"
	"strings"
)

type FutureErrorAbstract interface {
	Error() string
}

type MultipleFutureError struct {
	successCount int
	errorCount   int
	errList      []error
}

func (self *MultipleFutureError) Error() string {
	errStrList := make([]string, 0, len(self.errList))
	for _, v := range self.errList {
		errStrList = append(errStrList, v.Error())
	}
	errStr := strings.Join(errStrList, "\n")

	return fmt.Sprintf("[MultipleFutureError] err:\n%s", errStr)
}
