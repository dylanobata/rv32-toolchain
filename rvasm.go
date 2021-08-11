package main

import ( 
    "fmt"
    "bufio" 
    "os"
)

func check(e error) { if e != nil {panic(e)} }

func main() {
    if len(os.Args) != 2 {
        fmt.Println("Usage:", os.Args[0], "FILE.asm")
    }
    file, err := os.Open(os.Args[1])
    check(err)

    scanner := bufio.NewScanner(file)
    scanner.Split(bufio.ScanLines)
    var text []string

    for scanner.Scan() {
        fmt.Println(scanner.Text())
        text = append(text, scanner.Text())
    }
    file.Close()
    /*for _, each_ln := range text {
        fmt.Println(each_ln)
    }*/
}
