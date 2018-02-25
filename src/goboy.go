package GoBoy

//Frances was here!

type GameBoy struct {
    mram [8 * 1024]uint8
    vram [8 * 1024]uint8
}

type Register struct {
    regs [6]uint16
}

func (r *Register) getA() uint8 {
    return uint8(r.regs[0] >> 8)
}

func (r *Register) getF() uint8 {
    return uint8(r.regs[0])
}

func (r *Register) getB() uint8 {
    return uint8(r.regs[1] >> 8)
}

func (r *Register) getC() uint8 {
    return uint8(r.regs[1])
}

func (r *Register) getD() uint8 {
    return uint8(r.regs[2] >> 8)
}

func (r *Register) getE() uint8 {
    return uint8(r.regs[2])
}

func (r *Register) getH() uint8 {
    return uint8(r.regs[3] >> 8)
}

func (r *Register) getL() uint8 {
    return uint8(r.regs[3])
}

func (r *Register) getSP() uint16 {
    return r.regs[4]
}

func (r *Register) getPC() uint16 {
    return r.regs[5]
}

func (r *Register) getAF() uint16 {
    return r.regs[0]
}

func (r *Register) getBC() uint16 {
    return r.regs[1]
}

func (r *Register) getDE() uint16 {
    return r.regs[2]
}

func (r *Register) getHL() uint16 {
    return r.regs[3]
}
