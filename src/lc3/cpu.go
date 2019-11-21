package lc3

import (
	"bufio"
	"fmt"
	"os"
)

const (
	R_R0 = iota
	R_R1
	R_R2
	R_R3
	R_R4
	R_R5
	R_R6
	R_R7
	R_PC
	R_COND
)

const (
	REG_MODE = iota
	IM_MODE
)

const (
	FL_POS = 1 << 0 /* P */
	FL_ZRO = 1 << 1 /* Z */
	FL_NEG = 1 << 2 /* N */

)

const (
	OP_BR   = iota /* branch */
	OP_ADD         /* add  */
	OP_LD          /* load */
	OP_ST          /* store */
	OP_JSR         /* jump register */
	OP_AND         /* bitwise and */
	OP_LDR         /* load register */
	OP_STR         /* store register */
	OP_RTI         /* unused */
	OP_NOT         /* bitwise not */
	OP_LDI         /* load indirect */
	OP_STI         /* store indirect */
	OP_JMP         /* jump */
	OP_RES         /* reserved (unused) */
	OP_LEA         /* load effective address */
	OP_TRAP        /* execute trap */

)

const (
	TRAP_GETC  = 0x20
	TRAP_OUT   = 0x21
	TRAP_PUTS  = 0x22
	TRAP_IN    = 0x23
	TRAP_PUTSP = 0x24
	TRAP_HALT  = 0x25
)

const (
	MR_KBSR = 0xFE00
	MR_KBDR = 0xFE02
)

const (
	memories, registers int = 65535, 10
)

type CPU struct {
	MEMORY [memories]uint16
	REGS   [registers]uint16
}

var cpu *CPU = &CPU{}

var running int = 1

func read_char() byte {
	reader := bufio.NewReader(os.Stdin)
	char, err := reader.ReadByte()

	if err != nil {
		fmt.Println(err)
		running = 0
		return 0
	}

	return char
}

func read_memory(address uint16) uint16 {
	if address == MR_KBSR {

		cpu.MEMORY[MR_KBSR] = 1 << 15
		cpu.MEMORY[MR_KBDR] = uint16(read_char())

	}

	return cpu.MEMORY[address]
}

func sign_extend(n uint16, bits int) uint16 {

	if (n>>(bits-1))&1 == 1 {

		return n | (65535 << bits)
	}
	return n
}

func update_flags(register int) {
	if cpu.REGS[register] == 0 {
		cpu.REGS[R_COND] = FL_ZRO
	} else if (cpu.REGS[register]>>15)&1 == 1 {
		cpu.REGS[R_COND] = FL_NEG
	} else {
		cpu.REGS[R_COND] = FL_POS
	}
}

func print_registry() {
	for i := 0; i < 10; i++ {
		fmt.Println(i, " --> ", cpu.REGS[i])
	}
}

