package typeid

import (
	"geth-cody/ast"
)

var (
	ID_Void   = &TypeID{index: 0}
	ID_Bool   = &TypeID{index: 1}
	ID_Int    = &TypeID{index: 2}
	ID_Float  = &TypeID{index: 3}
	ID_Word   = &TypeID{index: 4}
	ID_Char   = &TypeID{index: 5}
	ID_Str    = &TypeID{index: 6}
	ID_Thread = &TypeID{index: 7}
	ID_TypeID = &TypeID{index: 8}
)

const maxIndex = (1 << 28) - 1 // 0xffffff
type ListIDs []uint64

// TypeID
//
//		CPPPPPPP.TTTTTTTT.TTTTTTTT.TTTTTTTT.LLLLLLLL.LLLLLLLL.LLLLLLLL.LLLLLLLL
//		C = Constant flag
//	 P = Pointer depth
//	 T = Index
//	 L = ListIndex
type TypeID struct {
	// index refers to the delcaration or primitive index.
	// Map[K, V] may have a list index of 16, for example.
	// Primitives always have the same indices.
	index uint32

	// listIndex represents generic argument types.
	// Non-Generic types have a list index of 0.
	// Map[string, int] might have a list index of 4, for example.
	listIndex uint32
}

func NewTypeID(index, listIndex uint32) ast.TypeID {
	return &TypeID{index: index, listIndex: listIndex}
}
func (tid *TypeID) Index() uint32 {
	return tid.index
}

func (tid *TypeID) DeclIndex() uint32 {
	return (tid.index << 8) >> 8
}

func (tid *TypeID) ListIndex() uint32 {
	return tid.listIndex
}

// ID is the 64-bit identifier for a concrete type.
func (tid *TypeID) ID() uint64 {
	return (uint64(tid.index) << 31) | uint64(tid.listIndex)
}

func (tid *TypeID) IsConstant() bool {
	return tid.index>>31 != 0
}

func (tid *TypeID) PointerDepth() uint8 {
	return uint8(tid.index>>24) & 0x7f
}

func (tid *TypeID) IsMarkable(c *Types) bool {
	return tid.DeclIndex() > ID_TypeID.index+uint32(len(c.Enums)+len(c.Structs))
}
