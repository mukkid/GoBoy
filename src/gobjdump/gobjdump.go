package main

import (
    "bufio"
    "encoding/hex"
    "fmt"
    "io"
    "strings"
)

var r8 = []string{
    "b",
    "c",
    "d",
    "e",
    "h",
    "l",
    "[hl]",
    "a",
}

var r16_sp = []string{
    "bc",
    "de",
    "hl",
    "sp",
}

var r16_af = []string{
    "bc",
    "de",
    "hl",
    "af",
}

var conditions = []string{
    "NZ",
    "Z",
    "NC",
    "C",
    "PO",
    "PE",
    "P",
    "M",
}

/* Consumes an immediate 8 bit value from the stream, updates the args buffer with it */
func imm8(r *bufio.Reader, instruction *[]uint8) string {
    nextByte, err := r.ReadByte()
    if err != nil {
        return ""
    }
    *instruction = append(*instruction, nextByte)
    return fmt.Sprintf("0x%02x", nextByte)
}

/* Consumes a signed immediate 8 bit value from the stream, updates the args buffer with it */
func imm8_s(r *bufio.Reader, instruction *[]uint8) string {
    nextByte, err := r.ReadByte()
    if err != nil {
        return ""
    }
    *instruction = append(*instruction, nextByte)
    return fmt.Sprintf("%d", int8(nextByte))
}

func imm16(r *bufio.Reader, instruction *[]uint8) string {
    imm := make([]uint8, 2)
    _, err := io.ReadFull(r, imm)
    if err != nil {
        return ""
    }
    *instruction = append(*instruction, imm[0], imm[1])
    return fmt.Sprintf("0x%02x%02x", imm[1], imm[0])
}

func imm16_addr(r *bufio.Reader, instruction *[]uint8) string {
    imm := make([]uint8, 2)
    _, err := io.ReadFull(r, imm)
    if err != nil {
        return ""
    }
    *instruction = append(*instruction, imm[0], imm[1])
    return fmt.Sprintf("[0x%02x%02x]", imm[1], imm[0])
}

func r16_af_addr(r *bufio.Reader, instruction *[]uint8) string {
    reg_index := ((*instruction)[0] & 0x30) >> 4
    return fmt.Sprintf("[%s]", r16_af[reg_index])
}

func r16_sp_addr(r *bufio.Reader, instruction *[]uint8) string {
    reg_index := ((*instruction)[0] & 0x30) >> 4
    return fmt.Sprintf("[%s]", r16_sp[reg_index])
}

func decodeDJNZ(r *bufio.Reader, instruction *[]uint8, mnemonic *[]string) {
    *mnemonic = append(*mnemonic, "djnz")
    /* Read operand (next byte) */
    *mnemonic = append(*mnemonic, imm8_s(r, instruction))
}

func decodeJR_E(r *bufio.Reader, instruction *[]uint8, mnemonic *[]string) {
    *mnemonic = append(*mnemonic, "jr")
    /* Read operand (next byte) */
    *mnemonic = append(*mnemonic, imm8(r, instruction))
}

func decodeJR_cond_E(r *bufio.Reader, instruction *[]uint8, mnemonic *[]string) {
    *mnemonic = append(*mnemonic, "jr")
    cond_index := ((*instruction)[0] & 0x38) >> 3 - 4
    *mnemonic = append(*mnemonic, conditions[cond_index])
    *mnemonic = append(*mnemonic, imm8(r, instruction))
}

func decodeLD_r16_nn(r *bufio.Reader, instruction *[]uint8, mnemonic *[]string) {
    reg_index := ((*instruction)[0] & 0x30) >> 4
    *mnemonic = append(*mnemonic, "ld")
    *mnemonic = append(*mnemonic, r16_sp[reg_index])
    *mnemonic = append(*mnemonic, imm16(r, instruction))
}

func decodeADD_hl_r16(r *bufio.Reader, instruction *[]uint8, mnemonic *[]string) {
    reg_index := ((*instruction)[0] & 0x30) >> 4
    *mnemonic = append(*mnemonic, "add")
    *mnemonic = append(*mnemonic, "hl")
    *mnemonic = append(*mnemonic, (r16_sp[reg_index]))
}

func decodeLD_BC_A(r *bufio.Reader, instruction *[]uint8, mnemonic *[]string) {
    *mnemonic = append(*mnemonic, "ld")
    *mnemonic = append(*mnemonic, "[bc]")
    *mnemonic = append(*mnemonic, "a")
}

func decodeLD_DE_A(r *bufio.Reader, instruction *[]uint8, mnemonic *[]string) {
    *mnemonic = append(*mnemonic, "ld")
    *mnemonic = append(*mnemonic, "[de]")
    *mnemonic = append(*mnemonic, "a")
}

