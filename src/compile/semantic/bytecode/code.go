package bytecode

import (
	"fmt"
	"math"
)

type op struct {
	index uint8
	name  string
}

func (c *op) String() string {
	return fmt.Sprintf("op{index:%d,name:%s}", c.index, c.name)
}

var ops = make([]op, 0, math.MaxUint8)

func newOp(name string) op {
	o := op{
		index: uint8(len(ops)),
		name:  name,
	}
	ops = append(ops, o)

	return o
}

var (
	OP_ARITH_ADD_INT      = newOp("OP_ARITH_ADD_INT")      // a+b
	OP_ARITH_SUBTRACT_INT = newOp("OP_ARITH_SUBTRACT_INT") // a-b
	OP_ARITH_MULTIPLY_INT = newOp("OP_ARITH_MULTIPLY_INT") // a*b
	OP_ARITH_DIVIDE_INT   = newOp("OP_ARITH_DIVIDE_INT")   // a/b
	OP_ARITH_MODULO_INT   = newOp("OP_ARITH_MODULO_INT")   // a%b
	OP_ARITH_NEGATE_INT   = newOp("OP_ARITH_NEGATE_INT")   // -a
	OP_ARITH_SHIFTRIGHT   = newOp("OP_ARITH_SHIFTRIGHT")   // a >>> b

	OP_ARITH_ADD_FLT      = newOp("OP_ARITH_ADD_FLT")      // a+b
	OP_ARITH_SUBTRACT_FLT = newOp("OP_ARITH_SUBTRACT_FLT") // a-b
	OP_ARITH_MULTIPLY_FLT = newOp("OP_ARITH_MULTIPLY_FLT") // a*b
	OP_ARITH_DIVIDE_FLT   = newOp("OP_ARITH_DIVIDE_FLT")   // a/b
	OP_ARITH_NEGATE_FLT   = newOp("OP_ARITH_NEGATE_FLT")   // -a

	OP_ARITH_CONCAT_STR = newOp("OP_ARITH_CONCAT_STR") // a + b

	OP_ARITH_CMP_INT = newOp("OP_ARITH_CMP_INT") // a<=>b
	OP_ARITH_CMP_FLT = newOp("OP_ARITH_CMP_FLT") // a<=>b
	OP_ARITH_CMP_STR = newOp("OP_ARITH_CMP_STR") // a<=>b
)

var (
	OP_LOGIC_AND = newOp("OP_LOGIC_AND") // a&&b
	OP_LOGIC_OR  = newOp("OP_LOGIC_OR")  // a||b
	OP_LOGIC_NOT = newOp("OP_LOGIC_NOT") // !a
)

var (
	OP_BIT_AND        = newOp("OP_BIT_AND")        // a&b
	OP_BIT_OR         = newOp("OP_BIT_OR")         // a|b
	OP_BIT_XOR        = newOp("OP_BIT_XOR")        // a^b
	OP_BIT_NOT        = newOp("OP_BIT_NOT")        // ~a
	OP_BIT_SHIFTRIGHT = newOp("OP_BIT_SHIFTRIGHT") // a >> b
	OP_BIT_SHIFTLEFT  = newOp("OP_BIT_SHIFTLEFT")  // a << b
)

var (
	OP_JUMP       = newOp("OP_JUMP")       // <-dest. Jumps to dest
	OP_JUMP_IFN   = newOp("OP_JUMP_IFN")   // <-top, <-dest. If !top, jumps to dest.
	OP_JUMP_STORE = newOp("OP_JUMP_STORE") // ->dest, <-top. Jumps to dest.
	OP_RETURN_N   = newOp("OP_RETURN_N")   // <-size, <-value, <-dest, <-scope, ->value. Pops return, jumps to dest, pushes return.
)
