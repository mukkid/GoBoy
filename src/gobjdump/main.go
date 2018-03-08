package main

import (
    "bufio"
    "github.com/pborman/getopt/v2"
    "os"
    "fmt"
    "encoding/hex"
    "encoding/binary"
    "strings"
)

var raw, gb *bool

func main_c(argv []string) int {
    if len(argv) < 2 {
        fmt.Printf("Usage: %s <options> <binary>\n", argv[0])
        return 1
    }

	/* Open binary */
	file, err := os.Open(argv[1])
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return 1
	}

    if *raw {
        bufReader := bufio.NewReader(file)
        return disassemblerLoop(bufReader, 0x0)
    } else { /* GameBoy Rom file */
        if gbPreamble(file) != 0 {
            return 1
        }
        return 0
    }
}

func gbPreamble(file *os.File) int {
    /* 
     * Code entry point is at 0x0100-0x0103
     * It is almost always nop followed by jp
     */
    var addr uint32 = 0x0100
    file.Seek(int64(addr), 0)
    r := bufio.NewReader(file)
    var (
        instruction []uint8
        mnemonic []string
        err error
    )
    for instruction, mnemonic, err = decodeInstruction(r);
        len(instruction) != 0 && instruction[0] == 0x00; /* while nops */
        instruction, mnemonic, err = decodeInstruction(r) {
        /* Generate hex encoding of instruction */
        instructionHex := make([]uint8, hex.EncodedLen(len(instruction)))
        hex.Encode(instructionHex, instruction)

        if err != nil {
            fmt.Printf("%s\n", err.Error())
            fmt.Printf("0x%016x: %-12s\n", addr, instructionHex)
            return 1
        }

        /* format - addr: <instruction bytes> <instruction mnemonic> */
        operands := ""
        if len(mnemonic) > 1 {
            operands = strings.Join(mnemonic[1:], ", ")
        }
        fmt.Printf("0x%016x: %-12s %-6s %s\n", addr, instructionHex, mnemonic[0], operands)
        addr += uint32(len(instruction))
    }
    addr += uint32(len(instruction))

    var target uint16
    switch instruction[0] {
    case 0xc3: /* jp */
        /* compute the offset of the jp */
        target = binary.LittleEndian.Uint16(instruction[1:])
        _, err = r.Discard(int(target) - int(addr))
        if err != nil {
            return 1
        }
        fmt.Printf("target: %d\n", target)
        return disassemblerLoop(r, uint32(target))
    }
    return 1
}

func main() {
    raw = getopt.BoolLong("raw", 'r', "Raw Z80 binary file")
    gb = getopt.BoolLong("gb", 0, "GameBoy ROM file")
    getopt.Parse()

    if *raw && *gb {
        fmt.Printf("Cannot specify both --raw and --gb\n")
        os.Exit(1)
    }

    /* Raw file by default */
    if !*raw && !*gb {
        *raw = true
    }

    os.Exit(main_c(append([]string{os.Args[0]}, getopt.Args()...)))
}