func decodeLD_nn_HL(r *bufio.Reader, instruction *[]uint8, mnemonic *[]string) {
    *mnemonic = append(*mnemonic, "ld")
    *mnemonic = append(*mnemonic, imm16_addr(r, instruction))
    *mnemonic = append(*mnemonic, "hl")
}

func decodeLD_nn_A(r *bufio.Reader, instruction *[]uint8, mnemonic *[]string) {
    *mnemonic = append(*mnemonic, "ld")
    *mnemonic = append(*mnemonic, imm16_addr(r, instruction))
    *mnemonic = append(*mnemonic, "a")
}

func decodeLD_A_BC(r *bufio.Reader, instruction *[]uint8, mnemonic *[]string) {
    *mnemonic = append(*mnemonic, "ld")
    *mnemonic = append(*mnemonic, "a")
    *mnemonic = append(*mnemonic, "[bc]")
}

func decodeLD_A_DE(r *bufio.Reader, instruction *[]uint8, mnemonic *[]string) {
    *mnemonic = append(*mnemonic, "ld")
    *mnemonic = append(*mnemonic, "a")
    *mnemonic = append(*mnemonic, "[de]")
}

func decodeLD_HL_nn(r *bufio.Reader, instruction *[]uint8, mnemonic *[]string) {
    *mnemonic = append(*mnemonic, "ld")
    *mnemonic = append(*mnemonic, "hl")
    *mnemonic = append(*mnemonic, imm16_addr(r, instruction))
}

func decodeLD_A_nn(r *bufio.Reader, instruction *[]uint8, mnemonic *[]string) {
    *mnemonic = append(*mnemonic, "ld")
    *mnemonic = append(*mnemonic, "a")
    *mnemonic = append(*mnemonic, imm16_addr(r, instruction))
}

func decodeINC_r16(r *bufio.Reader, instruction *[]uint8, mnemonic *[]string) {
    reg_index := ((*instruction)[0] & 0x30) >> 4
    *mnemonic = append(*mnemonic, "inc")
    *mnemonic = append(*mnemonic, r16_sp[reg_index])
}

func decodeDEC_r16(r *bufio.Reader, instruction *[]uint8, mnemonic *[]string) {
    reg_index := ((*instruction)[0] & 0x30) >> 4
    *mnemonic = append(*mnemonic, "dec")
    *mnemonic = append(*mnemonic, r16_sp[reg_index])
}

func decodeINC_r8(r *bufio.Reader, instruction *[]uint8, mnemonic *[]string) {
    reg_index := ((*instruction)[0] & 0x38) >> 3
    *mnemonic = append(*mnemonic, "inc")
    *mnemonic = append(*mnemonic, r8[reg_index])
}

func decodeDEC_r8(r *bufio.Reader, instruction *[]uint8, mnemonic *[]string) {
    reg_index := ((*instruction)[0] & 0x38) >> 3
    *mnemonic = append(*mnemonic, "dec")
    *mnemonic = append(*mnemonic, r8[reg_index])
}

func decodeLD_r8_n(r *bufio.Reader, instruction *[]uint8, mnemonic *[]string) {
    reg_index := ((*instruction)[0] & 0x38) >> 3
    *mnemonic = append(*mnemonic, "ld")
    *mnemonic = append(*mnemonic, r8[reg_index])
    *mnemonic = append(*mnemonic, imm8(r, instruction))
}

/*
 * Bumps the pointer in r
 * returns: the instruction bytes, the instruction mnemonic as an array of tokens
 */
