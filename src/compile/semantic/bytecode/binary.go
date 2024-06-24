package bytecode

// PrimitiveBinaryOp represents operations that pop two primitive values from the stack.
type PrimitiveBinaryOp struct {
	op op
}

func (pbo *PrimitiveBinaryOp) Bytes() []byte {
	return []byte{pbo.op.index}
}
