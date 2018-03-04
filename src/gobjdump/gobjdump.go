package main

import (
    "bufio"
    "encoding/hex"
    "fmt"
    "io"
    "strings"
)

type Z80AsmErrorType uint8

const (
    Z80AsmErrorIllegalInstruction Z80AsmErrorType = iota
    Z80AsmErrorUnimplementedInstruction
    Z80AsmErrorMalformedInstruction
    Z80AsmErrorUnknown
)

type Z80AsmError struct {
    errorType Z80AsmErrorType
}

func (e *Z80AsmError) Error() string {
    switch e.errorType {
    case Z80AsmErrorIllegalInstruction:
        return "Error: Illegal Instruction"
    case Z80AsmErrorUnimplementedInstruction:
        return "Error: Unimplemented Instruction"
    case Z80AsmErrorMalformedInstruction:
        return "Error: Malformed Instruction"
    default:
        return "Error: Unknown"
    }
}

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

var ALU = [][]string{
    []string{"add", "a"},
    []string{"adc", "a"},
    []string{"sub"},
    []string{"sbc", "a"},
    []string{"and"},
    []string{"xor"},
    []string{"or"},
    []string{"cp"},
}

/* Consumes an immediate 8 bit value from the stream, updates the args buffer with it */
func imm8(r *bufio.Reader, instruction *[]uint8) (string, error) {
    nextByte, err := r.ReadByte()
    if err != nil {
        if err == io.EOF {
            return "", &Z80AsmError{errorType: Z80AsmErrorMalformedInstruction}
        } else {
            return "", &Z80AsmError{errorType: Z80AsmErrorUnknown}
        }
    }
    *instruction = append(*instruction, nextByte)
    return fmt.Sprintf("0x%02x", nextByte), nil
}

/* Consumes a signed immediate 8 bit value from the stream, updates the args buffer with it */
func imm8_s(r *bufio.Reader, instruction *[]uint8) (string, error) {
    nextByte, err := r.ReadByte()
    if err != nil {
        if err == io.EOF {
            return "", &Z80AsmError{errorType: Z80AsmErrorMalformedInstruction}
        } else {
            return "", &Z80AsmError{errorType: Z80AsmErrorUnknown}
        }
    }
    *instruction = append(*instruction, nextByte)
    return fmt.Sprintf("%d", int8(nextByte)), nil
}

func imm16(r *bufio.Reader, instruction *[]uint8) (string, error) {
    imm := make([]uint8, 2)
    _, err := io.ReadFull(r, imm)
    if err != nil {
        if err == io.EOF || err == io.ErrUnexpectedEOF {
            return "", &Z80AsmError{errorType: Z80AsmErrorMalformedInstruction}
        } else {
            return "", &Z80AsmError{errorType: Z80AsmErrorUnknown}
        }
    }
    *instruction = append(*instruction, imm[0], imm[1])
    return fmt.Sprintf("0x%02x%02x", imm[1], imm[0]), nil
}

func imm16_addr(r *bufio.Reader, instruction *[]uint8) (string, error) {
    imm := make([]uint8, 2)
    _, err := io.ReadFull(r, imm)
    if err != nil {
        if err == io.EOF || err == io.ErrUnexpectedEOF {
            return "", &Z80AsmError{errorType: Z80AsmErrorMalformedInstruction}
        } else {
            return "", &Z80AsmError{errorType: Z80AsmErrorUnknown}
        }
    }
    *instruction = append(*instruction, imm[0], imm[1])
    return fmt.Sprintf("[0x%02x%02x]", imm[1], imm[0]), nil
}

func r16_af_addr(r *bufio.Reader, instruction *[]uint8) string {
    reg_index := ((*instruction)[0] & 0x30) >> 4
    return fmt.Sprintf("[%s]", r16_af[reg_index])
}

