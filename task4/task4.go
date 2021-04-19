package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

func main() {
	if len(os.Args) < 4 {
		panic("Please enter text")
	}
	originalString := os.Args[1]
	enterString := os.Args[2]
	keyString := os.Args[3]
	fmt.Printf("Start text: %s\n", originalString)
	fmt.Printf("Encrypted text: %s\n", enterString)
	fmt.Printf("Key words: %s\n", keyString)
	alfaSlice := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

	shift := findShift(enterString, keyString, alfaSlice)

	var resultString []string
	upperEnterString := strings.ToUpper(enterString)

	for _, runeChar := range upperEnterString {
		c := string(runeChar)
		decryptedCharIndex := -1
		decryptedCharIndex = getAlfaIndex(alfaSlice, c)

		if decryptedCharIndex >= 0 {
			resultString = append(resultString, alfaSlice[reCircularIndex(decryptedCharIndex, shift)])
		} else {
			resultString = append(resultString, c)
		}
	}

	fmt.Printf("Dencrypted text: %s\n", resultString)
}

func circularIndex(alfaRealIndex int, shift int) int {
	if (alfaRealIndex + shift) > 25 {
		return alfaRealIndex + shift - 26
	} else {
		return alfaRealIndex + shift
	}
}

func reCircularIndex(alfaRealIndex int, shift int) int {
	if (alfaRealIndex - shift) < 0 {
		return alfaRealIndex - shift + 26
	} else {
		return alfaRealIndex - shift
	}
}

func getAlfaIndex(alfaSlice []string, kChar string) int {
	for alfaIndex, alfaChar := range alfaSlice {
		if alfaChar == kChar {
			return alfaIndex
		}
	}
	return -1
}

func makeSliceFromString(text string) []string {
	text = strings.ToUpper(text)
	reg, _ := regexp.Compile("[^A-Z0-9]+")
	resultString := reg.ReplaceAllString(text, " ")
	return strings.Split(resultString, " ")
}

func findShift(enterString string, keyString string, alfaSlice []string) int {
	var textSlice = makeSliceFromString(enterString)
	var keyWordSlice = makeSliceFromString(keyString)

	var shift int

	equalLenWordsMap := make(map[string][]string)

	for _, itemKeyWord := range keyWordSlice {

		// iterate text encrypted words and collect map
		for _, itemText := range textSlice {
			if len(itemKeyWord) == len(itemText) {
				equalLenWordsMap[itemKeyWord] = append(equalLenWordsMap[itemKeyWord], itemText)
			}
		}
	}

	// iteration potential shifts
	for i := 1; i < 25; i++ {
		//k-key word, v-encrypted words slice, map iteration by key words
		for k, v := range equalLenWordsMap {
			var encryptedWordIndex int
			//iterate chars of key word
			for index, kChar := range k {

				var alfaRealIndex int
				var alfaShiftedIndex int

				alfaRealIndex = getAlfaIndex(alfaSlice, string(kChar))
				alfaShiftedIndex = circularIndex(alfaRealIndex, i)

				// if index of potential encrypted word defined
				if encryptedWordIndex != 0 {
					runes := []rune(v[encryptedWordIndex])

					vChar := string(runes[index])
					if vChar == alfaSlice[alfaShiftedIndex] {
						shift = i
						break
					}
				} else {
					//iterate all words with the same length
					for vIndex, vWord := range v {
						runes := []rune(vWord)

						vCharSlice := string(runes[index])
						if vCharSlice == alfaSlice[alfaShiftedIndex] {
							encryptedWordIndex = vIndex
							break
						}
					}
				}
			}
		}
	}

	fmt.Printf("Key is: %d\n", shift)
	return shift
}