func Run(memory [memories]uint16) {

	//m[0x3002] = 1000
	cpu.MEMORY = memory
	//fmt.Println(0x1261)
	//cpu.REGS[R_PC] = 0x3000

	for running == 1 {
		//print_registry()

		instruction := read_memory(cpu.REGS[R_PC])
		cpu.REGS[R_PC] = cpu.REGS[R_PC] + 1
		op := instruction >> 12

		switch op {
		case OP_BR:
			n := (instruction >> 11) & 1
			z := (instruction >> 10) & 1
			p := (instruction >> 9) & 1

			if (n&FL_NEG) == 1 || (z&FL_ZRO) == 1 || (p&FL_POS) == 1 {
				offset := instruction & 511
				cpu.REGS[R_PC] = cpu.REGS[R_PC] + sign_extend(offset, 9)
			}
		case OP_ADD:
			//fmt.Println("ADD")
			dr := (instruction >> 9) & 7
			sr1 := (instruction >> 6) & 7

			mode := (instruction >> 5) & 1

			if mode == 0 {
				sr2 := instruction & 7

				cpu.REGS[dr] = cpu.REGS[sr1] + cpu.REGS[sr2]

			} else {
				imm5 := instruction & 31

				cpu.REGS[dr] = cpu.REGS[sr1] + sign_extend(uint16(imm5), 5)
			}

			update_flags(int(dr))

		case OP_LD:

			//fmt.Println("LD")

			dr := (instruction >> 9) & 8

			offset := sign_extend(instruction&511, 9)

			cpu.REGS[dr] = read_memory(offset + cpu.REGS[R_PC])

			update_flags(int(dr))

		case OP_ST:
			//fmt.Println("ST")

			sr := (instruction >> 9) & 7

			offset := sign_extend(instruction&511, 9)

			cpu.MEMORY[offset+cpu.REGS[R_PC]] = cpu.REGS[sr]

		case OP_JSR:

			//fmt.Println("JSR")

			cpu.REGS[R_R7] = cpu.REGS[R_PC]

			mode := (instruction >> 11) & 1
			if mode == 0 {
				base := (instruction >> 6) & 7

				cpu.REGS[R_PC] = cpu.REGS[base]

			} else {
				offset := instruction & 2047

				cpu.REGS[R_PC] = cpu.REGS[R_PC] + sign_extend(offset, 11)

			}

		case OP_AND:
			//fmt.Println("AND")

			dr := (instruction >> 9) & 7
			sr1 := (instruction >> 6) & 7

			mode := (instruction >> 5) & 1

			if mode == 0 {
				sr2 := instruction & 7

				cpu.REGS[dr] = cpu.REGS[sr1] & cpu.REGS[sr2]

			} else {
				imm5 := instruction & 31
				cpu.REGS[dr] = cpu.REGS[sr1] & sign_extend(imm5, 31)

			}

			update_flags(int(dr))

		case OP_LDR:

			//fmt.Println("LDR")

			dr := (instruction >> 9) & 7

			base := (instruction >> 6) & 7
			offset := instruction & 63

			cpu.REGS[dr] = read_memory(base + sign_extend(offset, 6))

			update_flags(int(dr))

		case OP_STR:
			//fmt.Println("STR")

			sr := (instruction >> 9) & 7

			base := (instruction >> 6) & 7
			offset := instruction & 63

			cpu.MEMORY[base+sign_extend(offset, 6)] = cpu.REGS[sr]

		case OP_RTI:
			//fmt.Println("RTI")

		case OP_NOT:
			//fmt.Println("NOT")
			dr := (instruction >> 9) & 7
			sr := (instruction >> 6) & 7

			cpu.REGS[dr] = ^cpu.REGS[sr]

			update_flags(int(dr))

		case OP_LDI:
			//fmt.Println("LDI")
			dr := (instruction >> 9) & 7

			offset := instruction & 511

			cpu.REGS[dr] = read_memory(read_memory(cpu.REGS[R_PC] + sign_extend(offset, 9)))

			update_flags(int(dr))

		case OP_STI:
			//fmt.Println("STI")

			sr := (instruction >> 9) & 7

			offset := instruction & 511

			cpu.MEMORY[cpu.MEMORY[cpu.REGS[R_PC]+sign_extend(offset, 9)]] = cpu.REGS[sr]

		case OP_JMP:

			//fmt.Println("JMP")

			br := (instruction >> 6) & 7

			cpu.REGS[R_PC] = cpu.REGS[br]

		case OP_RES:

			//fmt.Println("RES")

		case OP_LEA:

			//fmt.Println("LEA")

			dr := (instruction >> 9) & 7

			offset := instruction & 511

			cpu.REGS[dr] = cpu.REGS[R_PC] + sign_extend(offset, 9)

			update_flags(int(dr))

		case OP_TRAP:

			//fmt.Println("TRAP")

			trapvect := instruction & 255

			switch trapvect {
			case TRAP_GETC:

				char := read_char()

				cpu.REGS[R_R0] = uint16(char)

			case TRAP_OUT:

				fmt.Printf("%c", byte(cpu.REGS[R_R0]))

			case TRAP_PUTS:
				start_address := cpu.REGS[R_R0]

				for i := 0; read_memory(start_address+uint16(i)) != 0x0; i++ {
					fmt.Printf("%c", byte(read_memory(start_address+uint16(i))))
				}

			case TRAP_IN:

				fmt.Printf("Enter a character: ")

				char := read_char()

				fmt.Printf("%c", char)
				cpu.REGS[R_R0] = uint16(char)

			case TRAP_PUTSP:

				start_address := cpu.REGS[R_R0]

				for i := 0; read_memory(start_address+uint16(i)) != 0x0; i++ {
					fmt.Printf("%c", byte(read_memory(start_address+uint16(i))&255))
					fmt.Printf("%c", byte((read_memory(start_address+uint16(i))>>8)&255))
				}

			case TRAP_HALT:
				fmt.Printf("HALT")

			default:
				fmt.Println("Unexpected Trap Code")
				running = 0
			}
			// cpu.REGS[R_R7] = cpu.REGS[R_PC]

			// cpu.REGS[R_PC] = cpu.MEMORY[sign_extend(trapvect, 8)]
		default:
			fmt.Println("Unexpected OpCode")
			running = 0
		}
	}

	fmt.Println("Programming Exited")

}
