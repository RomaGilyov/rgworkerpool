package bpool

import (
	"testing"
)

type TestTask struct {
	C int
}

func (tt *TestTask) Execute() {
	tt.C++
}

func TestWorkerPool(t *testing.T) {
	pool := NewPool(10)

	tt := &TestTask{C: 0}

	pool.Exec(tt)
	pool.Exec(tt)

	for i := 0; i < 20; i++ {
		pool.Exec(tt)
	}

	pool.Close()

	pool.Wait()

	if tt.C != 22 {
		t.Fatal("Task must increment C until 22")
	}
}
