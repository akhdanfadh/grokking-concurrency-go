package main

import "fmt"

type UnsyncedBankAccount struct {
	balance float32
}

func NewUnsyncedBankAccount(balance float32) *UnsyncedBankAccount {
	return &UnsyncedBankAccount{balance: balance}
}

func (a *UnsyncedBankAccount) Deposit(amount float32) error {
	if amount <= 0 {
		return fmt.Errorf("you can't deposit a negative amount of money")
	}
	a.balance += amount
	return nil
}

func (a *UnsyncedBankAccount) Withdraw(amount float32) error {
	if amount <= 0 || amount > a.balance {
		return fmt.Errorf("account does not have sufficient funds")
	}
	a.balance -= amount
	return nil
}

func (a *UnsyncedBankAccount) Balance() float32 {
	return a.balance
}
