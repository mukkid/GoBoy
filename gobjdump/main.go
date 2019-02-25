package main

import (
	"bytes"
	"fmt"
	"github.com/SrsBusiness/gobjdump"
	"github.com/pborman/getopt/v2"
	"io/ioutil"
	"os"
)

var raw, gb *bool

func main_c(argv []string) int {
	if len(argv) < 2 {
		fmt.Printf("Usage: %s <options> <binary>\n", argv[0])
		return 1
	}

	binData, err := ioutil.ReadFile(argv[1])
	fileLen := len(binData)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return 1
	}

	reader := bytes.NewReader(binData)

	if *raw {
		return gobjdump.DisassemblerLoop(reader, 0x0, uint32(fileLen))
	} else { /* GameBoy Rom file */
		if gobjdump.GBROMPreamble(reader) != 0 {
			return 1
		}
		return 0
	}
}

func main() {
	raw = getopt.BoolLong("raw", 'r', "Raw Z80 binary file")
	gb = getopt.BoolLong("gbrom", 0, "GameBoy ROM file")
	getopt.Parse()

	if *raw && *gb {
		fmt.Printf("Cannot specify both --raw and --gbrom\n")
		os.Exit(1)
	}

	/* Raw file by default */
	if !*raw && !*gb {
		*raw = true
	}

	os.Exit(main_c(append([]string{os.Args[0]}, getopt.Args()...)))
}
