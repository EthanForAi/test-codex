package main

import (
	"fmt"
	"math/rand"
	"time"
)

func printRandomNumber() {
	// Seed with current time to get different values across runs.
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	now := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("[%s] helloworld random: %d\n", now, r.Intn(100))
}

func main() {
	// Continuously print a random number every second.
	for {
		printRandomNumber()
		time.Sleep(time.Second)
	}
}
