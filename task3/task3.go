package main

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		panic("Please enter text")
	}
	enterString := os.Args[1]
	fmt.Printf("Start string: %s\n", enterString)
	reg, _ := regexp.Compile("[^а-яА-Яa-zA-Z0-9]+")
	resultString := reg.ReplaceAllString(enterString, " ")

	var slice = strings.Split(resultString, " ")
	mapString := CountWords(slice)

	SortMap(mapString, slice)
}

func CountWords(stringSlice []string) map[string]int {
	resultMap := make(map[string]int)
	for _, v := range stringSlice {
		resultMap[v]++
	}
	return resultMap
}

func SortMap(mapString map[string]int, sliceOrdered []string) {
	resultMap := map[int][]string{}
	var sortedSlice []int
	for k, v := range mapString {
		resultMap[v] = append(resultMap[v], k)
	}
	for k := range resultMap {
		sortedSlice = append(sortedSlice, k)
	}

	sort.Sort(sort.Reverse(sort.IntSlice(sortedSlice)))

	for _, intKey := range sortedSlice {

		stringSliceFromMap := resultMap[intKey]
		mapForInnerSort := make(map[int]string)
		var indexSliceForSort []int

		for _, stringV := range stringSliceFromMap {
			for index := range sliceOrdered {
				if sliceOrdered[index] == stringV {
					indexSliceForSort = append(indexSliceForSort, index)
					mapForInnerSort[index] = stringV
					break
				}
			}
		}
		sort.Ints(indexSliceForSort)

		for _, sortedIndex := range indexSliceForSort {
			fmt.Printf("%s(%d)\n", mapForInnerSort[sortedIndex], intKey)
		}
	}
}