func r16_sp_addr(r *bufio.Reader, instruction *[]uint8) string {
    reg_index := ((*instruction)[0] & 0x30) >> 4
    return fmt.Sprintf("[%s]", r16_sp[reg_index])
}

func decodeDJNZ(r *bufio.Reader, instruction *[]uint8, mnemonic *[]string) error {
    *mnemonic = append(*mnemonic, "djnz")
    /* Read operand (next byte) */
    operand, err := imm8_s(r, instruction)
    if err != nil {
        return err
    }
    *mnemonic = append(*mnemonic, operand)
    return nil
}

func decodeJR_E(r *bufio.Reader, instruction *[]uint8, mnemonic *[]string) error {
    *mnemonic = append(*mnemonic, "jr")
    /* Read operand (next byte) */
    operand, err := imm8(r, instruction)
    if err != nil {
        return err
    }
    *mnemonic = append(*mnemonic, operand)
    return nil
}

func decodeJR_cond_E(r *bufio.Reader, instruction *[]uint8, mnemonic *[]string) error {
    *mnemonic = append(*mnemonic, "jr")
    cond_index := ((*instruction)[0] & 0x38) >> 3 - 4
    *mnemonic = append(*mnemonic, conditions[cond_index])
    operand, err := imm8(r, instruction)
    if err != nil {
        return err
    }
    *mnemonic = append(*mnemonic, operand)
    return nil

}

