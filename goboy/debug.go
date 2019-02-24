package main

import "fmt"
import "strings"
import "os"
import "bufio"
import "strconv"


type FunctionFrame struct {
    addr uint16 /* Address of frame on stack */
    returnAddr uint16
}

type Debugger struct {
    gb *GameBoy
    breakpoints map[uint16]struct{} /* This is how sets work */
    paused bool
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
                if err == nil && addr >= 0 && addr <= 0xffff{
                    d.addBreakpoint(uint16(addr))
                } else {
                    fmt.Printf("Invalid address: %s\n", tokens[1])
                }
            case "d", "delete":
                addr, err := strconv.ParseUint(tokens[1], 0, 16)
                if err == nil && addr >= 0 && addr <= 0xffff{
                    d.deleteBreakpoint(uint16(addr))
                } else {
                    fmt.Printf("Invalid address: %s\n", tokens[1])
                }
            case "p", "print":
                if len(tokens) > 1 {
                    d.print(tokens[1])
                }
            }
        }
    }
}

var sig_chan = make(chan os.Signal, 1)

var debuggersSignal []*Debugger

func sigint_handler() {
    <-sig_chan
    for _, d := range debuggersSignal {
        d.paused = true
    }
}
