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
				/* LD [nn], sp */
			case 0x10:
				/*
				 * STOP
				 */
			case 0x18:
				/*
				 * jr E - jump to PC + E
				 */
			case 0x20:
				/* jr nz, nn */
			case 0x28:
				/* jr z, nn */
			case 0x30:
				/* jr nc, nn */
			case 0x38:
				/* jr c, nn */
			}
		case 0x01:
			/* switch on bit 3 */
			switch opCode & 0x08 {
			case 0x00:
				/* ld r16, nn */
			case 0x08:
				/* add hl, r16 */
			}
		case 0x02:
			/* switch on bit 3 */
			switch opCode & 0x08 {
			case 0x00:
				/* switch on bits 4-5 */
				switch opCode & 0x30 {
				case 0x00:
					/* ld [bc], a */
				case 0x10:
					/* ld [de], a */
				case 0x20:
					/* LDI [HL], A */
				case 0x30:
					/* LDD [HL], A */
				}
			case 0x08:
				/* switch on bits 4-5 */
				switch opCode & 0x30 {
				case 0x00:
					/* ld a, [bc] */
				case 0x10:
					/* ld a, [de] */
				case 0x20:
					/* ldi A, [HL] */
				case 0x30:
					/* ldd A, [HL] */
				}
			}
		case 0x03:
			/* switch on bit 3 */
			switch opCode & 0x08 {
			case 0x00:
				/* inc r16 */
			case 0x08:
				/* dec r16 */
			}
		case 0x04:
			/* inc r8 */
		case 0x05:
			/* dec r8 */
		case 0x06:
			/* ld r8, n */
		case 0x07:
			/* switch on bits 3-5 */
			switch opCode & 0x38 {
			case 0x00:
				/* RLCA */
			case 0x08:
				/* RRCA */
			case 0x10:
				/* RLA */
			case 0x18:
				/* RRA */
			case 0x20:
				/* DAA */
			case 0x28:
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
				/* ld r8, [hl] */
			}
		default:
			switch opCode & 0x38 {
			case 0x30:
				/* ld [hl], r8 */
			default:
				/* ld r8, r8 */
			}
		}
	case 0x80:
        switch opCode & 0x38 {
        case 0x00:
            switch opCode & 007 {
            case 0x06:
                /* add A, [HL] */
            default:
                /* add A, r8 */
            }
        case 0x08:
            switch opCode & 007 {
            case 0x06:
                /* adc A, [HL] */
            default:
                /* adc A, r8 */
            }
        case 0x10:
            switch opCode & 007 {
            case 0x06:
                /* sub [HL] */
            default:
                /* sub r8 */
            }
        case 0x18:
            switch opCode & 007 {
            case 0x06:
                /* sbc A, [HL] */
            default:
                /* sbc A, r8 */
            }
        case 0x20:
            switch opCode & 007 {
            case 0x06:
                /* and [HL] */
            default:
                /* and r8 */
            }
        case 0x28:
            switch opCode & 007 {
            case 0x06:
                /* xor [HL] */
            default:
                /* xor r8 */
            }
        case 0x30:
            switch opCode & 007 {
            case 0x06:
                /* or [HL] */
            default:
                /* or r8 */
            }
        case 0x38:
            switch opCode & 007 {
            case 0x06:
                /* cp [HL] */
            default:
                /* cp r8 */
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
				/* ret CC - conditional return */
			case 0x20:
				/* ld [0xff00 + n], A */
			case 0x28:
				/* add SP, n */
			case 0x30:
				/* ld A, [0xff00 + n] */
			case 0x38:
				/* ldhl SP, n */
			}
		case 0x01:
			switch opCode & 0x08 {
			case 0x00:
				/* pop r16 */
			case 0x08:
				switch opCode & 0x30 {
				case 0x00:
					/* ret */
				case 0x10:
					/* reti */
				case 0x20:
					/* jp hl */
				case 0x30:
					/* ld sp, hl */
				}
			}
		case 0x02:
			switch opCode & 0x38 {
			case 0x00:
				fallthrough
			case 0x08:
				fallthrough
			case 0x10:
				fallthrough
			case 0x18:
				/* JP cc, nn (conditional jump) */
			case 0x20:
				/* LD [0xff00 + C], A */
			case 0x28:
				/* LD [nn], A */
			case 0x30:
				/* LD A, [0xff00 + C] */
			case 0x38:
				/* LD A, [nn] */
			}
		case 0x03:
			switch opCode & 0x38 {
			case 0x00:
				/* jp nn */
			case 0x08:
				/* 0xcb prefix */
				prefix := opCode
				opCode = g.mainMemory.read(pc + 1)
				switch opCode & 0xc0 {
				case 0x00:
					/* assorted rotate & shift operations on register or memory */
                    switch opCode & 0x38 {
                    case 0x00:
                        switch opCode & 0x07 {
                        case 0x06:
                            /* rlc [HL] */
                        default:
                            /* rlc r8 */
                        }
                    case 0x08:
                        switch opCode & 0x07 {
                        case 0x06:
                            /* rrc [HL] */
                        default:
                            /* rrc r8 */
                        }
                    case 0x10:
                        switch opCode & 0x07 {
                        case 0x06:
                            /* rl [HL] */
                        default:
                            /* rl r8 */
                        }
                    case 0x18:
                        switch opCode & 0x07 {
                        case 0x06:
                            /* rr [HL] */
                        default:
                            /* rr r8 */
                        }
                    case 0x20:
                        switch opCode & 0x07 {
                        case 0x06:
                            /* sla [HL] */
                        default:
                            /* sla r8 */
                        }
                    case 0x28:
                        switch opCode & 0x07 {
                        case 0x06:
                            /* sra [HL] */
                        default:
                            /* sra r8 */
                        }
                    case 0x30:
                        switch opCode & 0x07 {
                        case 0x06:
                            /* swap [HL] */
                        default:
                            /* swap r8 */
                        }
                    case 0x38:
                        switch opCode & 0x07 {
                        case 0x06:
                            /* srl [HL] */
                        default:
                            /* srl r8 */
                        }
                    }
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
				/* Call cc - conditional call */
			default:
				/* Illegal */
			}
		case 0x05:
			switch opCode & 0x08 {
			case 0x00:
				/* push r16 */
			case 0x08:
				switch opCode & 0x30 {
				case 0x00:
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
                /* add A, nn */
            case 0x08:
                /* adc A, nn */
            case 0x10:
                /* sub nn */
            case 0x18:
                /* sbc a, nn */
            case 0x20:
                /* and nn */
            case 0x28:
                /* xor nn */
            case 0x30:
                /* or nn */
            case 0x38:
                /* cp nn */
            }
		case 0x07:
			/* rst p */
		}
	}
}
