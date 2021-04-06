package main

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		panic("Please enter a string")
	}
	enterString := os.Args[1]
	fmt.Printf("Start string: %s\n", enterString)
	decompression(compression(enterString))
}

func compression(enterString string) string {
	resultString := ""
	for _, charRune := range enterString {
		char := string(charRune)
		count := strings.Count(enterString, char)
		if count > 4 {
			n := count
			for n > 4 {
				repeatedString := strings.Repeat(char, n)
				if strings.Contains(enterString, repeatedString) {
					resultString = strings.ReplaceAll(enterString, repeatedString, join(n, char))
					enterString = resultString
				}
				n -= 1
			}
		}
	}
	fmt.Printf("Compress result string: %s\n", resultString)
	return resultString
}

func decompression(enterString string) {
	resultString := ""
	re := regexp.MustCompile(`#.+?#.`)
	combinationSlice := re.FindAllString(enterString, -1)
	if combinationSlice != nil && len(combinationSlice) > 0 {
		for _, combinationString := range combinationSlice {
			regForCount := regexp.MustCompile(`#.+?#`)
			countString := regForCount.FindString(combinationString)
			symbolIndex := len(combinationString) - 1
			count, _ := strconv.Atoi(strings.ReplaceAll(countString, "#", ""))
			resultString = strings.Replace(enterString, combinationString, strings.Repeat(string(combinationString[symbolIndex]), count), -1)
			enterString = resultString
		}
	}
	fmt.Printf("Decompress result string: %s", resultString)
}

func join(n int, char string) string {
	var b bytes.Buffer
	b.WriteString("#")
	b.WriteString(strconv.Itoa(n))
	b.WriteString("#")
	b.WriteString(char)
	return b.String()
}
