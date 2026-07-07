package util

import "os"

func IsPipe() bool {
	fileInfo, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return fileInfo.Mode()&os.ModeNamedPipe != 0
}
