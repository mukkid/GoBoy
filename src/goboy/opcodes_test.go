package GoBoy

import "testing"

func TestLD_r_r(t *testing.T) {
    gb := &GameBoy{
        Register: &Register{},
    }
    gb.set8Reg(B, 0xfe)
    gb.LD_r_r(0x50) // instruction for LD D, B
    if gb.get8Reg(B) != 0xfe {
        t.Fatal()
    }
    if gb.get8Reg(D) != 0xfe {
        t.Fatal()
    }
    if gb.get16Reg(PC) != 0x0001 {
        t.Fatal()
    }
}
