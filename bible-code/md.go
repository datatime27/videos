package main
import (
    "fmt"
    "regexp"
    "os"
    "strings"
)


func main() {
    dat, err := os.ReadFile("md.txt")
    if err != nil {
        panic(err)
    }
    nonWords := regexp.MustCompile(`[^A-Za-z]`)
  
    fmt.Printf("Moby Dick Full Length: %d\n", len(dat))
    text := string(dat)    
    text = strings.ToUpper(nonWords.ReplaceAllString(text, ""))
    
    fmt.Printf("Moby Dick all letters: %d\n", len(text))

    find := regexp.MustCompile(`VESTHESNOW.+?PENSABLETOTH`)
    code := find.FindStringSubmatch(text)[0]
    skip := 8057
    width := 36
    
    for i := 0; i <= len(code); i+= skip {
        fmt.Printf("%v\n", code[i:i+width])
    }
    fmt.Printf("%d\n", len(code))

}