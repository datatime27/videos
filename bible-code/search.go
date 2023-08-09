package main

import (
	"fmt"
	"math/rand"
	"regexp"
    "os"
	"sort"
    "strings"
    "time"
)

// https://pi.math.cornell.edu/~mec/2003-2004/cryptography/subs/frequencies.html
// https://www3.nd.edu/~busiforc/handouts/cryptography/letterfrequencies.html
var dictionary = map[byte]int {
    'E': 56,
    'A': 43,
    'R': 38,
    'I': 38,
    'O': 36,
    'T': 35,
    'N': 33,
    'S': 29,
    'L': 27,
    'C': 23,
    'U': 18,
    'D': 17,
    'P': 16,
    'M': 15,
    'H': 15,
    'G': 12,
    'B': 10,
    'F': 9,
    'Y': 9,
    'W': 6,
    'K': 5,
    'V': 5,
    'X': 1,
    'Z': 1,
    'J': 1,
    'Q': 1,
    }
	
const (
    NUM_WORKERS = 8
)

func buildText(length int) []byte {
    alphabet := []byte{}
    for letter, quantity := range dictionary{
        for i := 0; i < quantity; i++ {
            alphabet = append(alphabet, letter)
        }
    }
    
	b := []byte{}
	for i := 0; i < length; i++ {
		b = append(b, alphabet[rand.Intn(len(alphabet))])
	}
	return b
}

func readFile(filename string) []byte {
    dat, err := os.ReadFile(filename)
    if err != nil {
        panic(err)
    }
    nonWords := regexp.MustCompile(`[^A-Za-z]`)
  
    text := string(dat)    
    return []byte(strings.ToUpper(nonWords.ReplaceAllString(text, "")))

}

func find(s string, l []string) bool {
	for _, i := range l {
		if i == s {
			return true
		}
	}
	return false
}

type Discovery struct {
	name   string
	start  int
	end    int
	skip   int
	offset int
}

func (d Discovery) String() string {
	return fmt.Sprintf("\"%s\": %d-%d skip:%d ", d.name, d.start, d.end, d.skip)
}

func (d Discovery) findIntersections(text string, l []Discovery, maxIntersections int, requiredKeys map[string]bool) int {
    intersections := []Discovery{d}
    minIndex := d.start-100
    maxIndex := d.end+100
    uniqueKeys := map[string]bool{d.name:true}
    
    for _, entry := range l {
        if d.end < entry.start {
            continue
        }
        if d.start > entry.end {
            continue
        }
        
        if entry.start-100 < minIndex{
            minIndex = entry.start-100
        }
        if entry.end+100 > maxIndex{
            maxIndex = entry.end+100
        }
        intersections = append(intersections, entry)
        uniqueKeys[entry.name] = true
    }
    
    if len(requiredKeys) > 0 {
        for key := range requiredKeys {
            if !uniqueKeys[key] && !uniqueKeys[reverse(key)]{
                return maxIntersections
            }
        }
    }
    
    if minIndex < 0 {
        minIndex = 0
    }
    if maxIndex >= len(text) {
        maxIndex = len(text) -1
    }
    
    length := maxIndex-minIndex
    if len(intersections) >= maxIntersections && length > 2000 && length < 4000 {
        fmt.Printf("INTERSECTION: %d-%d length: %d %v \n", minIndex, maxIndex, length, intersections)
        fmt.Println(text[minIndex:maxIndex])
        fmt.Println()
        maxIntersections = len(uniqueKeys)
    }
    return maxIntersections
}

func worker(key *regexp.Regexp, text []byte, id int, skips <-chan int, results chan<- []Discovery) {
    var codeTime int64
    var searchTime int64
    for skip := range skips {
    	discoveryList := []Discovery{}
		for offset := 0; offset < skip; offset++ {
            startCodeTime := time.Now().UnixNano()
			var code []byte
            for i:=offset; i < len(text); i+=skip {
                code = append(code,text[i])
            }
			s := string(code)
			//fmt.Println(s)
            endCodeTime := time.Now().UnixNano()
			if matches := key.FindAllString(s, -1); len(matches) > 0 {
				indexes := key.FindAllStringIndex(s, -1)
				// l := []int{}
				for i, index := range indexes {
					start := index[0]*skip + offset
					end := (index[1]-1)*skip + offset
					// fmt.Printf("%d %d %d %d-%d\n", skip, offset, index[1],
					// 	index[0]*skip+offset, (index[1]-1)*skip+offset)
					match := matches[i]

					d := Discovery{
						name:   match,
						start:  start,
						end:    end,
						skip:   skip,
						offset: offset,
					}
                    discoveryList = append(discoveryList,d)
				}
			}
            endSearchTime := time.Now().UnixNano()
            
            codeTime += endCodeTime-startCodeTime
            searchTime += endSearchTime-endCodeTime
		}
        results <- discoveryList
        
        //if skip % 800 == 0 {
        //    fmt.Printf("codeTime(M): %g searchTime(M) %g\n",float32(codeTime)/1000000.0, float32(searchTime)/1000000.0)
        //}
    }
}

func combos(textLen, wordLen int) int{
    total := 0
    for i := 1; i < textLen/wordLen; i++ {
        total += (textLen/i)-wordLen+1
    }
    return total
}

func getKey(words []string) *regexp.Regexp {
    key := ""
    for _, word := range words {
        key = fmt.Sprintf("%s|(%s)|(%s)", key, word, reverse(word))
    }
    return regexp.MustCompile(strings.Trim(key,"|"))
}

func reverse(s string) string {
    rns := []rune(s) // convert to rune
    for i, j := 0, len(rns)-1; i < j; i, j = i+1, j-1 {
        rns[i], rns[j] = rns[j], rns[i]
    }
    return string(rns)
}

