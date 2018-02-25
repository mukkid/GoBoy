package GoBoy

//Frances was here!

type GameBoy struct {
    mram [8 * 1024]uint8
    vram [8 * 1024]uint8
}

type Reg8ID int
type Reg16ID int

const (
    A Reg8ID = 0
    F Reg8ID = 1
    B Reg8ID = 2
    C Reg8ID = 3
    D Reg8ID = 4
    E Reg8ID = 5
    H Reg8ID = 6
    L Reg8ID = 7
)

const (
    AF Reg16ID = 0
    BC Reg16ID = 1
    DE Reg16ID = 2
    HL Reg16ID = 3
    SP Reg16ID = 4
    PC Reg16ID = 5
)


type Register struct {
    regs [6]uint16
}

func (r *Register) get16Reg(id Reg16ID) uint16 {
    return r.regs[id]
}

func (r *Register) set16Reg(id Reg16ID, value uint16) {
    r.regs[id] = value
}

func (r *Register) get8Reg(id Reg8ID) uint8 {
    block := id / 2  // block is the index of the 16bit version
    end := id % 2  // end indicates the high (0) or low (1) byte of the 16bit register
    if end == 0 {
        return uint8(r.regs[block] >> 8)
    } else {
        return uint8(r.regs[block])
    }
}

func (r *Register) set8Reg(id Reg8ID, value uint8) {
    block := id / 2
    end := id % 2
    if end == 0 {
        r.regs[block] &= 0x00ff
        r.regs[block] |= uint16(value << 8)
    } else {
        r.regs[block] &= 0xff00
        r.regs[block] |= uint16(value)
    }
}
