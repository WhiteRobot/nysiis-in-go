package main

import (
	"os"
	"regexp"
	"strings"
)

func NYSIIS(name string, trueNYSIIS bool) string {
	
	/*
	Translate first characters of name: MAC → MCC; KN → NN; K → C; PH, PF → FF; SCH → SSS
    Translate last characters of name: EE, IE → Y; DT, RT, RD, NT, ND → D
    First character of key = first character of name.
    Translate remaining characters by following rules, incrementing by one character each time:
        EV → AF else A, E, I, O, U → A
        Q → G; Z → S; M → N
        KN → N else K → C
        SCH → SSS; PH → FF
        H → If previous or next is non-vowel, previous.
        W → If previous is vowel, A.
        Add current to key if current is not same as the last key character.
    If last character is S, remove it.
    If last characters are AY, replace with Y.
    If last character is A, remove it.
    Append translated key to value from step 3 (removed first character)
    If longer than 6 characters, truncate to first 6 characters. (only needed for true NYSIIS, some versions use the full key)
	*/
	
	translated := strings.ToUpper(name)
	
	// Translate first characters of name
	
	// MAC → MCC; SCH → SSS
	switch translated[0:3] {
		case "MAC":
			translated = "MCC" + translated[3:]
		case "SCH":
			translated = "SSS" + translated[3:]
	}
	
	// KN → NN; PH, PF → FF
	switch translated[0:2] {
		case "KN":
			translated = "NN" + translated[2:]
		case "PH", "PF":
			translated = "FF" + translated[2:]
	}
	
	// K → C
	if translated[0] == 'K' {
		translated = "C" + translated[1:]
	}
	
	// Translate last characters of name: EE, IE → Y; DT, RT, RD, NT, ND → D
	switch translated[len(translated)-2:] {
		case "EE", "IE":
			translated = translated[0:len(translated)-2] + "Y"
		case "DT", "RT", "RD", "NT", "ND":
			translated = translated[0:len(translated)-2] + "D"
	}
	
	lastCharacter := string(translated[0])
	translated = translated[1:]
	currentCharacter := ""
	append := ""
	translatedLength := len(translated)
	
	// Append translated key to value from step 3 (removed first character)
	key := lastCharacter
	
	vowels := regexp.MustCompile("[AEIOU]")
	
	// Translate remaining characters by following rules, incrementing by one character each time
	for i:=0; i < translatedLength; i++ { 
		
		currentCharacter = string(translated[i])
		
		if i < (translatedLength - 2) && translated[i:i+2] == "EV" {
			// EV → AF
			append = "AF"
		} else if vowels.MatchString(currentCharacter) {
			// A, E, I, O, U → A
			append = "A"
		} else if currentCharacter == "Q" {
			// Q → G
			append = "G"
		} else if currentCharacter == "Z" {
			// Z → S
			append = "S"
		} else if currentCharacter == "M" {
			// M → N
			append = "N"
		} else if i < (translatedLength - 2) && translated[i:i+2] == "KN" {
			// KN → N
			append = "N"
		} else if currentCharacter == "K" {
			// K → C
			append = "C"
		} else if i < (translatedLength - 3) && translated[i:i+3] == "SCH" {
			// SCH → SSS
			append = "SSS"
		} else if i < (translatedLength - 2) && translated[i:i+2] == "PH" {
			// PH → FF
			append = "FF"
		} else if currentCharacter == "H" && (vowels.MatchString(lastCharacter) == false || vowels.MatchString(currentCharacter) == false) {
			// H → If previous or next is non-vowel, previous.
			append = lastCharacter
		} else if currentCharacter == "W" && vowels.MatchString(lastCharacter) {
			// W → If previous is vowel, A.
			append = "A"
		} else {
			append = currentCharacter
		}
		
		if key[len(key) - 1:] != append {
			// Add current to key if current is not same as the last key character.
			key += append
		}
		
		lastCharacter = currentCharacter
		
	}// end while
	
	if key[len(key)-1] == 'S' {
		// If last character is S, remove it.
		key = key[0:len(key)-1]
	}else if key[len(key)-2:] == "AY" {
		// If last characters are AY, replace with Y.
		key = key[0:len(key)-2] + "Y"
	} else if key[len(key)-1] == 'A' {
		// If last character is A, remove it.
		key = key[0:len(key)-1]
	}
	
	// If longer than 6 characters, truncate to first 6 characters. (only needed for true NYSIIS, some versions use the full key)
	if trueNYSIIS && len(key) > 6 {
		return key[0:6]
	} else {
		return key
	}
}

func main() {
	name := os.Args[1]
	
	println(NYSIIS(name, true))
}
