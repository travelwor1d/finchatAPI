package logerr

import (
	"fmt"
	"os"
)

func LogError(a ...interface{}) {
	fmt.Fprint(os.Stderr, a...)
}

func LogErrorf(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
}
