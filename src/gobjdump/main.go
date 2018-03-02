package main

import (
    "os"
    "bufio"
    "fmt"
)

func main_c(argv []string) int {
    if len(argv) < 2 {
        fmt.Printf("Usage: %s <binary>\n", argv[0])
        return 1
    }

    /* Open binary */
    file, err := os.Open(argv[1])
    if err != nil {
        fmt.Printf("%s\n", err.Error())
        return 1
    }

    bufReader := bufio.NewReader(file)
    return disassemblerLoop(bufReader)
}

func main() {
    os.Exit(main_c(os.Args))
}
