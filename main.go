package main

import (
	"fmt"
	"time"
)

const TransferLimit = 10000

type User struct {
	ID uint
	Name string
	Accounts []Account
}

type Account struct {
	ID string
	Balance float64
}

type Transaction struct {
	ID uint
	Amount float64
	Origin *Account
	Destiny *Account
	Date time.Time
}

// Here we save the amount of the deposits & withdrawals for each Account
var historical map[string][]float64

func main (){
	historical = make(map[string][]float64)
	user1,user2 := generateUsers()
	fmt.Println("User 1 -> ",user1)
	fmt.Println("User 2 -> ",user2)
	fmt.Println()

	// Print current balance of all the accounts from user 1 & user 2
	fmt.Printf("User 1 Account %v: %v\n", user1.Accounts[0].ID, GetBalance(user1.Accounts[0]))
	fmt.Printf("User 1 Account %v: %v\n", user1.Accounts[1].ID, GetBalance(user1.Accounts[1]))
	fmt.Printf("User 2 Account %v: %v\n", user2.Accounts[0].ID, GetBalance(user1.Accounts[0]))
	fmt.Printf("User 2 Account %v: %v\n", user2.Accounts[1].ID, GetBalance(user1.Accounts[1]))
	fmt.Println()

	// Transfer $5,000 from User 1 Account 1 to User 1 Account 2
	tx := Transaction{
		ID: uint(time.Now().Unix()),
		Amount: float64(5000),
		Origin: &user1.Accounts[0],
		Destiny: &user1.Accounts[1],
		Date: time.Now(),
	}
	if err := MakeTransaction(tx); err != nil {
		fmt.Println("ERROR: ", err)
	}
	fmt.Println()

	// Print current balance of User 1 Account 1 & User 1 Account 2
	fmt.Printf("User 1 Account %v: %v\n", user1.Accounts[0].ID, GetBalance(user1.Accounts[0]))
	fmt.Printf("User 1 Account %v: %v\n", user1.Accounts[1].ID, GetBalance(user1.Accounts[1]))
	fmt.Println()

	// Transfer $15,000 from User 1 Account 2 to User 2 Account 1 (can use multiple transactions)
	txAmount := 15_000
	remainder := txAmount % TransferLimit
	oneWhole := txAmount - remainder
	for _, amount := range []int{oneWhole, remainder} {
		tx := Transaction{
			ID:      uint(time.Now().Unix()),
			Amount:  float64(amount),
			Origin:  &user1.Accounts[1],
			Destiny: &user2.Accounts[0],
			Date:    time.Now(),
		}
		if err := MakeTransaction(tx); err != nil {
			fmt.Println("ERROR: ", err)
		}
	}
	fmt.Println()

	// Print current balance of User 1 Account 2 & User 2 Account 1
	fmt.Printf("User 1 Account %v: %v\n", user1.Accounts[1].ID, GetBalance(user1.Accounts[1]))
	fmt.Printf("User 2 Account %v: %v\n", user2.Accounts[0].ID, GetBalance(user2.Accounts[0]))
	fmt.Println()

	// Print historic deposits & withdrawals for each user
	GetHistoricalTransactions(user1)
	GetHistoricalTransactions(user2)
	fmt.Println()

	// Try to transfer $20,000 from User 1 Account 1 to User 2 Account 1
	tx = Transaction{
		ID: uint(time.Now().Unix()),
		Amount: 20_000.0,
		Origin: &user1.Accounts[0],
		Destiny: &user2.Accounts[0],
		Date: time.Now(),
	}
	if err := MakeTransaction(tx); err != nil {
		fmt.Println("ERROR: ", err)
	}
	fmt.Println()

	// Try to transfer $11,000 from User 2 Account 1 to User 1 Account 1
	tx = Transaction{
		ID: uint(time.Now().Unix()),
		Amount: 11_000.0,
		Origin: &user2.Accounts[0],
		Destiny: &user1.Accounts[0],
		Date: time.Now(),
	}
	if err := MakeTransaction(tx); err != nil {
		fmt.Println("ERROR: ", err)
	}
	fmt.Println()
}

func GetBalance(account Account) float64 {
	return account.Balance
}

func GetHistoricalTransactions(user User){
	fmt.Printf("Transaction History for User %v:\n", user.ID)
	for _, account := range user.Accounts {
		fmt.Printf("Account %v: %v\n", account.ID, historical[account.ID])
	}
}

func MakeTransaction(transaction Transaction) error{
	//get amount
	amount := transaction.Amount
	var errMsg string
	if transaction.Origin.Balance < amount {
		errMsg = "low balance"
	} else if amount > TransferLimit {
		errMsg = "exceeds transfer limit"
	}

	if errMsg == "" {
		//subtract from origin
		transaction.Origin.Balance -= amount
		//credit destination
		transaction.Destiny.Balance += amount

		//include transcation detail in tx hist
		originHist := historical[transaction.Origin.ID]
		historical[transaction.Origin.ID] = append(originHist, -amount)
		destHist := historical[transaction.Destiny.ID]
		historical[transaction.Destiny.ID] = append(destHist, amount)
	} else {
		return fmt.Errorf(errMsg)
	}
	return nil
}

func generateUsers() (user1,user2 User){
	// USER 1
	user1 = User{
		ID:       1,
		Name:     "User #1",
	}
	account1User1 := Account{
		ID:      "1111-1111-1111-1111",
		Balance: 20000,
	}
	account2User1 := Account{
		ID:     "2222-2222-2222-2222",
		Balance: 20000,
	}
	user1.Accounts = append(user1.Accounts,account1User1)
	user1.Accounts = append(user1.Accounts,account2User1)
	// USER 2
	user2 = User{
		ID:       2,
		Name:     "User #2",
	}
	account1User2 := Account{
		ID:      "3333-3333-3333-3333",
		Balance: 20000,
	}
	account2User2 := Account{
		ID:     "4444-4444-4444-4444",
		Balance: 20000,
	}
	user2.Accounts = append(user2.Accounts,account1User2)
	user2.Accounts = append(user2.Accounts,account2User2)
	return
}