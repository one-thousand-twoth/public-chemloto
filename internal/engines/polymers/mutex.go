package polymers

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
	"sync"
)

type debugMutex struct {
	name string
	mux  sync.Mutex
}

func (d *debugMutex) Lock() {
	d.print("Lock")
	d.mux.Lock()
	d.print("Lock after")
}

func (d *debugMutex) Unlock() {
	d.mux.Unlock()
	d.print("Unlock")
}

func (d *debugMutex) print(l string) {
	_, fn, line, _ := runtime.Caller(2)
	fmt.Printf("Locker: %s:%d -> %v:%v -> %v\n", fn, line, d.name, l, getGID())
}

func getGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}
