package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"

	//lint:ignore ST1001 Ignoring dot-imports for the example-dir usage
	. "github.com/j-flat/go-veikkaus/goveikkaus"
)

func main() {
	r := bufio.NewReader(os.Stdin)
	fmt.Println("Veikkaus username: ")
	username, _ := r.ReadString('\n')

	fmt.Println("Veikkaus password: ")
	bytePassword, _ := term.ReadPassword(int(os.Stdin.Fd()))
	password := string(bytePassword)

	ctx := context.Background()

	client := NewClient(nil)

	_, _, err := client.Auth.Login(ctx, strings.TrimSpace(username), strings.TrimSpace(password))

	if err != nil {
		fmt.Printf("\nerror: %v\n", err)
		return
	}

	fmt.Println("Login successful")

	if err := client.Auth.Logout(); err != nil {
		fmt.Printf("Logout was not successful. Error: %s", err)
	} else {
		fmt.Println("Logout succesfully")
	}
}
