package typeid

import (
	"encoding/binary"
	"fmt"
	"geth-cody/ast"
	"geth-cody/io"
	"math"
	"sync"
)

type ListKey string

func newListKey(ids []uint64) ListKey {
	bytes := make([]byte, len(ids)*4)
	for _, id := range ids {
		bytes = binary.LittleEndian.AppendUint64(bytes, id)
	}

	return ListKey(bytes)
}

type Types struct {
	Enums      []ast.Declaration
	Structs    []ast.Declaration
	Classes    []ast.Declaration
	Abstracts  []ast.Declaration
	Interfaces []ast.Declaration
	declMu     sync.Mutex

	ListsMap map[ListKey]uint32
	Lists    []ListIDs
	listsMu  sync.RWMutex
}

func NewTypes() *Types {
	return &Types{ListsMap: map[ListKey]uint32{}}
}

// region Next
func (c *Types) hasTooManyTypes() io.Error {
	if int(ID_TypeID.index)+len(c.Enums)+len(c.Structs)+
		len(c.Classes)+len(c.Abstracts)+len(c.Interfaces) == maxIndex {
		return io.NewError("Total number of declarations are limited to 2^28-1")
	}

	return nil
}

func (c *Types) next(decl ast.Declaration, decls []ast.Declaration) (uint32, io.Error) {
	c.declMu.Lock()
	defer c.declMu.Unlock()

	if err := c.hasTooManyTypes(); err != nil {
		return 0, err
	}
	decls = append(decls, decl)

	return uint32(len(decls) - 1), nil
}

func (c *Types) NextEnumIndex(d ast.Declaration) (uint32, io.Error) {
	return c.next(d, c.Enums)
}

func (c *Types) NextStructIndex(d ast.Declaration) (uint32, io.Error) {
	return c.next(d, c.Structs)
}

func (c *Types) NextClassIndex(d ast.Declaration) (uint32, io.Error) {
	return c.next(d, c.Classes)
}

func (c *Types) NextAbstractIndex(d ast.Declaration) (uint32, io.Error) {
	return c.next(d, c.Abstracts)
}

func (c *Types) NextInterfaceIndex(d ast.Declaration) (uint32, io.Error) {
	return c.next(d, c.Interfaces)
}

func (c *Types) NextListIndex(ids []uint64) (uint32, io.Error) {
	key := newListKey(ids)

	c.listsMu.RLock()
	if id, ok := c.ListsMap[key]; ok {
		c.listsMu.RUnlock()
		return id, nil
	}
	c.listsMu.RUnlock()

	c.listsMu.Lock()
	defer c.listsMu.Unlock()
	nextList := uint32(len(c.Lists))
	if nextList == math.MaxUint32 {
		return 0, io.NewError("Total number of unique list types is limited to 2^32-1")
	}

	c.ListsMap[key] = nextList
	c.Lists = append(append(c.Lists, []uint64{uint64(len(ids))}), ids)

	return nextList, nil
}

// endregion

// region Index

func (*Types) EnumIndex(index uint32) uint32 {
	return ID_TypeID.index + index
}

func (c *Types) StructIndex(index uint32) uint32 {
	return ID_TypeID.index + uint32(len(c.Enums)) + index
}

func (c *Types) ClassIndex(index uint32) uint32 {
	return ID_TypeID.index + uint32(len(c.Enums)+len(c.Structs)) + index
}

func (c *Types) AbstractIndex(index uint32) uint32 {
	return ID_TypeID.index + uint32(len(c.Enums)+len(c.Structs)+len(c.Classes)) + index
}

func (c *Types) InterfaceIndex(index uint32) uint32 {
	return ID_TypeID.index + uint32(len(c.Enums)+len(c.Structs)+len(c.Classes)+len(c.Abstracts)) + index
}

func (c *Types) ListIndex(ids []uint64) uint32 {
	key := newListKey(ids)
	id, ok := c.ListsMap[key]
	if !ok {
		panic(fmt.Sprintf("list key not found: %+v", ids))
	}

	return id
}

func (c *Types) MaxIndex() uint64 {
	return uint64(ID_TypeID.index) + uint64(len(c.Enums)+len(c.Structs)+
		len(c.Classes)+len(c.Abstracts)+len(c.Interfaces)+len(c.Lists))
}

// endregion
