package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

type BankAccount interface {
	Deposit(amount float32) error
	Withdraw(amount float32) error
	Balance() float32
}

type ATM struct {
	account BankAccount
}

func (a *ATM) transaction() {
	_ = a.account.Deposit(10)

	// simulate real action and encourage interleaving
	time.Sleep(1 * time.Millisecond)
	runtime.Gosched() // yield the processor to allow other goroutines to run

	_ = a.account.Withdraw(10)
}

func (a *ATM) Run() {
	a.transaction()
}

func testATMs(account BankAccount, atmNumber int) {
	var wg sync.WaitGroup
	wg.Add(atmNumber)
	for range atmNumber {
		atm := &ATM{account: account}
		go func() {
			defer wg.Done()
			atm.Run()
		}()
	}
	wg.Wait()
}

func main() {
	const atmNumber = 1000

	unsynced := NewUnsyncedBankAccount(0)
	testATMs(unsynced, atmNumber)
	fmt.Println("Balance of unsynced account after concurrent transactions:")
	fmt.Printf("Actual: %.0f\nExpected: 0\n\n", unsynced.Balance())

	synced := NewSyncedBankAccount(0)
	testATMs(synced, atmNumber)
	fmt.Println("Balance of synced account after concurrent transactions:")
	fmt.Printf("Actual: %.0f\nExpected: 0\n", synced.Balance())
}
