package main

func (g *GameBoy) Step() {
	/* 3 is the max length of an instruction (I think) */
	pc := g.regs[PC]
	opCode := g.mainMemory.read(pc)

	/* Switch on bits 6-7 */
	switch opCode & 0xc0 {
	case 0x00:
		/* Switch on bits 0-2 */
		switch opCode & 0x07 {
		case 0x00:
			/* Switch on bits 3-5 */
			switch opCode & 0x38 {
			case 0x00:
				/* nop */
				instruction := []uint8{opCode}
				g.NOP(instruction)
			case 0x08:
				g.LD_nn_sp(g.mainMemory.readN(pc, 3))
				/* LD [nn], sp */
			case 0x10:
				/*
				 * STOP
				 */
			case 0x18:
				g.JR_e(g.mainMemory.readN(pc, 2))
				/*
				 * jr E - jump to PC + E
				 */
			case 0x20:
				/* jr nz, nn */
				fallthrough
			case 0x28:
				/* jr z, nn */
				fallthrough
			case 0x30:
				/* jr nc, nn */
				fallthrough
			case 0x38:
				g.JR_cc_e(g.mainMemory.readN(pc, 2))
				/* jr c, nn */
			}
		case 0x01:
			/* switch on bit 3 */
			switch opCode & 0x08 {
			case 0x00:
				g.LD_dd_nn(g.mainMemory.readN(pc, 3))
				/* ld r16, nn */
			case 0x08:
				g.ADD_hl_ss(g.mainMemory.readN(pc, 1))
				/* add hl, r16 */
			}
		case 0x02:
			/* switch on bit 3 */
			switch opCode & 0x08 {
			case 0x00:
				/* switch on bits 4-5 */
				switch opCode & 0x30 {
				case 0x00:
					g.LD_bc_a(g.mainMemory.readN(pc, 1))
					/* ld [bc], a */
				case 0x10:
					g.LD_de_a(g.mainMemory.readN(pc, 1))
					/* ld [de], a */
				case 0x20:
					g.LD_hli_a(g.mainMemory.readN(pc, 1))
					/* LDI [HL], A */
				case 0x30:
					g.LD_hld_a(g.mainMemory.readN(pc, 1))
					/* LDD [HL], A */
				}
			case 0x08:
				/* switch on bits 4-5 */
				switch opCode & 0x30 {
				case 0x00:
					g.LD_a_bc(g.mainMemory.readN(pc, 1))
					/* ld a, [bc] */
				case 0x10:
					g.LD_a_de(g.mainMemory.readN(pc, 1))
					/* ld a, [de] */
				case 0x20:
					g.LD_a_hli(g.mainMemory.readN(pc, 1))
					/* ldi A, [HL] */
				case 0x30:
					g.LD_a_hld(g.mainMemory.readN(pc, 1))
					/* ldd A, [HL] */
				}
			}
		case 0x03:
			/* switch on bit 3 */
			switch opCode & 0x08 {
			case 0x00:
				g.INC_ss(g.mainMemory.readN(pc, 1))
				/* inc r16 */
			case 0x08:
				g.DEC_ss(g.mainMemory.readN(pc, 1))
				/* dec r16 */
			}
		case 0x04:
			g.INC_r(g.mainMemory.readN(pc, 1))
			/* inc r8 */
		case 0x05:
			g.DEC_r(g.mainMemory.readN(pc, 1))
			/* dec r8 */
		case 0x06:
			g.LD_r_n(g.mainMemory.readN(pc, 2))
			/* ld r8, n */
		case 0x07:
			/* switch on bits 3-5 */
			switch opCode & 0x38 {
			case 0x00:
				g.RLCA(g.mainMemory.readN(pc, 1))
				/* RLCA */
			case 0x08:
				g.RRCA(g.mainMemory.readN(pc, 1))
				/* RRCA */
			case 0x10:
				g.RLA(g.mainMemory.readN(pc, 1))
				/* RLA */
			case 0x18:
				g.RRA(g.mainMemory.readN(pc, 1))
				/* RRA */
			case 0x20:
				g.DAA(g.mainMemory.readN(pc, 1))
				/* DAA */
			case 0x28:
				g.CPL(g.mainMemory.readN(pc, 1))
				/* CPL */
			case 0x30:
				/* SCF */
			case 0x38:
				/* CCF */
			}
		}
	case 0x40:
		switch opCode & 0x07 {
		case 0x06:
			switch opCode & 0x38 {
			case 0x30:
				/* halt */
			default:
				g.LD_r_hl(g.mainMemory.readN(pc, 1))
				/* ld r8, [hl] */
			}
		default:
			switch opCode & 0x38 {
			case 0x30:
				g.LD_hl_r(g.mainMemory.readN(pc, 1))
				/* ld [hl], r8 */
			default:
				g.LD_r_r(g.mainMemory.readN(pc, 1))
				/* ld r8, r8 */
			}
		}
	case 0x80:
		/* assorted ALU instructions on A and register/memory location */
	case 0xc0:
		switch opCode & 0x07 {
		case 0x00:
			switch opCode & 0x38 {
			case 0x00:
				fallthrough
			case 0x08:
				fallthrough
			case 0x10:
				fallthrough
			case 0x18:
				g.RET_cc(g.mainMemory.readN(pc, 1))
				/* ret CC - conditional return */
			case 0x20:
				g.LD_n_a(g.mainMemory.readN(pc, 2))
				/* ld [0xff00 + n], A */
			case 0x28:
				g.ADD_sp_e(g.mainMemory.readN(pc, 2))
				/* add SP, n */
			case 0x30:
				g.LD_a_n(g.mainMemory.readN(pc, 2))
				/* ld A, [0xff00 + n] */
			case 0x38:
				g.LDHL_sp_e(g.mainMemory.readN(pc, 2))
				/* ldhl SP, n */
			}
		case 0x01:
			switch opCode & 0x08 {
			case 0x00:
				g.POP_qq(g.mainMemory.readN(pc, 1))
				/* pop r16 */
			case 0x08:
				switch opCode & 0x30 {
				case 0x00:
					g.RET(g.mainMemory.readN(pc, 1))
					/* ret */
				case 0x10:
					g.RETI(g.mainMemory.readN(pc, 1))
					/* reti */
				case 0x20:
					g.JP_hl(g.mainMemory.readN(pc, 1))
					/* jp hl */
				case 0x30:
					g.LD_sp_hl(g.mainMemory.readN(pc, 1))
					/* ld sp, hl */
				}
			}
		case 0x02:
			/* jp cc, nn - conditional absolute jump */
			switch opCode & 0x38 {
			case 0x00:
				fallthrough
			case 0x08:
				fallthrough
			case 0x10:
				fallthrough
			case 0x18:
				g.JP_cc_nn(g.mainMemory.readN(pc, 3))
				/* JP cc (conditional jump) */
			case 0x20:
				g.LD_c_a(g.mainMemory.readN(pc, 1))
				/* LD [0xff00 + C], A */
			case 0x28:
				g.LD_nn_a(g.mainMemory.readN(pc, 3))
				/* LD [nn], A */
			case 0x30:
				g.LD_a_c(g.mainMemory.readN(pc, 1))
				/* LD A, [0xff00 + C] */
			case 0x38:
				g.LD_a_nn(g.mainMemory.readN(pc, 3))
				/* LD A, [nn] */
			}
		case 0x03:
			switch opCode & 0x38 {
			case 0x00:
				g.JP_nn(g.mainMemory.readN(pc, 3))
				/* jp nn */
			case 0x08:
				/* 0xcb prefix */
				prefix := opCode
				opCode = g.mainMemory.read(pc + 1)
				switch opCode & 0xc0 {
				case 0x00:
					/* assorted rotate & shift operations on register or memory */
				case 0x40:
					/* bit b, r8 */
					instruction := []uint8{prefix, opCode}
					g.BIT_b_r(instruction)
				case 0x80:
					/* res b, r8 */
				case 0xc0:
					/* set b, r8 */
				}
			case 0x10:
				/* Illegal */
			case 0x18:
				/* Illegal */
			case 0x20:
				/* Illegal */
			case 0x28:
				/* Illegal */
			case 0x30:
				/* di */
			case 0x38:
				/* ei */
			}
		case 0x04:
			/* call cc, nn - conditional call */
			switch opCode & 0x38 {
			case 0x00:
				fallthrough
			case 0x08:
				fallthrough
			case 0x10:
				fallthrough
			case 0x18:
				g.CALL_cc_nn(g.mainMemory.readN(pc, 3))
				/* Call cc - conditional call */
			default:
				/* Illegal */
			}
		case 0x05:
			switch opCode & 0x08 {
			case 0x00:
				g.PUSH_qq(g.mainMemory.readN(pc, 1))
				/* push r16 */
			case 0x08:
				switch opCode & 0x30 {
				case 0x00:
					g.CALL_nn(g.mainMemory.readN(pc, 3))
					/* call nn */
				case 0x10:
					/* Illegal */
				case 0x20:
					/* Illegal */
				case 0x30:
					/* Illegal */
				}
			}
		case 0x06:
			/* assorted ALU instructions on A and immediate operand */
		case 0x07:
			g.RST(g.mainMemory.readN(pc, 1))
			/* rst p */
		}
	}
}
