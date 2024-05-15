package util

import "fmt"

const reset = "\033[0m"
const Black = "\033[30m"
const Red = "\033[31m"
const Green = "\033[32m"
const Yellow = "\033[33m"
const Blue = "\033[34m"
const Magenta = "\033[35m"
const Cyan = "\033[36m"
const White = "\033[37m"

func Fore(color string, ln string) string {
	return fmt.Sprint(color + ln + reset)
}
