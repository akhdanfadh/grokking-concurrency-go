package main

import "sync"

type RWLocker interface {
	AcquireRead()
	ReleaseRead()
	AcquireWrite()
	ReleaseWrite()
}

// RWLock is a lock that allows multiple readers or one writer at a time.
type RWLock struct {
	readers   int
	readLock  sync.Mutex
	writeLock sync.Mutex
}

func NewRWLock() *RWLock { return &RWLock{} }

// AcquireRead acquires the read lock for the current thread.
// If there is a writer waiting for the lock, the method blocks until the writer releases the lock.
func (l *RWLock) AcquireRead() {
	l.readLock.Lock()
	l.readers++
	if l.readers == 1 {
		l.writeLock.Lock()
	}
	l.readLock.Unlock()
}

// ReleaseRead releases the read lock held by the current thread.
// If there are no more readers holding the lock, the method releases the write lock.
func (l *RWLock) ReleaseRead() {
	if l.readers < 1 {
		panic("ReleaseRead called without a matching AcquireRead")
	}

	l.readLock.Lock()
	l.readers--
	if l.readers == 0 {
		l.writeLock.Unlock()
	}
	l.readLock.Unlock()
}

// AcquireWrite acquires the write lock for the current thread.
// If there is a reader or a writer holding the lock, the method blocks until the lock is released.
func (l *RWLock) AcquireWrite() {
	l.writeLock.Lock()
}

// ReleaseWrite releases the write lock held by the current thread.
func (l *RWLock) ReleaseWrite() {
	l.writeLock.Unlock()
}
