package main

import (
	"bufio"
	"fmt"
	"github.com/AMythicDev/controlplane/pii-detection"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter text: ")
	if scanner.Scan() {
		input := scanner.Text()
		body, err := piidetection.RecognizePII(input)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}
		rstr, err := piidetection.RedactPII(body, input)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}
		fmt.Printf("%s\n", rstr)
	}
}
