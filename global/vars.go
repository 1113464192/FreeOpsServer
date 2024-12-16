package global

import (
	"errors"
	"golang.org/x/sync/semaphore"
	"sync"
)

var (
	RootPath      string
	Sem           *semaphore.Weighted
	MaxWebSSH     uint64
	WebSSHCounter uint64
	mu            sync.Mutex
	OpsSSHKey     []byte
)

func IncreaseWebSSHConn() error {
	mu.Lock()
	defer mu.Unlock()
	if WebSSHCounter >= MaxWebSSH {
		return errors.New("已达到最大webSSH数量")
	}

	WebSSHCounter++

	return nil
}

func ReduceWebSSHConn() {
	mu.Lock()
	defer mu.Unlock()

	if WebSSHCounter > 0 {
		WebSSHCounter--
	}
}
