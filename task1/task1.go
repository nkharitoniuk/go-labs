package main

import ("fmt"
		"os"
		"strconv"
)

func main() {

	var err error
	
	if len(os.Args) < 2{
		panic("Please enter a number")
	}
	argument := os.Args[1]	
	var number int
	
	if number, err = strconv.Atoi(argument); err != nil {
        panic(err)
    }
    
	rank := 10
	for number/rank > 0 {
		rank = rank * 10
	}
	if (number * number) % rank == number {
		fmt.Printf("%v is automorphic number", number)
	} else {
		fmt.Printf("%v is not automorphic number", number)
	}
}