func decodeInstruction(r *bufio.Reader) ([]uint8, []string) {
    /* If EOF, return empty string */
    var instruction []uint8
    nextByte, err := r.ReadByte()
    if err != nil {
        return nil, nil
    }

    instruction = append(instruction, nextByte)
    var mnemonic []string;

    switch nextByte {
    /* prefix */
    case 0xcc:
    case 0xdd:
    case 0xed:
    case 0xfd:
    /* unprefixed */
    default:
        /* Switch on bits 6-7 */
        switch nextByte & 0xc0 {
            case 0x00:
                /* Switch on bits 0-2 */
                switch nextByte & 0x07 {
                case 0x00:
                    /* Switch on bits 3-5 */
                    switch nextByte & 0x38 {
                    case 0x00:
                        /* nop */
                        mnemonic = append(mnemonic, "nop")
                    case 0x08:
                        /* ex af,af' */
                        /* Not implemented in GB */
                        mnemonic = append(mnemonic, "ex")
                        mnemonic = append(mnemonic, "af")
                        mnemonic = append(mnemonic, "af'")
                    case 0x10:
                        /*
                         * djnz x
                         * Not implemented in GB
                         */
                         decodeDJNZ(r, &instruction, &mnemonic)
                    case 0x18:
                        /*
                         * jr E - jump to PC + E
                         */
                         decodeJR_E(r, &instruction, &mnemonic)
                    default:
                        /* jr nz|z|nc|c, E*/
                        decodeJR_cond_E(r, &instruction, &mnemonic)
                    }
                case 0x01:
                    /* switch on bit 3 */
                    switch nextByte & 0x08 {
                    case 0x00:
                        /* ld rp[p], nn */
                        decodeLD_r16_nn(r, &instruction, &mnemonic)
                    case 0x08:
                        /* add hl, rp[p] */
                        decodeADD_hl_r16(r, &instruction, &mnemonic)
                    }
                case 0x02:
                    /* switch on bit 3 */
                    switch nextByte & 0x08 {
                    case 0x00:
                        /* switch on bits 4-5 */
                        switch nextByte & 0x30 {
                            case 0x00:
                                /* ld [bc], a */
                                decodeLD_BC_A(r, &instruction, &mnemonic)
                            case 0x10:
                                /* ld [de], a */
                                decodeLD_DE_A(r, &instruction, &mnemonic)
                            case 0x20:
                                /* ld [nn], hl */
                                decodeLD_nn_HL(r, &instruction, &mnemonic)
                            case 0x30:
                                /* ld [nn], a */
                                decodeLD_nn_A(r, &instruction, &mnemonic)
                        }
                    case 0x08:
                        /* switch on bits 4-5 */
                        switch nextByte & 0x30 {
                        case 0x00:
                            /* ld a, [bc] */
                            decodeLD_A_BC(r, &instruction, &mnemonic)
                        case 0x10:
                            /* ld a, [de] */
                            decodeLD_A_DE(r, &instruction, &mnemonic)
                        case 0x20:
                            /* ld hl, [nn] */
                            decodeLD_HL_nn(r, &instruction, &mnemonic)
                        case 0x30:
                            /* ld a, [nn] */
                            decodeLD_A_nn(r, &instruction, &mnemonic)
                        }
                    }
                case 0x03:
                    /* switch on bit 3 */
                    switch nextByte & 0x08 {
                    case 0x00:
                        /* inc r16 */
                        decodeINC_r16(r, &instruction, &mnemonic)
                    case 0x08:
                        /* dec r16 */
                        decodeDEC_r16(r, &instruction, &mnemonic)
                    }
                case 0x04:
                    /* inc r8 */
                    decodeINC_r8(r, &instruction, &mnemonic)
                case 0x05:
                    /* dec r8 */
                    decodeDEC_r8(r, &instruction, &mnemonic)
                case 0x06:
                    /* ld r8, n */
                    decodeLD_r8_n(r, &instruction, &mnemonic)
                case 0x07:
                    /* switch on bits 3-5 */
                    switch nextByte & 0x38 {
                    case 0x00:
                        /* RLCA */
                        mnemonic = append(mnemonic, "rlca")
                    case 0x08:
                        /* RRCA */
                        mnemonic = append(mnemonic, "rrca")
                    case 0x10:
                        /* RLA */
                        mnemonic = append(mnemonic, "rla")
                    case 0x18:
                        /* RRA */
                        mnemonic = append(mnemonic, "rra")
                    case 0x20:
                        /* DAA */
                        mnemonic = append(mnemonic, "daa")
                    case 0x28:
                        /* CPL */
                        mnemonic = append(mnemonic, "cpl")
                    case 0x30:
                        /* SCF */
                        mnemonic = append(mnemonic, "scf")
                    case 0x38:
                        /* CCF */
                        mnemonic = append(mnemonic, "ccf")
                    }
                }
            case 0x40:
            case 0x80:
            case 0xc0:
        }
    }
    return instruction, mnemonic
}

func disassemblerLoop(r *bufio.Reader) {
    var addr uint32 = 0x0
    for instruction, mnemonic := decodeInstruction(r);
        len(instruction) != 0;
        instruction, mnemonic = decodeInstruction(r) {
        /* Generate hex encoding of instruction */
        instructionHex := make([]uint8, hex.EncodedLen(len(instruction)))
        hex.Encode(instructionHex, instruction)

        /* format - addr: <instruction bytes> <instruction mnemonic> */
        operands := ""
        if len(mnemonic) > 1 {
            operands = strings.Join(mnemonic[1:], ", ")
        }
        fmt.Printf("0x%016x: %-12s %-6s %s\n", addr, instructionHex, mnemonic[0], operands)
        addr += uint32(len(instruction))
    }
}
