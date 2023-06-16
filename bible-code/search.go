package main

import (
	"fmt"
	"math/rand"
	"regexp"
	"sort"
)

// "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// https://pi.math.cornell.edu/~mec/2003-2004/cryptography/subs/frequencies.html
// https://www3.nd.edu/~busiforc/handouts/cryptography/letterfrequencies.html
var characters = string("EEEEEEEEEEEAAAAAAAARRRRRRRIIIIIIIOOOOOOOTTTTTTNNNNNN" +
	"SSSSSLLLLLCCCCUUUDDDPPPMMMHHHGGBBFYWKVXZJQ" +
	"000111222333444555666777888999")

func buildText(length int) []byte {
	b := []byte{}
	for i := 1; i < length; i++ {
		b = append(b, characters[rand.Intn(len(characters))])
	}
	return b
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
	return fmt.Sprintf("\"%s\": %d-%d %d+%d", d.name, d.start, d.end, d.skip, d.offset)
}

func main() {
	// keys := map[string]bool{
	// 	"to": true,
	// }
	text := buildText(1000000)
	// text := []byte("abcd2C0O2V0IzD")
	// text := []byte("GEMYCOMMANDMENTSMYSTATUTESANDMYLAWSANDISAACDWELTINGERARANDTHEMENOFTHEPLACEASKEDHIMOFHISWIFEANDHESAIDSHEISMYSISTERFORHEFEAREDTOSAYSHEISMYWIFELESTSAIDHETHEMENOFTHEPLACESHOULDKILLMEFORREBEKAHBECAUSESHEWASFAIRTOLOOKUPONANDITCAMETOPASSWHENHEHADBEENTHEREALONGTIMETHATABIMELECHKINGOFTHEPHILISTINESLOOKEDOUTATAWINDOWANDSAWANDBEHOLDISAACWASSPORTINGWITHREBEKAHHISWIFEANDABIMELECHCALLEDISAACANDSAIDBEHOLDOFASURETYSHEISTHYWIFEANDHOWSAIDSTTHOUSHEISMYSISTERANDISAACSAIDUNTOHIMBECAUSEISAIDLESTIDIEFORHERANDABIMELECHSAIDWHATISTHISTHOUHASTDON")
	// key := regexp.MustCompile("(2020)|(COVID)|(VIRUS)|(SICK)|(ILLNESS)|(PANDEMIC)|(CORONA)|(CHINA)")
	// key := regexp.MustCompile("(BIBLE)|(CODE)|(EDOC)|(ELBIB)")
	// key := regexp.MustCompile("(BIBLE)|(ELBIB)")
	key := regexp.MustCompile("(ILLNESS)|(SSENLLI)|(PANDEMIC)|(CIMEDNAP)|(CORONA)|(ANOROC)")

	discoveryList := []Discovery{}
	numDiscoveries := 0

	for skip := 1; skip < 300; skip++ {
		for offset := 0; offset < skip; offset++ {
			var code []byte
			for index, char := range text {
				if index%skip == offset {
					code = append(code, char)
				}
			}
			s := string(code)
			//fmt.Println(s)
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
					discoveryList = append(discoveryList, d)
					fmt.Printf("%v\n", d)
					numDiscoveries++
				}
				// fmt.Printf("%d/%d %d  %v\n", skip, length/4, offset, l)
			}
		}
		// fmt.Printf("skip: %d num: %d\n", skip, numDiscoveries)
	}

	sort.SliceStable(discoveryList, func(i, j int) bool {
		return discoveryList[i].start < discoveryList[j].start
	})

	for _, d := range discoveryList {
		fmt.Printf("%v\n", d)
	}
	// fmt.Println(string(text))
}
