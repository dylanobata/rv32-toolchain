package main

import (
    "fmt"
    "strings"
    "bufio"
    "io"
    "os"
    "strconv"
    "errors"
    "log"
)

func check(e error) { if e != nil {panic(e)} }

func SplitOn(r rune) bool { return r == ',' || r == ' ' || r == '\t' || r == '(' || r == ')'} // delimiters to split on

func isValidImmediate(s string) (int64, error) {
    var imm1, imm2, imm3 int64
    var err1, err2, err3 error
    imm1, err1 = strconv.ParseInt(s, 10, 32) // check if s is a decimal number

    if strings.HasPrefix(s, "0x") {
        imm2, err2 = strconv.ParseInt(s[2:], 16, 64) // check if s is hex
    } else if strings.HasPrefix(s, "-0x") {
        imm2, err2 = strconv.ParseInt(string(s[0]) + s[3:], 16, 64) // ignore the "0x" part but include the '-'
    } else if strings.HasPrefix(s, "0b") {
        imm3, err3 = strconv.ParseInt(s[2:], 2, 64) // check if s is binary
    } else if strings.HasPrefix(s, "-0b") {
        imm3, err3 = strconv.ParseInt(string(s[0]) + s[3:], 2, 64)
    }

    if err1 != nil && err2 != nil && err3 != nil {
        return 0, errors.New("Invalid immediate value")
    } else if err1 == nil {
        return imm1, nil
    } else if err2 == nil {
        return imm2, nil
    } else {
        return imm3, nil
    }
}

