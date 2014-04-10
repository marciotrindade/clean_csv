package main

import "fmt"
import "regexp"

func main() {
	regex := `^([^@\s]+)@((?:[-a-z0-9]+\.)+[a-z]{2,})$`
	emails := [3]string{
		"marciotrindade@gmail.com",
		"marciotrindade@gmail",
		"marciotrindade@gmail.com.br",
	}

	for _, email := range emails {
		match, _ := regexp.MatchString(regex, email)
		fmt.Println(match)
	}
}
