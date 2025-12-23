package main

import "sync"

// RWLockFair is an RWLock that ensures fair ordering between readers and writers.
// It prevents writer starvation by serializing arrival to the lock,
// ensuring that once a writer appears, no new readers can overtake it.
type RWLockFair struct {
	*RWLock
	orderLock sync.Mutex
}

func NewRWLockFair() *RWLockFair { return &RWLockFair{RWLock: NewRWLock()} }

func (f *RWLockFair) AcquireRead() {
	f.orderLock.Lock()
	f.RWLock.AcquireRead()
	f.orderLock.Unlock()
}

func (f *RWLockFair) AcquireWrite() {
	f.orderLock.Lock()
	f.RWLock.AcquireWrite()
	f.orderLock.Unlock()
}
