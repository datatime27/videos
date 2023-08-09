package main

import (
	"fmt"
)


func main() {
    textLen := 1000* 1000
    wordLen := 7
    
    total := 0
    skipMax := textLen/wordLen+1
    total = textLen*skipMax
    /*
    for skip := 1; skip <= skipMax; skip++ {
    total += textLen
    
		for offset := 0; offset < skip; offset++ {
            codeLen := (textLen)/skip
            //fmt.Printf("skip %d, offset %d, codeLen %d, wordLen %d\n", skip, offset, codeLen, wordLen)
            
            if codeLen < wordLen {
                continue
            }
            
            total += codeLen// - wordLen
        }
    
    }
    */
    fmt.Printf("%d letter combinations %d\n", wordLen, total)
}