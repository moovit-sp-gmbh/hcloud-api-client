package internal

import (
	"fmt"
)

var debug = false

func ToggleDebug(d bool) {
	debug = d
}

func Debug(msg string) {
	if debug {
		fmt.Printf("%s\n", msg)
	}
}
