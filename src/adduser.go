package main

import (
	"code.google.com/p/go.crypto/bcrypt"
	"fmt"
	"github.com/howeyc/gopass"
	"log"
)

func main() {
	fmt.Printf("Username: ")
	username := gopass.GetPasswd()
	fmt.Printf("Password: ")
	password := gopass.GetPasswd()

	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Printf("%s:%s\n", username, hashedPassword)
}
