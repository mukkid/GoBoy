package main

import (
    "bufio"
    "bytes"
    "encoding/hex"
    "fmt"
)


type reg8_id uint

const (
    B reg8_id = 0
    C
    D
    E
    H
    L
    /* 6 unused? */
    A = 0x07
)

/* Consumes an immediate 8 bit value from the stream, updates the args buffer with it */
func imm8(r *bufio.Reader, instruction *[]uint8, args *bytes.Buffer) {
    nextByte, err := r.ReadByte()
    if err != nil {
        return
    }
    *instruction = append(*instruction, nextByte)
    args.WriteString("0x")
    arg := []uint8{nextByte}
    offset := make([]uint8, hex.EncodedLen(len(arg)))
    hex.Encode(offset, arg)
    args.Write(offset)
}

func decodeDJNZ(r *bufio.Reader, instruction *[]uint8, mnemonic *bytes.Buffer, args *bytes.Buffer) {
    mnemonic.WriteString("djnz")
    /* Read operand (next byte) */
    imm8(r, instruction, args)
}

func decodeJR_E(r *bufio.Reader, instruction *[]uint8, mnemonic *bytes.Buffer, args *bytes.Buffer) {
    mnemonic.WriteString("jr")
    /* Read operand (next byte) */
    imm8(r, instruction, args)
}

func decodeJR_nz_E(r *bufio.Reader, instruction *[]uint8, mnemonic *bytes.Buffer, args *bytes.Buffer) {
    mnemonic.WriteString("jr")
    args.WriteString("nz, ")
    imm8(r, instruction, args)
}

func decodeJR_z_E(r *bufio.Reader, instruction *[]uint8, mnemonic *bytes.Buffer, args *bytes.Buffer) {
    mnemonic.WriteString("jr")
    args.WriteString("z, ")
    imm8(r, instruction, args)
}

func decodeJR_nc_E(r *bufio.Reader, instruction *[]uint8, mnemonic *bytes.Buffer, args *bytes.Buffer) {
    mnemonic.WriteString("jr")
    args.WriteString("nc, ")
    imm8(r, instruction, args)
}

func decodeJR_c_E(r *bufio.Reader, instruction *[]uint8, mnemonic *bytes.Buffer, args *bytes.Buffer) {
    mnemonic.WriteString("jr")
    args.WriteString("c, ")
    imm8(r, instruction, args)
}
/*
 * Bumps the pointer in r
 * returns: the instruction bytes, the instruction mnemonic, the instruction operands
 */
func decodeInstruction(r *bufio.Reader) ([]uint8, string, string) {
    /* If EOF, return empty string */
    var instruction []uint8
    nextByte, err := r.ReadByte()
    if err != nil {
        return nil, "", ""
    }

    instruction = append(instruction, nextByte)
    var mnemonic bytes.Buffer;
    var args bytes.Buffer

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
                        mnemonic.WriteString("nop")
                    case 0x08:
                        /* ex af,af' */
                        /* Not implemented in GB */
                        mnemonic.WriteString("ex")
                        args.WriteString("af, af'")
                    case 0x10:
                        /*
                         * djnz x
                         * Not implemented in GB
                         */
                         decodeDJNZ(r, &instruction, &mnemonic, &args)
                    case 0x18:
                        /*
                         * jr E - jump to PC + E
                         */
                         decodeJR_E(r, &instruction, &mnemonic, &args)
                    case 0x20:
                        /* jr nz, E*/
                        decodeJR_nz_E(r, &instruction, &mnemonic, &args)
                    case 0x28:
                        /* jr z, E*/
                        decodeJR_z_E(r, &instruction, &mnemonic, &args)
                    case 0x30:
                        /* jr nc, E*/
                        decodeJR_nc_E(r, &instruction, &mnemonic, &args)
                    case 0x38:
                        /* jr c, E*/
                        decodeJR_c_E(r, &instruction, &mnemonic, &args)
                    }
                case 0x01:
                case 0x02:
                case 0x03:
                case 0x04:
                case 0x05:
                case 0x06:
                case 0x07:
                }
            case 0x40:
            case 0x80:
            case 0xc0:
        }
    }
    return instruction, mnemonic.String(), args.String()
}

func disassemblerLoop(r *bufio.Reader) {
    var addr uint32 = 0x0
    for instruction, mnemonic, operands := decodeInstruction(r);
        len(instruction) != 0;
        instruction, mnemonic, operands = decodeInstruction(r) {
        /* Generate hex encoding of instruction */
        instructionHex := make([]uint8, hex.EncodedLen(len(instruction)))
        hex.Encode(instructionHex, instruction)

        /* format - addr: <instruction bytes> <instruction mnemonic> */
        fmt.Printf("0x%016x: %-12s %-6s %s\n", addr, instructionHex, mnemonic, operands)
        addr += uint32(len(instruction))
    }
}