func generateKeys(requiredKeys, optionalKeys []string) (*regexp.Regexp, map[string]bool) {
    sb := []string {}
    requiredKeyMap := map[string]bool{}
    
    for _, key := range requiredKeys {
        sb = append(sb, fmt.Sprintf("(%s)",key))
        sb = append(sb, fmt.Sprintf("(%s)",reverse(key)))
        requiredKeyMap[key] = true
    }
    for _, key := range optionalKeys {
        sb = append(sb, fmt.Sprintf("(%s)",key))
        sb = append(sb, fmt.Sprintf("(%s)",reverse(key)))
    }
    return regexp.MustCompile(strings.Join(sb,"|")), requiredKeyMap
}

func main() {

    
    // ASSASSIN in Moby Dick
    text := readFile("md.txt")
    text = []byte(regexp.MustCompile(`VESTHESNOW.+?PENSABLETOTH`).FindStringSubmatch(string(text))[0])
    pattern, requiredKeys := generateKeys([]string{
        "RABIN", "CHEST", "SHOT", "SHOTDEAD", "DEAD", "IGAL", "AMIR", "EYAL", "LEHI","THREE"}, []string{})
/*
    // BIBLE CODE in KJV
	text := []byte("GEMYCOMMANDMENTSMYSTATUTESANDMYLAWSANDISAACDWELTINGERARANDTHEMENOFTHEPLACEASKEDHIMOFHISWIFEANDHESAIDSHEISMYSISTERFORHEFEAREDTOSAYSHEISMYWIFELESTSAIDHETHEMENOFTHEPLACESHOULDKILLMEFORREBEKAHBECAUSESHEWASFAIRTOLOOKUPONANDITCAMETOPASSWHENHEHADBEENTHEREALONGTIMETHATABIMELECHKINGOFTHEPHILISTINESLOOKEDOUTATAWINDOWANDSAWANDBEHOLDISAACWASSPORTINGWITHREBEKAHHISWIFEANDABIMELECHCALLEDISAACANDSAIDBEHOLDOFASURETYSHEISTHYWIFEANDHOWSAIDSTTHOUSHEISMYSISTERANDISAACSAIDUNTOHIMBECAUSEISAIDLESTIDIEFORHERANDABIMELECHSAIDWHATISTHISTHOUHASTDON")
    pattern, requiredKeys := generateKeys([]string{"BIBLE","CODE"}, []string{})



    // COVID in Moby Dick
    text := readFile("md.txt")
    pattern, requiredKeys := generateKeys([]string{"WUHAN","COVID", "CHINA", "SICK", "VIRUS"}, []string{"ILLNESS", "PANDEMIC", "CORONA"})

    // Data Time Patreon in Moby Dick
    text := readFile("md.txt")
    pattern, requiredKeys := generateKeys([]string{"PATREON"}, []string{"DATA", "TIME", "SUBS"})

    // COVID in random letters
	text := buildText(1000 * 1000)
    pattern, requiredKeys := generateKeys([]string{"WUHAN","COVID", "CHINA", "SICK", "VIRUS"}, []string{"ILLNESS", "PANDEMIC", "CORONA"})

    // ASSASSIN in random letters
	text := buildText(1000 * 1000)
    pattern, requiredKeys := generateKeys([]string{"RABIN", "CHEST", "SHOT", "DEAD", "IGAL", "AMIR", "EYAL", "LEHI"}, []string{"ASSASSIN"})
    */



    
	//key := regexp.MustCompile("(BIBLE)|(ELBIB)")
	//key := regexp.MustCompile("(ILLNESS)|(PANDEMIC)|(CORONA)|(VIRUS)")
	//key := regexp.MustCompile("(DATA.{0,4}TIME)")
    //key := regexp.MustCompile("(YITZHAK)|(KAHZTIY)|(ASSASSINATE)|(ETANISSASSA)|(ASSASSIN)|(NISSASSA)") // |(RABIN)|(NIBAR)
    //key := getKey([]string{"CHEST", "SHOT", "RABIN", "DEAD", "SHOTDEAD", "IGAL", "AMIR", "EYAL", "LEHI"})
    
	discoveryList := []Discovery{}
    
    skipMin := 1
    skipMax := len(text)-5
    skips := make(chan int, skipMax)
    results := make(chan []Discovery, skipMax)
    
    fmt.Printf("combinations %d\n",combos(len(text),7))

    for w := 1; w <= NUM_WORKERS; w++ {
        go worker(pattern, text, w, skips, results)
    }
    fmt.Printf("%s: Searching up to skip code %d\n", time.Now().Format(time.UnixDate), skipMax)
	for skip := skipMin; skip < skipMax; skip++ {
        skips <- skip
		// fmt.Printf("skip: %d num: %d\n", skip, numDiscoveries)
	}
    close(skips)
    
    textString := string(text)

    maxIntersections := 0
    for skip := skipMin; skip < skipMax; skip++ {
        if skip % 1000 == 0{
            fmt.Printf("%s: Current Skip Code: %d\n", time.Now().Format(time.UnixDate), skip)
        }
        for _, d := range <-results {
            //fmt.Printf("%v\n", d)
            maxIntersections = d.findIntersections(textString, discoveryList, maxIntersections, requiredKeys)
            discoveryList = append(discoveryList,d)
            }

    }

	sort.SliceStable(discoveryList, func(i, j int) bool {
		return discoveryList[i].start < discoveryList[j].start
	})

    fmt.Printf("\nOrdered Results:\n")
	for _, d := range discoveryList {
		fmt.Printf("%v\n", d)
	}
    fmt.Printf("%s: Done\n", time.Now().Format(time.UnixDate))

}