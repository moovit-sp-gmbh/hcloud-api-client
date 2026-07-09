package cli

import "fmt"

func Prompt(text string) string {
	fmt.Printf("%s: ", text)
	var input string
	fmt.Scan(&input)

	return input
}
