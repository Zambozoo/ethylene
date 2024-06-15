package bytecode

import "sync"

type Bytecode interface{}

type Bytecodes struct {
	mu sync.Mutex
}

func (b *Bytecodes) Add(bytecodes *Bytecodes) {
	b.mu.Lock()
	defer b.mu.Unlock()
	panic("implement me")
}

func (b *Bytecodes) Size() int64 {
	panic("implement me")
}
