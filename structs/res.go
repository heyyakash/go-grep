package structs

import "sync"

type Result struct {
	Mutex sync.RWMutex
	Lines []string
}

func NewResultHolder() *Result {
	return &Result{
		Mutex: sync.RWMutex{},
		Lines: []string{},
	}
}

func (r *Result) AddLine(line string) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	r.Lines = append(r.Lines, line)
}

func (r *Result) GetLines() []string {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	return r.Lines
}
