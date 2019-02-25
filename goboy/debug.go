package main

import "bytes"
import "fmt"
import "strings"
import "os"
import "bufio"
import "strconv"
import . "github.com/SrsBusiness/gobjdump"
import "regexp"

type FunctionFrame struct {
	addr       uint16 /* Address of frame on stack */
	returnAddr uint16
}

type Debugger struct {
	gb          *GameBoy
	breakpoints map[uint16]struct{} /* This is how sets work */
	paused      bool
	ROMReader   *bytes.Reader
}

func isBreakpoint(m map[uint16]struct{}, breakpoint uint16) bool {
	_, isMember := m[breakpoint]
	return isMember
}

func (d *Debugger) cont() {
	d.paused = false
	d.next()
	for !isBreakpoint(d.breakpoints, d.gb.get16Reg(PC)) && !d.paused {
		d.next()
	}
}

func (d *Debugger) addBreakpoint(addr uint16) {
	d.breakpoints[addr] = struct{}{}
}

func (d *Debugger) deleteBreakpoint(addr uint16) {
	delete(d.breakpoints, addr)
}

func (d *Debugger) next() {
	/* Probably don't care about d.paused here */
	d.gb.handleInterrupt()
	d.gb.Step()
}

func (d *Debugger) run() {
	/* TODO: reinitialize to clean state i.e. clear registers, reload ROM, reset memory */
	d.cont()
}

var regNames8 map[string]Reg8ID = map[string]Reg8ID{
	"a": A,
	"b": B,
	"c": C,
	"d": D,
	"e": E,
	"f": F,
	"h": H,
	"l": L,
}
var regNames16 map[string]Reg16ID = map[string]Reg16ID{
	"bc": BC,
	"de": DE,
	"hl": HL,
	"sp": SP,
	"pc": PC,
	"af": AF,
}

func (d *Debugger) printAllRegs() {
	fmt.Printf(
		`A:  0x%04x
B:  0x%04x
C:  0x%04x
D:  0x%04x
E:  0x%04x
F:  0x%04x
H:  0x%04x
L:  0x%04x
BC: 0x%04x
DE: 0x%04x
HL: 0x%04x
SP: 0x%04x
PC: 0x%04x
AF: 0x%04x
`,
		d.gb.get8Reg(A),
		d.gb.get8Reg(B),
		d.gb.get8Reg(C),
		d.gb.get8Reg(D),
		d.gb.get8Reg(E),
		d.gb.get8Reg(F),
		d.gb.get8Reg(H),
		d.gb.get8Reg(L),
		d.gb.get16Reg(BC),
		d.gb.get16Reg(DE),
		d.gb.get16Reg(HL),
		d.gb.get16Reg(SP),
		d.gb.get16Reg(PC),
		d.gb.get16Reg(AF),
	)
}

func (d *Debugger) print(id string) {
	switch id {
	case "bc", "de", "hl", "sp", "pc", "af":
		fmt.Printf("0x%04x\n", d.gb.get16Reg(regNames16[id]))
	case "a", "b", "c", "d", "e", "f", "h", "l":
		fmt.Printf("0x%04x\n", d.gb.get8Reg(regNames8[id]))
	case "regs", "registers":
		d.printAllRegs()
	}
}

func (d *Debugger) prompt() {
	fmt.Printf(">>> ")
}

func (d *Debugger) printMemory(addr, numBytes uint16) {
	/* 8 bytes per line */
	bytes := d.gb.mainMemory.readN(addr, numBytes)
	var i uint16
	for i = 0; i < numBytes; i += 8 {
		fmt.Printf("0x%04x:", addr+i)
		line := bytes[i : i+8]
		if numBytes-i < 8 {
			line = bytes[i:]
		}
		for _, b := range line {
			fmt.Printf(" 0x%02x", b)
		}
		fmt.Printf("\n")
	}
}

func (d *Debugger) printInstructions(addr, numInstructions uint16) {
	d.ROMReader.Seek(int64(addr), 0)
	var i uint16 = 0
	for gbInstruction, addr := DecodeInstruction(d.ROMReader, uint32(addr)); i < numInstructions && gbInstruction != nil; gbInstruction, addr = DecodeInstruction(d.ROMReader, uint32(addr)) {
		fmt.Printf("%s\n", gbInstruction.ToStr())
		i++
	}
}

var print_memory_regex = regexp.MustCompile(`^x/([0-9]*)([xi]*)$`)

func debugLoop(d *Debugger) {
	reader := bufio.NewReader(os.Stdin)
	for {
		/* TODO: handle error */
		d.prompt()
		cmd, _ := reader.ReadString('\n')
		tokens := strings.Fields(strings.ToLower(cmd))
		if len(tokens) == 0 {
			continue
		}
		switch len(tokens) {
		case 1:
			switch tokens[0] {
			case "r", "run":
				d.run()
			case "c", "continue":
				d.cont()
			case "n", "next":
				d.gb.handleInterrupt()
				d.gb.Step()
			case "q", "quit":
				return
			}
		case 2:
			switch tokens[0] {
			case "b", "break":
				addr, err := strconv.ParseUint(tokens[1], 0, 16)
				if err == nil && addr >= 0 && addr <= 0xffff {
					d.addBreakpoint(uint16(addr))
				} else {
					fmt.Printf("Invalid address: %s\n", tokens[1])
				}
			case "d", "delete":
				addr, err := strconv.ParseUint(tokens[1], 0, 16)
				if err == nil && addr >= 0 && addr <= 0xffff {
					d.deleteBreakpoint(uint16(addr))
				} else {
					fmt.Printf("Invalid address: %s\n", tokens[1])
				}
			case "p", "print":
				if len(tokens) > 1 {
					d.print(tokens[1])
				}
			case "x":
				addr, err := strconv.ParseUint(tokens[1], 0, 16)
				if err != nil {
					fmt.Printf("Invalid address: %s\n", tokens[1])
					break
				}
				d.printMemory(uint16(addr), 1)

			default:
				/* boolean switch */
				switch {
				case print_memory_regex.MatchString(tokens[0]):
					/* Print memory with options e.g. x/8x 0x28b */
					options := print_memory_regex.FindAllStringSubmatch(tokens[0], -1)[0]
					addr, err := strconv.ParseUint(tokens[1], 0, 16)
					if err != nil {
						fmt.Printf("Invalid address: %s\n", tokens[1])
						break
					}
					var num uint64 = 1
					if options[1] != "" {
						num, _ = strconv.ParseUint(options[1], 0, 16)
					}
					format := options[2]
					switch format {
					case "i":
						d.printInstructions(uint16(addr), uint16(num))
					case "x":
						fallthrough
					default:
						d.printMemory(uint16(addr), uint16(num))
					}
				}
			}
		}
	}
}

func NewDebugger(gb *GameBoy) *Debugger {
	return &Debugger{
		gb:          gb,
		breakpoints: make(map[uint16]struct{}),
		paused:      true,
		ROMReader:   Gb.mainMemory.cartridge.reader(),
	}
}

var sig_chan = make(chan os.Signal, 1)

var debuggersSignal []*Debugger

func sigint_handler() {
	for {
		<-sig_chan
		for _, d := range debuggersSignal {
			d.paused = true
		}
	}
}
