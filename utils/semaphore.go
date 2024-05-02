package utils

type Semaphore struct {
	bc chan struct{}
}

func NewSemaphore(n int) *Semaphore {
	return &Semaphore{bc: make(chan struct{}, n)}
}

func (s *Semaphore) Acquire() {
	s.bc <- struct{}{}
}

func (s *Semaphore) Release() {
	<-s.bc
}
