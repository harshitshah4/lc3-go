package main

import (
	"os"

	lc3 "./lc3"
	utils "./lc3/utils"
)

func main() {
	var memory [65535]uint16

	args := os.Args[1:]
	utils.Load(&memory, args[0])
	lc3.Run(memory)
}
