package main

import (
	"fmt"

	"github.com/kadaan/tracerr"
)

func main() {
	if err := nilError(); err != nil {
		tracerr.PrintSourceColor(err)
	} else {
		fmt.Println("no error")
	}
}

func nilError() error {
	return tracerr.Wrap(nil)
}