func decodeLD_r16_nn(r *bufio.Reader, instruction *[]uint8, mnemonic *[]string) error {
    reg_index := ((*instruction)[0] & 0x30) >> 4
    *mnemonic = append(*mnemonic, "ld")
    *mnemonic = append(*mnemonic, r16_sp[reg_index])
    operand, err := imm16(r, instruction)
    if err != nil {
        return err
    }
    *mnemonic = append(*mnemonic, operand)
    return nil
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

func decodeLD_nn_HL(r *bufio.Reader, instruction *[]uint8, mnemonic *[]string) error {
    *mnemonic = append(*mnemonic, "ld")
    operand, err := imm16_addr(r, instruction)
    if err != nil {
        return err
    }
    *mnemonic = append(*mnemonic, operand)
    *mnemonic = append(*mnemonic, "hl")
    return nil
}

func decodeLD_nn_A(r *bufio.Reader, instruction *[]uint8, mnemonic *[]string) error {
    *mnemonic = append(*mnemonic, "ld")
    operand, err := imm16_addr(r, instruction)
    if err != nil {
        return err
    }
    *mnemonic = append(*mnemonic, operand)
    *mnemonic = append(*mnemonic, "a")
    return nil
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

func decodeLD_HL_nn(r *bufio.Reader, instruction *[]uint8, mnemonic *[]string) error {
    *mnemonic = append(*mnemonic, "ld")
    *mnemonic = append(*mnemonic, "hl")
    operand, err := imm16_addr(r, instruction)
    if err != nil {
        return err
    }
    *mnemonic = append(*mnemonic, operand)
    return nil
}

func decodeLD_A_nn(r *bufio.Reader, instruction *[]uint8, mnemonic *[]string) error {
    *mnemonic = append(*mnemonic, "ld")
    *mnemonic = append(*mnemonic, "a")
    operand, err := imm16_addr(r, instruction)
    if err != nil {
        return err
    }
    *mnemonic = append(*mnemonic, operand)
    return nil
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

func decodeLD_r8_n(r *bufio.Reader, instruction *[]uint8, mnemonic *[]string) error {
    reg_index := ((*instruction)[0] & 0x38) >> 3
    *mnemonic = append(*mnemonic, "ld")
    *mnemonic = append(*mnemonic, r8[reg_index])
    operand, err := imm8(r, instruction)
    if err != nil {
        return err
    }
    *mnemonic = append(*mnemonic, operand)
    return nil
}

func decodeLD_r8_r8(r *bufio.Reader, instruction *[]uint8, mnemonic *[]string) {
    reg_dst := ((*instruction)[0] & 0x38) >> 3
    reg_src := (*instruction)[0] & 0x7
    *mnemonic = append(*mnemonic, "ld")
    *mnemonic = append(*mnemonic, r8[reg_dst])
    *mnemonic = append(*mnemonic, r8[reg_src])
}

func decodeALU_r8(r *bufio.Reader, instruction *[]uint8, mnemonic *[]string) {
    alu_op := ((*instruction)[0] & 0x38) >> 3
    reg_index := (*instruction)[0] & 0x07
    *mnemonic = append(*mnemonic, ALU[alu_op]...)
    *mnemonic = append(*mnemonic, r8[reg_index])
}

func decodeRET_cc(r *bufio.Reader, instruction *[]uint8, mnemonic *[]string) {
    cc := ((*instruction)[0] & 0x38) >> 3
    *mnemonic = append(*mnemonic, "ret")
    *mnemonic = append(*mnemonic, conditions[cc])
}

func decodePOP_r16(r *bufio.Reader, instruction *[]uint8, mnemonic *[]string) {
    reg_index := ((*instruction)[0] & 0x30) >> 4
    *mnemonic = append(*mnemonic, "pop")
    *mnemonic = append(*mnemonic, r16_af[reg_index])
}

func decodeJP_HL(r *bufio.Reader, instruction *[]uint8, mnemonic *[]string) {
    *mnemonic = append(*mnemonic, "jp")
    *mnemonic = append(*mnemonic, "hl")
}

func decodeLD_SP_HL(r *bufio.Reader, instruction *[]uint8, mnemonic *[]string) {
    *mnemonic = append(*mnemonic, "ld")
    *mnemonic = append(*mnemonic, "sp")
    *mnemonic = append(*mnemonic, "hl")
}

func decodeJP_cc_nn(r *bufio.Reader, instruction *[]uint8, mnemonic *[]string) error {
    cc := ((*instruction)[0] & 0x38) >> 3
    *mnemonic = append(*mnemonic, "jp")
    *mnemonic = append(*mnemonic, conditions[cc])
    operand, err := imm16(r, instruction)
    if err != nil {
        return err
    }
    *mnemonic = append(*mnemonic, operand)
    return nil
}

func decodeJP_nn(r *bufio.Reader, instruction *[]uint8, mnemonic *[]string) error {
    *mnemonic = append(*mnemonic, "jp")
    operand, err := imm16(r, instruction)
    if err != nil {
        return err
    }
    *mnemonic = append(*mnemonic, operand)
    return nil
}

func decodeOUT_n_A(r *bufio.Reader, instruction *[]uint8, mnemonic *[]string) error {
    *mnemonic = append(*mnemonic, "out")
    operand, err := imm8(r, instruction)
    if err != nil {
        return err
    }
    *mnemonic = append(*mnemonic, fmt.Sprintf("[%s]", operand))
    *mnemonic = append(*mnemonic, "a")
    return nil
}

func decodeIN_a_n(r *bufio.Reader, instruction *[]uint8, mnemonic *[]string) error {
    *mnemonic = append(*mnemonic, "in")
    *mnemonic = append(*mnemonic, "a")
    operand, err := imm8(r, instruction)
    if err != nil {
        return err
    }
    *mnemonic = append(*mnemonic, fmt.Sprintf("[%s]", operand))
    return nil
}

func decodeEX_SP_HL(r *bufio.Reader, instruction *[]uint8, mnemonic *[]string) {
    *mnemonic = append(*mnemonic, "ex")
    *mnemonic = append(*mnemonic, "[sp]")
    *mnemonic = append(*mnemonic, "hl")
}

func decodeEX_DE_HL(r *bufio.Reader, instruction *[]uint8, mnemonic *[]string) {
    *mnemonic = append(*mnemonic, "ex")
    *mnemonic = append(*mnemonic, "de")
    *mnemonic = append(*mnemonic, "hl")
}

func decodeCALL_cc_nn(r *bufio.Reader, instruction *[]uint8, mnemonic *[]string) error {
    cc := ((*instruction)[0] & 0x38) >> 3
    *mnemonic = append(*mnemonic, "call")
    *mnemonic = append(*mnemonic, conditions[cc])
    operand, err := imm16(r, instruction)
    if err != nil {
        return err
    }
    *mnemonic = append(*mnemonic, operand)
    return nil
}

func decodePUSH_r16(r *bufio.Reader, instruction *[]uint8, mnemonic *[]string) {
    reg_index := ((*instruction)[0] & 0x30) >> 4
    *mnemonic = append(*mnemonic, "push")
    *mnemonic = append(*mnemonic, r16_af[reg_index])
}

func decodeCALL_nn(r *bufio.Reader, instruction *[]uint8, mnemonic *[]string) error {
    *mnemonic = append(*mnemonic, "call")
    operand, err := imm16(r, instruction)
    if err != nil {
        return err
    }
    *mnemonic = append(*mnemonic, operand)
    return nil
}

func decodeALU_n(r *bufio.Reader, instruction *[]uint8, mnemonic *[]string) error {
    alu_op := ((*instruction)[0] & 0x38) >> 3
    *mnemonic = append(*mnemonic, ALU[alu_op]...)
    operand, err := imm8(r, instruction)
    if err != nil {
        return err
    }
    *mnemonic = append(*mnemonic, operand)
    return nil
}

func decodeRST(r *bufio.Reader, instruction *[]uint8, mnemonic *[]string) {
    t := ((*instruction)[0] & 0x38) >> 3
    *mnemonic = append(*mnemonic, "rst")
    *mnemonic = append(*mnemonic, fmt.Sprintf("0x%02x", t * 8))
}

/*
 * Bumps the pointer in r
 * returns: the instruction bytes, the instruction mnemonic as an array of tokens
 */
func decodeInstruction(r *bufio.Reader) ([]uint8, []string, error) {
    /* If EOF, return empty string */
    var instruction []uint8
    nextByte, err := r.ReadByte()
    if err != nil {
        if err == io.EOF {
            return nil, nil, nil
        }
    }

    instruction = append(instruction, nextByte)
    var mnemonic []string;

    switch nextByte {
    /* prefix */
    case 0xcb:
        err = &Z80AsmError{errorType: Z80AsmErrorUnimplementedInstruction}
    case 0xdd:
        err = &Z80AsmError{errorType: Z80AsmErrorUnimplementedInstruction}
    case 0xed:
        err = &Z80AsmError{errorType: Z80AsmErrorUnimplementedInstruction}
    case 0xfd:
        err = &Z80AsmError{errorType: Z80AsmErrorUnimplementedInstruction}
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
                    err = decodeDJNZ(r, &instruction, &mnemonic)
                case 0x18:
                    /*
                     * jr E - jump to PC + E
                     */
                    err = decodeJR_E(r, &instruction, &mnemonic)
                default:
                    /* jr nz|z|nc|c, E*/
                    err = decodeJR_cond_E(r, &instruction, &mnemonic)
                }
            case 0x01:
                /* switch on bit 3 */
                switch nextByte & 0x08 {
                case 0x00:
                    /* ld rp[p], nn */
                    err = decodeLD_r16_nn(r, &instruction, &mnemonic)
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
                            err = decodeLD_nn_HL(r, &instruction, &mnemonic)
                        case 0x30:
                            /* ld [nn], a */
                            err = decodeLD_nn_A(r, &instruction, &mnemonic)
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
                        err = decodeLD_HL_nn(r, &instruction, &mnemonic)
                    case 0x30:
                        /* ld a, [nn] */
                        err = decodeLD_A_nn(r, &instruction, &mnemonic)
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
                err = decodeLD_r8_n(r, &instruction, &mnemonic)
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
            switch nextByte & 0x07{
            case 0x6:
                switch nextByte & 0x38 {
                case 0x30:
                    mnemonic = append(mnemonic, "halt")
                default:
                    /* ld r, r' */
                    decodeLD_r8_r8(r, &instruction, &mnemonic)
                }
            default:
                /* ld r, r' */
                decodeLD_r8_r8(r, &instruction, &mnemonic)
            }
        case 0x80:
            /* assorted ALU instructions */
            decodeALU_r8(r, &instruction, &mnemonic)
        case 0xc0:
            switch nextByte & 0x07 {
                case 0x00:
                    /* ret CC */
                    decodeRET_cc(r, &instruction, &mnemonic)
                case 0x01:
                    switch nextByte & 0x08 {
                    case 0x00:
                        decodePOP_r16(r, &instruction, &mnemonic)
                    case 0x08:
                        switch nextByte & 0x30 {
                        case 0x00:
                            mnemonic = append(mnemonic, "ret")
                        case 0x10:
                            mnemonic = append(mnemonic, "exx")
                        case 0x20:
                            decodeJP_HL(r, &instruction, &mnemonic)

                        case 0x30:
                            decodeLD_SP_HL(r, &instruction, &mnemonic)
                        }
                    }
                case 0x02:
                    err = decodeJP_cc_nn(r, &instruction, &mnemonic)
                case 0x03:
                    switch nextByte & 0x38 {
                    case 0x00:
                        err = decodeJP_nn(r, &instruction, &mnemonic)
                    case 0x08:
                        /* 0xcb prefix */
                    case 0x10:
                        err = decodeOUT_n_A(r, &instruction, &mnemonic)
                    case 0x18:
                        err = decodeIN_a_n(r, &instruction, &mnemonic)
                    case 0x20:
                        decodeEX_SP_HL(r, &instruction, &mnemonic)
                    case 0x28:
                        decodeEX_DE_HL(r, &instruction, &mnemonic)
                    case 0x30:
                        mnemonic = append(mnemonic, "di")
                    case 0x38:
                        mnemonic = append(mnemonic, "ei")
                    }
                case 0x04:
                    err = decodeCALL_cc_nn(r, &instruction, &mnemonic)
                case 0x05:
                    switch nextByte & 0x08 {
                    case 0x00:
                        decodePUSH_r16(r, &instruction, &mnemonic)
                    case 0x08:
                        switch nextByte & 0x30 {
                        case 0x00:
                            err = decodeCALL_nn(r, &instruction, &mnemonic)
                        case 0x10:
                            /* DD prefix */
                        case 0x20:
                            /* ED prefix */
                        case 0x30:
                            /* FD prefix */
                        }
                    }
                case 0x06:
                    err = decodeALU_n(r, &instruction, &mnemonic)
                case 0x07:
                    decodeRST(r, &instruction, &mnemonic)
            }
        }
    }
    return instruction, mnemonic, err
}

func disassemblerLoop(r *bufio.Reader) int {
    var addr uint32 = 0x0
    for instruction, mnemonic, err := decodeInstruction(r);
        len(instruction) != 0;
        instruction, mnemonic, err = decodeInstruction(r) {
        /* Generate hex encoding of instruction */
        instructionHex := make([]uint8, hex.EncodedLen(len(instruction)))
        hex.Encode(instructionHex, instruction)

        if err != nil {
            fmt.Printf("%s\n", err.Error())
            fmt.Printf("0x%016x: %-12s\n", addr, instructionHex)
            return 1
        }

        /* format - addr: <instruction bytes> <instruction mnemonic> */
        operands := ""
        if len(mnemonic) > 1 {
            operands = strings.Join(mnemonic[1:], ", ")
        }
        fmt.Printf("0x%016x: %-12s %-6s %s\n", addr, instructionHex, mnemonic[0], operands)
        addr += uint32(len(instruction))
    }
    return 0
}
