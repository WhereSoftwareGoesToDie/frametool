package main

import (
	"fmt"
	"os"
)

func Errorf(format string, v ...interface{}) {
	f := format + "\n"
	fmt.Fprintf(os.Stderr, f, v...)
}
