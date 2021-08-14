package main
import (
    "fmt"
    "strings"
    "bufio"
    "os"
    "strconv"
    "log"
)

func check(e error) { if e != nil {panic(e)} }

func SplitOn(r rune) bool { return r == ',' || r == ' ' || r == '\t' } // delimiters to split on

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

    scanner := bufio.NewScanner(file)
    scanner.Split(bufio.ScanLines)
    var code []string
    var instruction uint32
    var address uint32 = 0
    lineCounter := 1

    // set up write file for machine code comparison
    f, err := os.Create("asm-tests/asm-u-bin/sub-ns-mc-u.txt")
    if err != nil { log.Fatal(err) }
    defer f.Close()

    // first pass
    for scanner.Scan() {
        //fmt.Println(scanner.Text())
        line := strings.Split(scanner.Text(), "#")[0] // get any text before the comment "#" and ignore any text after it
        //fmt.Println(line)
        code = strings.FieldsFunc(line, SplitOn) // split into n strings 
        if len(code) == 0 { continue }
        switch(code[0]) { // code[0] is operation
        // code[1] is rd, code[2] is rs1 code[3] is rs2 
        case "addi":
            if len(code) != 4 { fmt.Println("Missing argument on line: ", lineCounter) }
            //if code[4]
            imm, _ := strconv.Atoi(code[3])
            instruction = uint32(imm)<<20 | regBin[code[2]]<<15 | regBin[code[1]]<<7 | opBin[code[0]]

        case "add":
            if len(code) != 4 { fmt.Println("Missing argument on line: ", lineCounter) }
            instruction = regBin[code[3]]<<20 | regBin[code[2]]<<15 | regBin[code[1]]<<7 | opBin[code[0]]

        case "sub":
            if len(code) != 4 { fmt.Println("Missing argument on line: ", lineCounter) }
            instruction = regBin[code[3]]<<20 | regBin[code[2]]<<15 | regBin[code[1]]<<7 | opBin[code[0]]

        case "sll":
            if len(code) != 4 { fmt.Println("Missing argument on line: ", lineCounter) }
            instruction = regBin[code[3]]<<20 | regBin[code[2]]<<15 | regBin[code[1]]<<7 | opBin[code[0]]

        case "slt":
            if len(code) != 4 { fmt.Println("Missing argument on line: ", lineCounter) }
            instruction = regBin[code[3]] << 20 | regBin[code[2]]<<15 | regBin[code[1]]<<7 | opBin[code[0]]

        case "sltu":
            if len(code) != 4 { fmt.Println("Missing argument on line: ", lineCounter) }
            instruction = regBin[code[3]]<<20 | regBin[code[2]]<<15 | regBin[code[1]]<<7 | opBin[code[0]]

        case "xor":
            if len(code) != 4 { fmt.Println("Missing argument on line: ", lineCounter) }
            instruction = regBin[code[3]]<<20 | regBin[code[2]]<<15 | regBin[code[1]]<<7 | opBin[code[0]]

        case "srl":
            if len(code) != 4 { fmt.Println("Missing argument on line: ", lineCounter) }
            instruction = regBin[code[3]]<<20 | regBin[code[2]]<<15 | regBin[code[1]]<<7 | opBin[code[0]]

        case "sra":
            if len(code) != 4 { fmt.Println("Missing argument on line: ", lineCounter) }
            instruction = regBin[code[3]]<<20 | regBin[code[2]]<<15 | regBin[code[1]]<<7 | opBin[code[0]]

        case "or":
            if len(code) != 4 { fmt.Println("Missing argument on line: ", lineCounter) }
            instruction = regBin[code[3]]<<20 | regBin[code[2]]<<15 | regBin[code[1]]<<7 | opBin[code[0]]

        case "and":
            if len(code) != 4 { fmt.Println("Missing argument on line: ", lineCounter) }
            instruction = regBin[code[3]]<<20 | regBin[code[2]]<<15 | regBin[code[1]]<<7 | opBin[code[0]]

        default:
            fmt.Println("Syntax Error on line: ", lineCounter)
        }

        lineCounter++
        address += 4
        fmt.Printf("Address: 0x%08x   Instruction:  0x%08x\n", address, instruction)

        // write machine code to file for comparisons
        f.WriteString(fmt.Sprintf("0x%08x\n", instruction))
    }

    file.Close()
}
