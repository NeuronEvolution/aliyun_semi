package schedule

import "fmt"

var DebugEnabled = false

func SetDebug(v bool) {
	DebugEnabled = v
}

func debugLog(format string, a ...interface{}) {
	if DebugEnabled {
		fmt.Printf(format+"\n", a...)
	}
}
