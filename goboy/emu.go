package main

func (g *GameBoy) Step() {
	/* 3 is the max length of an instruction (I think) */
	pc := g.regs[PC]
	opCode := g.mainMemory.read(pc)
	var cycles int = 0
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
				cycles = g.NOP(instruction)
			case 0x08:
				cycles = g.LD_nn_sp(g.mainMemory.readN(pc, 3))
				/* LD [nn], sp */
			case 0x10:
				/*
				 * STOP
				 */
			case 0x18:
				cycles = g.JR_e(g.mainMemory.readN(pc, 2))
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
				cycles = g.JR_cc_e(g.mainMemory.readN(pc, 2))
				/* jr c, nn */
			}
		case 0x01:
			/* switch on bit 3 */
			switch opCode & 0x08 {
			case 0x00:
				cycles = g.LD_dd_nn(g.mainMemory.readN(pc, 3))
				/* ld r16, nn */
			case 0x08:
				cycles = g.ADD_hl_ss(g.mainMemory.readN(pc, 1))
				/* add hl, r16 */
			}
		case 0x02:
			/* switch on bit 3 */
			switch opCode & 0x08 {
			case 0x00:
				/* switch on bits 4-5 */
				switch opCode & 0x30 {
				case 0x00:
					cycles = g.LD_bc_a(g.mainMemory.readN(pc, 1))
					/* ld [bc], a */
				case 0x10:
					cycles = g.LD_de_a(g.mainMemory.readN(pc, 1))
					/* ld [de], a */
				case 0x20:
					cycles = g.LDI_hl_a(g.mainMemory.readN(pc, 1))
					/* LDI [HL], A */
				case 0x30:
					cycles = g.LDD_hl_a(g.mainMemory.readN(pc, 1))
					/* LDD [HL], A */
				}
			case 0x08:
				/* switch on bits 4-5 */
				switch opCode & 0x30 {
				case 0x00:
					cycles = g.LD_a_bc(g.mainMemory.readN(pc, 1))
					/* ld a, [bc] */
				case 0x10:
					cycles = g.LD_a_de(g.mainMemory.readN(pc, 1))
					/* ld a, [de] */
				case 0x20:
					cycles = g.LDI_a_hl(g.mainMemory.readN(pc, 1))
					/* ldi A, [HL] */
				case 0x30:
					cycles = g.LDD_a_hl(g.mainMemory.readN(pc, 1))
					/* ldd A, [HL] */
				}
			}
		case 0x03:
			/* switch on bit 3 */
			switch opCode & 0x08 {
			case 0x00:
				cycles = g.INC_ss(g.mainMemory.readN(pc, 1))
				/* inc r16 */
			case 0x08:
				cycles = g.DEC_ss(g.mainMemory.readN(pc, 1))
				/* dec r16 */
			}
		case 0x04:
			cycles = g.INC_r(g.mainMemory.readN(pc, 1))
			/* inc r8 */
		case 0x05:
			cycles = g.DEC_r(g.mainMemory.readN(pc, 1))
			/* dec r8 */
		case 0x06:
			cycles = g.LD_r_n(g.mainMemory.readN(pc, 2))
			/* ld r8, n */
		case 0x07:
			/* switch on bits 3-5 */
			switch opCode & 0x38 {
			case 0x00:
				cycles = g.RLCA(g.mainMemory.readN(pc, 1))
				/* RLCA */
			case 0x08:
				cycles = g.RRCA(g.mainMemory.readN(pc, 1))
				/* RRCA */
			case 0x10:
				cycles = g.RLA(g.mainMemory.readN(pc, 1))
				/* RLA */
			case 0x18:
				cycles = g.RRA(g.mainMemory.readN(pc, 1))
				/* RRA */
			case 0x20:
				cycles = g.DAA(g.mainMemory.readN(pc, 1))
				/* DAA */
			case 0x28:
				cycles = g.CPL(g.mainMemory.readN(pc, 1))
				/* CPL */
			case 0x30:
				cycles = g.SCF(g.mainMemory.readN(pc, 1))
				/* SCF */
			case 0x38:
				cycles = g.CCF(g.mainMemory.readN(pc, 1))
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
				cycles = g.LD_r_hl(g.mainMemory.readN(pc, 1))
				/* ld r8, [hl] */
			}
		default:
			switch opCode & 0x38 {
			case 0x30:
				cycles = g.LD_hl_r(g.mainMemory.readN(pc, 1))
				/* ld [hl], r8 */
			default:
				cycles = g.LD_r_r(g.mainMemory.readN(pc, 1))
				/* ld r8, r8 */
			}
		}
	case 0x80:
		/* assorted ALU instructions on A and register/memory location */
		switch opCode & 0x38 {
		case 0x00:
			// most ADD instructions
			switch opCode & 0x07 {
			case 0x06:
				cycles = g.ADD_a_hl(g.mainMemory.readN(pc, 1))
				// add A, [HL]
			default:
				// add A, r
				cycles = g.ADD_a_r(g.mainMemory.readN(pc, 1))
			}
		case 0x08:
			// most ADC instructions
			switch opCode & 0x07 {
			case 0x06:
				// adc A, [HL]
				cycles = g.ADC_a_hl(g.mainMemory.readN(pc, 1))
			default:
				cycles = g.ADC_a_r(g.mainMemory.readN(pc, 1))
			}
		case 0x10:
			// SUB instructions
			switch opCode & 0x07 {
			case 0x06:
				cycles = g.SUB_a_hl(g.mainMemory.readN(pc, 1))
			default:
				cycles = g.SUB_a_r(g.mainMemory.readN(pc, 1))
			}
		case 0x18:
			// SBC instructions
			switch opCode & 0x07 {
			case 0x06:
				cycles = g.SBC_a_hl(g.mainMemory.readN(pc, 1))
			default:
				cycles = g.SUB_a_r(g.mainMemory.readN(pc, 1))
			}
		case 0x20:
			// AND instructions
			switch opCode & 0x07 {
			case 0x06:
				cycles = g.AND_a_hl(g.mainMemory.readN(pc, 1))
			default:
				cycles = g.AND_a_r(g.mainMemory.readN(pc, 1))
			}
		case 0x28:
			// XOR instructions
			switch opCode & 0x07 {
			case 0x06:
				cycles = g.XOR_a_hl(g.mainMemory.readN(pc, 1))
			default:
				cycles = g.XOR_a_r(g.mainMemory.readN(pc, 1))
			}
		case 0x30:
			// OR instructions
			switch opCode & 0x07 {
			case 0x06:
				cycles = g.OR_a_hl(g.mainMemory.readN(pc, 1))
			default:
				cycles = g.OR_a_r(g.mainMemory.readN(pc, 1))
			}
		case 0x38:
			switch opCode & 0x07 {
			case 0x06:
				cycles = g.CP_a_hl(g.mainMemory.readN(pc, 1))
			default:
				cycles = g.CP_a_r(g.mainMemory.readN(pc, 1))
			}
		}
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
				cycles = g.RET_cc(g.mainMemory.readN(pc, 1))
				/* ret CC - conditional return */
			case 0x20:
				cycles = g.LD_n_a(g.mainMemory.readN(pc, 2))
				/* ld [0xff00 + n], A */
			case 0x28:
				cycles = g.ADD_sp_e(g.mainMemory.readN(pc, 2))
				/* add SP, n */
			case 0x30:
				cycles = g.LD_a_n(g.mainMemory.readN(pc, 2))
				/* ld A, [0xff00 + n] */
			case 0x38:
				cycles = g.LDHL_sp_e(g.mainMemory.readN(pc, 2))
				/* ldhl SP, n */
			}
		case 0x01:
			switch opCode & 0x08 {
			case 0x00:
				cycles = g.POP_qq(g.mainMemory.readN(pc, 1))
				/* pop r16 */
			case 0x08:
				switch opCode & 0x30 {
				case 0x00:
					cycles = g.RET(g.mainMemory.readN(pc, 1))
					/* ret */
				case 0x10:
					cycles = g.RETI(g.mainMemory.readN(pc, 1))
					/* reti */
				case 0x20:
					cycles = g.JP_hl(g.mainMemory.readN(pc, 1))
					/* jp hl */
				case 0x30:
					cycles = g.LD_sp_hl(g.mainMemory.readN(pc, 1))
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
				cycles = g.JP_cc_nn(g.mainMemory.readN(pc, 3))
				/* JP cc (conditional jump) */
			case 0x20:
				cycles = g.LD_c_a(g.mainMemory.readN(pc, 1))
				/* LD [0xff00 + C], A */
			case 0x28:
				cycles = g.LD_nn_a(g.mainMemory.readN(pc, 3))
				/* LD [nn], A */
			case 0x30:
				cycles = g.LD_a_c(g.mainMemory.readN(pc, 1))
				/* LD A, [0xff00 + C] */
			case 0x38:
				cycles = g.LD_a_nn(g.mainMemory.readN(pc, 3))
				/* LD A, [nn] */
			}
		case 0x03:
			switch opCode & 0x38 {
			case 0x00:
				cycles = g.JP_nn(g.mainMemory.readN(pc, 3))
				/* jp nn */
			case 0x08:
				/* 0xcb prefix */
				opCode = g.mainMemory.read(pc + 1)
				switch opCode & 0xc0 {
				case 0x00:
					/* assorted rotate & shift operations on register or memory */
					switch opCode & 0x38 {
					case 0x00:
						switch opCode & 0x07 {
						case 0x06:
							cycles = g.RLC_hl(g.mainMemory.readN(pc, 2))
						default:
							cycles = g.RLC_r(g.mainMemory.readN(pc, 2))
						}
					case 0x08:
						switch opCode & 0x07 {
						case 0x06:
							cycles = g.RRC_hl(g.mainMemory.readN(pc, 2))
						default:
							cycles = g.RRC_r(g.mainMemory.readN(pc, 2))
						}
					case 0x10:
						switch opCode & 0x07 {
						case 0x06:
							cycles = g.RL_hl(g.mainMemory.readN(pc, 2))
						default:
							cycles = g.RL_r(g.mainMemory.readN(pc, 2))
						}
					case 0x18:
						switch opCode & 0x07 {
						case 0x06:
							cycles = g.RR_hl(g.mainMemory.readN(pc, 2))
						default:
							cycles = g.RR_r(g.mainMemory.readN(pc, 2))
						}
					case 0x20:
						switch opCode & 0x07 {
						case 0x06:
							cycles = g.SLA_hl(g.mainMemory.readN(pc, 2))
						default:
							cycles = g.SLA_r(g.mainMemory.readN(pc, 2))
						}
					case 0x28:
						switch opCode & 0x07 {
						case 0x06:
							cycles = g.SRA_hl(g.mainMemory.readN(pc, 2))
						default:
							cycles = g.SRA_r(g.mainMemory.readN(pc, 2))
						}
					case 0x30:
						switch opCode & 0x07 {
						case 0x06:
							cycles = g.SWAP_hl(g.mainMemory.readN(pc, 2))
						default:
							cycles = g.SWAP_r(g.mainMemory.readN(pc, 2))
						}
					case 0x38:
						switch opCode & 0x07 {
						case 0x06:
							cycles = g.SRL_hl(g.mainMemory.readN(pc, 2))
						default:
							cycles = g.SRL_r(g.mainMemory.readN(pc, 2))
						}
					}
				case 0x40:
					/* bit b, r8 */
					switch opCode & 0x07 {
					case 0x06:
						cycles = g.BIT_b_hl(g.mainMemory.readN(pc, 2))
					default:
						cycles = g.BIT_b_r(g.mainMemory.readN(pc, 2))
					}
				case 0x80:
					/* res b, r8 */
					switch opCode & 0x07 {
					case 0x06:
						cycles = g.RES_b_hl(g.mainMemory.readN(pc, 2))
					default:
						cycles = g.RES_b_r(g.mainMemory.readN(pc, 2))
					}
				case 0xc0:
					/* set b, r8 */
					switch opCode & 0x07 {
					case 0x06:
						cycles = g.SET_b_hl(g.mainMemory.readN(pc, 2))
					default:
						cycles = g.SET_b_r(g.mainMemory.readN(pc, 2))
					}
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
				cycles = g.DI(g.mainMemory.readN(pc, 1))
				/* di */
			case 0x38:
				cycles = g.EI(g.mainMemory.readN(pc, 1))
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
				cycles = g.CALL_cc_nn(g.mainMemory.readN(pc, 3))
				/* Call cc - conditional call */
			default:
				/* Illegal */
			}
		case 0x05:
			switch opCode & 0x08 {
			case 0x00:
				cycles = g.PUSH_qq(g.mainMemory.readN(pc, 1))
				/* push r16 */
			case 0x08:
				switch opCode & 0x30 {
				case 0x00:
					cycles = g.CALL_nn(g.mainMemory.readN(pc, 3))
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
			switch opCode & 0x38 {
			case 0x00:
				cycles = g.ADD_a_n(g.mainMemory.readN(pc, 2))
			case 0x08:
				cycles = g.ADC_a_n(g.mainMemory.readN(pc, 2))
			case 0x10:
				cycles = g.SUB_a_n(g.mainMemory.readN(pc, 2))
			case 0x18:
				cycles = g.SBC_a_n(g.mainMemory.readN(pc, 2))
			case 0x20:
				cycles = g.AND_a_n(g.mainMemory.readN(pc, 2))
			case 0x28:
				cycles = g.XOR_a_n(g.mainMemory.readN(pc, 2))
			case 0x30:
				cycles = g.OR_a_n(g.mainMemory.readN(pc, 2))
			case 0x38:
				cycles = g.CP_a_n(g.mainMemory.readN(pc, 2))
			}
		case 0x07:
			cycles = g.RST(g.mainMemory.readN(pc, 1))
			/* rst p */
		}
	}

	g.TSCStart += uint64(cycles)

	/*
	 * This loop could occur at the top of the function as well
	 * Do not start executing instruction until TSCStart + cycles
	 * We need to check g.Paused in case the Debugger is paused
	 * with SIGINT, otherwise this may infinitely loop.
	 */
	for g.TSC < g.TSCStart && !g.Paused {
	}
}

func (g *GameBoy) handleInterrupt() {
	if !g.interruptEnabled {
		return
	}
	interrupts_enabled := g.mainMemory.read(0xffff)
	interrupts_request := g.mainMemory.read(0xff0f)
	interrupts := interrupts_enabled & interrupts_request
	if interrupts > 0x00 {
		bit := interrupts & -interrupts
		switch bit {
		case 0x01:
			// VBLANK
			g.interruptJumpHelper(0x0040)
		case 0x02:
			// LCD Stat
			g.interruptJumpHelper(0x0048)
		case 0x04:
			// TIMER
			g.interruptJumpHelper(0x0050)
		case 0x08:
			// SERIAL
			g.interruptJumpHelper(0x0058)
		case 0x10:
			// JOYPAD
			g.interruptJumpHelper(0x0060)
		}
		interrupts_request ^= bit
		g.mainMemory.write(0xff0f, interrupts_request)
	}
}

func (g *GameBoy) interruptJumpHelper(target uint16) {
	g.interruptEnabled = false
	val := g.get16Reg(PC)
	lowVal := uint8(val)
	highVal := uint8(val >> 8)
	g.regs[SP] -= 1
	g.mainMemory.write(g.get16Reg(SP), highVal)
	g.regs[SP] -= 1
	g.mainMemory.write(g.get16Reg(SP), lowVal)
	g.regs[PC] = target
}

/* TODO: investigate good ways to test these Ticker loops */
func (Gb *GameBoy) LCDLoop() {
	for _ = range Gb.LCDClock.C {
		// NOTE: Check debugger pause flag here
		if !Gb.Paused {
			/* LY - 0xff44 */
			LY := Gb.mainMemory.ioregs[0x44]
			LY = (LY + 1) % 0x9a // LY increments from 0 (0x00) to 153 (0x99) and then repeats
			Gb.mainMemory.ioregs[0x44] = LY
			/*
			 * LYC  - 0xff45
			 * STAT - 0xff41
			 * LYC and LC are continuously compared with each other. When
			 * both values are identical, the coincident bit in the STAT register becomes
			 * set, and (if enabled) a STAT interrupt is requested.
			 */
			LYC := Gb.mainMemory.ioregs[0x45]
			if LYC == LY {
				/* Set bit 6 in STAT */
				STAT := Gb.mainMemory.ioregs[0x41]
				STAT |= 0x40
			}

		}
	}
}

func (Gb *GameBoy) TSCLoop() {
	for _ = range Gb.CPUClock.C {
		if !Gb.Paused {
			Gb.TSC += 4
		}
	}
}
