package ioutils

import (
	"fmt"
	"github.com/nsf/termbox-go"
)

func ClearLine(y int) {
	termbox.Sync()
	width, _ := termbox.Size()
	for i := 0; i < width; i++ {
		termbox.SetCell(i, y, ' ', termbox.ColorDefault, termbox.ColorDefault)
	}
	termbox.Flush()
}

func ConsoleOutLn(v string) {
	fmt.Println(v)
}

func ConsoleOut(v string, y int) {
	print(v, y)
}

func print(message string, y int) {
	ClearLine(y)
	runes := []rune(message)
	for x, r := range runes {
		termbox.SetCell(x, y, r, termbox.ColorDefault, termbox.ColorDefault)
	}
	termbox.Flush()
}
