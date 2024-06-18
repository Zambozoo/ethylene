package decl

import (
	"geth-cody/ast"
	"geth-cody/io"

	"go.uber.org/zap"
)

type GenericDecl struct {
	TypesMap   map[string]ast.DeclType // Generic type parameters
	Types      []ast.DeclType
	TypesCount int
}

func newGenericDecl() GenericDecl {
	return GenericDecl{
		TypesMap: map[string]ast.DeclType{},
	}
}

func (gd *GenericDecl) PutGeneric(name string, generic ast.DeclType) io.Error {
	if _, exists := gd.TypesMap[name]; exists {
		return io.NewError("Duplicate generic type parameter",
			zap.String("name", name),
			zap.Any("location", generic.Location()),
		)
	}
	gd.TypesMap[name] = generic
	gd.Types = append(gd.Types, generic)
	gd.TypesCount++
	return nil
}

func (gd *GenericDecl) GenericsMap() map[string]ast.DeclType {
	return gd.TypesMap
}

func (gd *GenericDecl) Generics() []ast.DeclType {
	return gd.Types
}

func (gd *GenericDecl) GenericsCount() int {
	return gd.TypesCount
}
