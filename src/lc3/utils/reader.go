package utils

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

func Load(memory *[^uint16(0)]uint16, path string) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return err
	}

	origin := binary.BigEndian.Uint16(b[:2])

	for i := 2; i < len(b); i += 2 {
		memory[origin] = binary.BigEndian.Uint16(b[i : i+2])
		origin++
	}

	return nil
}

// Does not work as expected
func ReadFile(path string) [65535]uint16 {

	var file [65535]uint16

	data, err := ioutil.ReadFile(path)

	if err != nil {
		fmt.Println(err)

	}

	var origin uint16 = swap16(uint16(data[0]))

	fmt.Println(origin)

	for i := 1; i < len(data); i++ {
		file[origin+uint16(i)] = swap16(uint16(data[i]))
		fmt.Println(swap16(uint16(data[i])))
	}

	return file
}

func swap16(value uint16) uint16 {
	return (value << 8) | (value >> 8)
}

func GetChar() (byte, error) {
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()

	b := make([]byte, 1)
	os.Stdin.Read(b)
	return b[0], nil
}

// IsKeyPressed checks if a key was pressed
func IsKeyPressed() bool {
	fi, _ := os.Stdin.Stat()
	return fi.Size() > 0
}
