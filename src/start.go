package main

import (
	lc3 "./lc3"
	utils "./lc3/utils"
)

func main() {
	var memory [65535]uint16

	utils.Load(&memory, "C:\\Users\\harsh\\OneDrive\\Desktop\\c-test\\src\\2048.obj")
	lc3.Run(memory)
}
