package cli

import (
	"bufio"
	"fmt"
	"os"
)

func Prompt(text string) (string, error) {
	fmt.Printf("%s: ", text)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return input[:len(input)-1], nil
}