func main() {
    regBin := map[string] uint32{
        "x0" : 0b00000,  "zero" : 0b00000,
        "x1" : 0b00001,  "ra"   : 0b00001,
        "x2" : 0b00010,  "sp"   : 0b00010,
        "x3" : 0b00011,  "gp"   : 0b00011,
        "x4" : 0b00100,  "tp"   : 0b00100,
        "x5" : 0b00101,  "t0"   : 0b00101,
        "x6" : 0b00110,  "t1"   : 0b00110,
        "x7" : 0b00111,  "t2"   : 0b00111,
        "x8" : 0b01000,  "s0"   : 0b01000, "fp": 0b01000,
        "x9" : 0b01001,  "s1"   : 0b01001,
        "x10": 0b01010,  "a0"   : 0b01010,
        "x11": 0b01011,  "a1"   : 0b01011,
        "x12": 0b01100,  "a2"   : 0b01100,
        "x13": 0b01101,  "a3"   : 0b01101,
        "x14": 0b01110,  "a4"   : 0b01110,
        "x15": 0b01111,  "a5"   : 0b01111,
        "x16": 0b10000,  "a6"   : 0b10000,
        "x17": 0b10001,  "a7"   : 0b10001,
        "x18": 0b10010,  "s2"   : 0b10010,
        "x19": 0b10011,  "s3"   : 0b10011,
        "x20": 0b10100,  "s4"   : 0b10100,
        "x21": 0b10101,  "s5"   : 0b10101,
        "x22": 0b10110,  "s6"   : 0b10110,
        "x23": 0b10111,  "s7"   : 0b10111,
        "x24": 0b11000,  "s8"   : 0b11000,
        "x25": 0b11001,  "s9"   : 0b11001,
        "x26": 0b11010,  "s10"  : 0b11010,
        "x27": 0b11011,  "s11"  : 0b11011,
        "x28": 0b11100,  "t3"   : 0b11100,
        "x29": 0b11101,  "t4"   : 0b11101,
        "x30": 0b11110,  "t5"   : 0b11110,
        "x31": 0b11111,  "t6"   : 0b11111,
    }

    opBin := map[string] uint32 {
        "lui"   : 0b00000000000000000000000000110111,
        "auipc" : 0b00000000000000000000000000010111,
        "jal"   : 0b00000000000000000000000001101111,
        "jalr"  : 0b00000000000000000000000001100111,

        "beq"   : 0b00000000000000000000000001100011,
        "bne"   : 0b00000000000000000001000001100011,
        "blt"   : 0b00000000000000000100000001100011,
        "bge"   : 0b00000000000000000101000001100011,
        "bltu"  : 0b00000000000000000110000001100011,
        "bgeu"  : 0b00000000000000000111000001100011,

        "lb"    : 0b00000000000000000000000000000011,
        "lh"    : 0b00000000000000000001000000000011,
        "lw"    : 0b00000000000000000010000000000011,
        "lbu"   : 0b00000000000000000100000000000011,
        "lhu"   : 0b00000000000000000101000000000011,

        "sb"    : 0b00000000000000000000000000100011,
        "sh"    : 0b00000000000000000001000000100011,
        "sw"    : 0b00000000000000000010000000100011,

        "addi"  : 0b00000000000000000000000000010011,
        "slti"  : 0b00000000000000000010000000010011,
        "sltiu" : 0b00000000000000000011000000010011,
        "xori"  : 0b00000000000000000100000000010011,
        "ori"   : 0b00000000000000000110000000010011,
        "andi"  : 0b00000000000000000111000000010011,

        "slli"  : 0b00000000000000000001000000010011,
        "srli"  : 0b00000000000000000101000000010011,
        "srai"  : 0b01000000000000000101000000010011,

        "add"   : 0b00000000000000000000000000110011,
        "sub"   : 0b01000000000000000000000000110011,
        "sll"   : 0b00000000000000000001000000110011,
        "slt"   : 0b00000000000000000010000000110011,
        "sltu"  : 0b00000000000000000011000000110011,
        "xor"   : 0b00000000000000000100000000110011,
        "srl"   : 0b00000000000000000101000000110011,
        "sra"   : 0b01000000000000000101000000110011,
        "or"    : 0b00000000000000000110000000110011,
        "and"   : 0b00000000000000000111000000110011,

        "ecall" : 0b00000000000000000000000001110011,
        "ebreak": 0b00000010000000000000000001110011,

    }

    if len(os.Args) != 2 { fmt.Println("Usage:", os.Args[0], "FILE.s") }
    file, err := os.Open(os.Args[1])
    check(err)
    defer file.Close()

    scanner := bufio.NewScanner(file) // stores content from file
    scanner.Split(bufio.ScanLines)

    var code []string
    var instruction uint32
    var address uint32 = 0
    lineCounter := 1

    symbolTable := make(map[string]int64, 100)
    const UNKNOWN = -1
//    literalPool := make(map[string]int64, 100)
    for scanner.Scan() { // first pass
        line := strings.Split(scanner.Text(), "#")[0] // get any text before the comment "#" and ignore any text after it
        code = strings.FieldsFunc(line, SplitOn) // break code into its operation and operands 
        if len(code) == 0 { // filter out whitespace
            lineCounter++
            continue
        }
        switchOnOp := code[0] // check if first entry of code is a label or an op
        if strings.HasSuffix(switchOnOp, ":") {
            label := strings.TrimSuffix(code[0], ":")
            symbolTable[label] = int64(address) // if label exists in symbolTable, update value to valid address
            if len(code) >= 2 { // opcode is in code[1] if code[0] is a label
                switchOnOp = code[1]
                code = code[1:]
            } else { continue }
        }

        switch(switchOnOp) {
        case "lui", "auipc", "jal": // Instruction format:  op  rd, imm     or      label: op  rd, imm 
            if len(code) != 3 && len(code) != 4 {
                fmt.Println("Incorrect argument count on line: ", lineCounter)
                os.Exit(0)
            }
            if len(code) == 4 && !strings.HasSuffix(code[0], ":") && len(code[0]) > 1 {
                fmt.Printf("%s not a valid label\n", code[0])
                os.Exit(0)
            }
           if len(code) == 4 {
                label := strings.TrimSuffix(code[0], ":")
                _, exists := symbolTable[label]
                if exists {
                    symbolTable[label] = int64(address) // if label exists in symbolTable, update value to valid address
                }
            }

        case "beq", "bne", "blt", "bge", "bltu", "bgeu": // Instruction format:  op rs1, rs2, label     or     label: op    rs1,rs2, label
            if len(code) != 4 && len(code) != 5 {
                fmt.Println("Incorrect argument count on line: ", lineCounter)
                os.Exit(0)
            }
            if len(code) == 5 && !strings.HasSuffix(code[0], ":") && len(code[0]) > 1{
                fmt.Printf("%s not a valid label\n", code[0])
                os.Exit(0)
            }
            if len(code) == 4 {
                _, exists := symbolTable[code[3]]
                if !exists {
                    symbolTable[code[3]] = UNKNOWN // if symbol is not in symbolTable, create entry with flag -1
                }
            }
            if len(code) == 5 {
                label := strings.TrimSuffix(code[0], ":")
                _, exists := symbolTable[code[0]]
                if exists { // check if label exists in symbol table
                    symbolTable[label] = int64(address)
                }
                _, exists = symbolTable[code[4]]
                if !exists {
                    symbolTable[code[4]] = UNKNOWN
                }
            }

        case "lb", "lh", "lw", "lbu", "lhu": // Instruction format: op rd, imm(rs1)     or      label: op rd, imm(rs1)
            if len(code) != 4 && len(code) != 5 {
                fmt.Println("Incorrect argument count on line: ", lineCounter)
                os.Exit(0)
            }
            if len(code) == 5 && !strings.HasSuffix(code[0], ":") && len(code[0]) > 1 {
                fmt.Printf("%s not a valid label\n", code[0])
                os.Exit(0)
            }
            if len(code) == 5 {
               label := strings.TrimSuffix(code[0], ":")
                _, exists := symbolTable[code[0]]
                if exists {
                    symbolTable[label] = int64(address) // if label exists in symbolTable, update value to valid address
                }
            }

        case "sb", "sh", "sw": // Instruction format: op rs2, imm(rs1)      or      label: op rs2, imm(rs1) 
            if len(code) != 4 && len(code) != 5 {
                fmt.Println("Incorrect argument count on line: ", lineCounter)
                os.Exit(0)
            }
            if len(code) == 5 && !strings.HasSuffix(code[0], ":") && len(code[0]) > 1 {
                fmt.Printf("%s not a valid label\n", code[0])
                os.Exit(0)
            }
            if len(code) == 5 {
                label := strings.TrimSuffix(code[0], ":")
                _, exists := symbolTable[code[0]]
                if exists {
                    symbolTable[label] = int64(address) // if label exists in symbolTable, update value to valid address
                }
            }

        case "addi", "slti", "sltiu", "xori", "ori", "andi", "jalr": // Instruction format: op rd, rs1, imm     or      label:  op rd, rs1, imm
            if len(code) != 4 && len(code) != 5 {
                fmt.Println("Incorrect argument count on line: ", lineCounter)
                os.Exit(0)
            }
            if len(code) == 5 && !strings.HasSuffix(code[0], ":") && len(code[0]) > 1 {
                fmt.Printf("%s not a valid label\n", code[0])
                os.Exit(0)
            }
            if len(code) == 5 {
                label := strings.TrimSuffix(code[0], ":")
                _, exists := symbolTable[label]
                if exists {
                    symbolTable[label] = int64(address)
                }
            }

        case "slli", "srli", "srai": // Instruction format: op rd, rs1, imm     or      label: rd, rs1, imm
            if len(code) != 4 && len(code) != 5 {
                fmt.Println("Incorrect argument count on line: ", lineCounter)
                os.Exit(0)
            }
            if len(code) == 5 && !strings.HasSuffix(code[0], ":") && len(code[0]) > 1 {
                fmt.Printf("%s not a valid label\n", code[0])
                os.Exit(0)
            }
            if len(code) == 5 {
                label := strings.TrimSuffix(code[0], ":")
                _, exists := symbolTable[label]
                if exists {
                    symbolTable[label] = int64(address) // if label exists in symbolTable, update value to valid address
                }
            }

        case "add", "sub", "sll", "slt", "sltu", "xor", "srl", "sra", "or", "and": // Instruction format: op rd, rs1, rs2       or      label: op rd, rs1, rs2
            if len(code) != 4 && len(code) != 5 {
                fmt.Println("Incorrect argument count on line: ", lineCounter)
                os.Exit(0)
            }
            if len(code) == 5 && !strings.HasSuffix(code[0], ":") && len(code[0]) > 1 {
                fmt.Printf("%s not a valid label\n", code[0])
                os.Exit(0)
            }

        case "ecall", "ebreak": // Instruction format: op       or      label: op
            if len(code) != 1 && len(code) != 2 {
                fmt.Println("Incorrect argument count on line: ", lineCounter)
                os.Exit(0)
            }
            if len(code) == 2 && !strings.HasSuffix(code[0], ":") && len(code[0]) > 1 {
                fmt.Printf("%s not a valid label\n", code[0])
                os.Exit(0)
            }
            // check if imm has a constant definition 

        default:
            fmt.Println("Syntax Error on line: ", lineCounter)
            os.Exit(0)
        }
        lineCounter++
        address += 4

    }

    for key, element := range symbolTable {
        fmt.Println("Key:", key, "Element:", element)
    }
    // reset file to start and reinitialize scanner
    _, err = file.Seek(0, io.SeekStart)
    scanner = bufio.NewScanner(file)
    scanner.Split(bufio.ScanLines)

    // set up write file for machine code comparison
    f, err := os.Create("asm-tests/asm-u-bin/beq-bne-mc-u.txt")
    if err != nil { log.Fatal(err) }
    defer f.Close()

    address = 0
    lineCounter = 1
    for scanner.Scan() { // second pass
        line := strings.Split(scanner.Text(), "#")[0] // get any text before the comment "#" and ignore any text after it
        code = strings.FieldsFunc(line, SplitOn) // split into  operation, operands, and/or labels
        if len(code) == 0 { // code is whitespace, ignore it
            lineCounter++
            continue
        }
        switchOnOp := code[0] // check if first entry of code is a label or an op
        if strings.HasSuffix(switchOnOp, ":") {
            if len(code) >= 2 { // opcode is in code[1] if code[0] is a label
                switchOnOp = code[1]
                code = code[1:] // reindex array so that op is at index 0 
            } else {
                continue
            }
        }
        switch(switchOnOp) { // switch on operation 
        case "lui", "auipc":
            if len(code) != 3 { fmt.Println("Incorrect argument count on line: ", lineCounter) }
            imm, err := isValidImmediate(code[2])
            op, opFound := opBin[code[0]]
            rd, rdFound := regBin[code[1]]
            if err != nil {
                fmt.Printf("Error on line %d: %s\n", lineCounter, err)
                os.Exit(0)
            }
            if imm > 1048575 || imm < 0 {
                fmt.Printf("Error on line %d: Immediate value out of range (should be between 0 and 1048575)\n", lineCounter)
                os.Exit(0)
            }
            if !opFound || !rdFound {
                fmt.Println("Invalid register on line", lineCounter)
                os.Exit(0)
            }
            instruction = uint32(imm)<<12 | rd<<7 | op

        case "jal":
            if len(code) != 3 { fmt.Println("Incorrect argument count on line: ", lineCounter) }
            op, opFound := opBin[code[0]]
            rd, rdFound := regBin[code[1]]
            label, labelFound := symbolTable[code[2]]
            if !labelFound {
                fmt.Println("Error: label not found")
                os.Exit(0)
            }
            if !opFound && !rdFound {
               fmt.Println("Invalid register on line", lineCounter)
               os.Exit(0)
            }
            label = label - int64(address)
            instruction = (uint32(label) & 0x80000)<<11 | (uint32(label) & 0x7FE)<<20 | (uint32(label) & 0x400)<<19 | (uint32(label) & 0x7F800)<<11 | rd<<7 | op

        case "beq", "bne", "blt", "bge", "bltu", "bgeu": // op rs1, rs2, imm
            if len(code) != 4 { fmt.Println("Incorrect argument count on line: ", lineCounter) }
            label, labelFound := symbolTable[code[3]]
            if !labelFound {
                fmt.Printf("On line %d: label not found\n", lineCounter,)
                os.Exit(0)
            }
            label = label - int64(address)
            fmt.Println("Label:", label)
            op, opFound := opBin[code[0]]
            rs1, rs1Found := regBin[code[1]]
            rs2, rs2Found := regBin[code[2]]
            if  opFound && rs1Found && rs2Found {
                instruction =  (uint32(label) & 0x800)<<20 | ((uint32(label) & 0x7E0))<<20 | rs2<<20 | rs1<<15 | (uint32(label) & 0x1E)<<7 | (uint32(label) & 0x400)>>3 | op
            } else if !rs1Found || !rs2Found {
                    fmt.Println("Invalid register on line", lineCounter)
                    os.Exit(0)
            }

        case "lb", "lh", "lw", "lbu", "lhu": // op rd, imm(rs1)
            if len(code) != 4 {
                fmt.Println("Incorrect argument count on line: ", lineCounter)
            }
            imm, err := isValidImmediate(code[2])
            if err != nil {
                fmt.Printf("Error on line %d: %s\n", lineCounter, err)
                os.Exit(0)
            }
            op, opFound := opBin[code[0]]
            rd, rdFound := regBin[code[1]]
            rs1, rs1Found := regBin[code[3]]
            if  opFound && rdFound && rs1Found {
                instruction = uint32(imm)<<20 | rs1<<15 | rd<<7 | op
            } else {
                if !opFound {
                    fmt.Println("Invalid operation on line", lineCounter)
                    os.Exit(0)
                } else if !rdFound || !rs1Found {
                    fmt.Println("Invalid register on line", lineCounter)
                    os.Exit(0)
                }
            }

        case "sb", "sh", "sw": // op rs2, imm(rs1)
            if len(code) != 4 {
                fmt.Println("Incorrect argument count on line: ", lineCounter)
            }
            imm, err := isValidImmediate(code[2])
            if err != nil {
                fmt.Printf("Error on line %d: %s\n", lineCounter, err)
                os.Exit(0)
            }
            op, opFound := opBin[code[0]]
            rs2, rs2Found := regBin[code[1]]
            rs1, rs1Found := regBin[code[3]]
            if  opFound && rs1Found && rs2Found {
                instruction = (uint32(imm) & 0xFE0)<<20 | rs2<<20 | rs1<<15 | (uint32(imm) & 0x1F)<<7 | op
            } else if !rs1Found || !rs2Found {
                    fmt.Println("Invalid register on line", lineCounter)
                    os.Exit(0)
                }

        case "addi", "slti", "sltiu", "xori", "ori", "andi", "jalr": // op rd, rs1, immediate
            if len(code) != 4 { fmt.Println("Incorrect argument count on line: ", lineCounter) }
            imm, err := isValidImmediate(code[3])
            if err != nil {
                fmt.Printf("Error on line %d: %s\n", lineCounter, err)
                os.Exit(0)
            }
            if imm > 2047 || imm < -2048 {
                fmt.Printf("Error on line %d: Immediate value out of range (should be between -2048 and 2047)\n", lineCounter)
                os.Exit(0)
            }
            op, opFound := opBin[code[0]]
            rd, rdFound := regBin[code[1]]
            rs1, rs1Found := regBin[code[2]]
            if  opFound && rdFound && rs1Found {
                instruction = uint32(imm)<<20 | rs1<<15 | rd<<7 | op
            }  else if !rdFound || !rs1Found {
                    fmt.Println("Invalid register on line", lineCounter)
                    os.Exit(0)
            }

        case "slli", "srli", "srai": // op rd, rs1, immediate
            if len(code) != 4 {
                fmt.Println("Incorrect argument count on line: ", lineCounter)
                os.Exit(0)
            }
            imm, err := isValidImmediate(code[3])
            if err != nil {
                fmt.Printf("Error on line %d: %s\n", lineCounter, err)
                os.Exit(0)
            }
            if imm > 31 {
                fmt.Printf("Error on line %d: Immediate value out of range (should be between 0 and 31)")
                os.Exit(0)
            }
            op, opFound := opBin[code[0]]
            rd, rdFound := regBin[code[1]]
            rs1, rs1Found := regBin[code[2]]
            if  opFound && rdFound && rs1Found {
                instruction = uint32(imm)<<20 | rs1<<15 | rd<<7 | op
            }  else if !rdFound || !rs1Found {
                    fmt.Println("Invalid register on line", lineCounter)
                    os.Exit(0)
            }

        case "add", "sub", "sll", "slt", "sltu", "xor", "srl", "sra", "or", "and": // op rd, rs1, rs2
            if len(code) != 4 {
                fmt.Println("Incorrect argument count on line: ", lineCounter)
                os.Exit(0)
             }
            op, opFound := opBin[code[0]]
            rd, rdFound := regBin[code[1]]
            rs1, rs1Found := regBin[code[2]]
            rs2, rs2Found := regBin[code[3]]
            if  opFound && rdFound && rs1Found && rs2Found {
                instruction = rs2<<20 | rs1<<15 | rd<<7 | op // code[0]=op, code[1]=rd, code[2]=rs1 code[3]=rs2 
            }  else if !rdFound || !rs1Found || !rs2Found {
                    fmt.Println("Invalid register on line", lineCounter)
                    os.Exit(0)
            }

        case "ecall", "ebreak":
            if len(code) != 1 {
                fmt.Println("Too many arguments on line: ", lineCounter)
                os.Exit(0)
            }
            instruction = opBin[code[0]]

        default:
            fmt.Println("Syntax Error on line: ", lineCounter)
            os.Exit(0)
        }
        fmt.Printf("Address: 0x%08x     Line: %d     Instruction:  0x%08x\n", address, lineCounter, instruction)
        lineCounter++
        address += 4

        //write machine code to file for comparisons
        f.WriteString(fmt.Sprintf("0x%08x\n", instruction))
    }
}
