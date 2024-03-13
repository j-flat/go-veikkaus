package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/term"

	//lint:ignore ST1001 Ignoring dot-imports for the example-dir usage
	. "github.com/j-flat/go-veikkaus/goveikkaus"
)

func formatCurrency(amount int, currencySymbol string) string {
	// Convert the integer amount to a string
	amountStr := strconv.Itoa(amount)

	// Insert commas to separate thousands
	formattedAmount := insertCommas(amountStr)

	// Insert currency symbol
	formattedAmount += currencySymbol

	return formattedAmount
}

func insertCommas(s string) string {
	// Split the string into integer and fractional parts
	parts := strings.Split(s, "")
	integerPart := parts[:len(parts)-2]
	fractionalPart := parts[len(parts)-2:]

	// Insert commas to separate thousands
	for i := len(integerPart) - 3; i > 0; i -= 3 {
		integerPart = append(integerPart[:i], append([]string{","}, integerPart[i:]...)...)
	}

	// Combine integer and fractional parts with a comma separator
	return strings.Join(integerPart, "") + "," + strings.Join(fractionalPart, "")
}

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

	balance, _, err := client.Auth.AccountBalance(ctx)

	if err != nil {
		fmt.Printf("Could not return account-balance for the user. ERR: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Account status:", balance.Status)
	fmt.Printf("Account balance %s\n", formatCurrency(balance.Balances.Cash.Balance, balance.Balances.Cash.Currency))
	fmt.Println("Now logging out....")

	if err := client.Auth.Logout(); err != nil {
		fmt.Printf("Logout was not successful. Error: %s", err)
		os.Exit(1)
	} else {
		fmt.Println("Logout succesfully")
	}
}
