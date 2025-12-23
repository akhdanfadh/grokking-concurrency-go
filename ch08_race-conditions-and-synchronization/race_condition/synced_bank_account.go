package main

import (
	"sync"
)

type SyncedBankAccount struct {
	inner *UnsyncedBankAccount
	mu    sync.Mutex
}

func NewSyncedBankAccount(balance float32) *SyncedBankAccount {
	return &SyncedBankAccount{inner: NewUnsyncedBankAccount(balance)}
}

func (a *SyncedBankAccount) Deposit(amount float32) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.inner.Deposit(amount)
}

func (a *SyncedBankAccount) Withdraw(amount float32) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.inner.Withdraw(amount)
}

func (a *SyncedBankAccount) Balance() float32 {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.inner.Balance()
}
