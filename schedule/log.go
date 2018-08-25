package schedule

import (
	"fmt"
	"time"
)

func (r *ResourceManagement) log(format string, a ...interface{}) {
	fmt.Printf("["+r.Dataset+"]["+time.Now().Format(time.RFC3339)+"]"+format, a...)
}
