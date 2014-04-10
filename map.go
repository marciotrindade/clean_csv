package main

import "fmt"

func main() {
	emails := make(map[string]bool)

	emails["marciotrindade@gmail.com"] = true

	if emails["marciotrindade@gmail.com"] {
		fmt.Println("emails exists")
	} else {
		fmt.Println("emails doesn't exist")
	}

	fmt.Println(emails["marciotrindade@gmail.com"])
}
