package ioutils

import (
	"fmt"
)

func ClearLine() {
	fmt.Print("\033[2K")
	goToStart()
}

func ConsoleOutLn(v string) {
	fmt.Println(v)
}

func NewLine() {
	fmt.Println()
}

func ConsoleOut(v string) {
	print(v)
}

func goToStart() {
	fmt.Print("\033[999D")
}

func print(message string) {
	ClearLine()
	fmt.Print(message)
}
