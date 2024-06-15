package decl

import (
	"geth-cody/ast"
	"geth-cody/compile/data"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
	"geth-cody/strs"

	"go.uber.org/zap"
)

type BaseDecl struct {
	StartToken token.Token
	EndToken   token.Token

	Name_   token.Token
	IsClass bool

	Methods_       map[string]ast.Method
	StaticMethods_ map[string]ast.Method

	Members_       map[string]ast.Member
	StaticMembers_ map[string]ast.Member

	Declarations_ map[string]ast.DeclField
}

func newDecl() BaseDecl {
	return BaseDecl{
		Methods_:       map[string]ast.Method{},
		StaticMethods_: map[string]ast.Method{},

		Members_:       map[string]ast.Member{},
		StaticMembers_: map[string]ast.Member{},

		Declarations_: map[string]ast.DeclField{},
	}
}

func (d *BaseDecl) Name() *token.Token {
	return &d.Name_
}
func (d *BaseDecl) Members() map[string]ast.Member {
	return d.Members_
}
func (d *BaseDecl) StaticMembers() map[string]ast.Member {
	return d.StaticMembers_
}
func (d *BaseDecl) Methods() map[string]ast.Method {
	return d.Methods_
}
func (d *BaseDecl) StaticMethods() map[string]ast.Method {
	return d.StaticMethods_
}
func (d *BaseDecl) Declarations() map[string]ast.DeclField {
	return d.Declarations_
}

func (d *BaseDecl) AddField(f ast.Field) io.Error {
	name := f.Name().Value
	if decl, ok := f.(ast.DeclField); ok {
		if name == d.Name_.Value {
			return io.NewError("inner decl name duplicates outer decl",
				zap.Any("decl", name),
				zap.Any("location", decl.Location()),
			)
		} else if _, exists := d.Declarations_[name]; exists {
			return io.NewError("duplicate decl name",
				zap.Any("decl", name),
				zap.Any("location", decl.Location()),
			)
		}
		d.Declarations_[decl.Name().Value] = decl
		return nil
	}

	_, methodExists := d.Methods_[name]
	_, staticMethodExists := d.StaticMethods_[name]
	_, memberExists := d.Members_[name]
	_, staticMemberExists := d.StaticMembers_[name]
	if methodExists || staticMethodExists || memberExists || staticMemberExists {
		return io.NewError("duplicate field name",
			zap.Any("member", name),
			zap.Any("location", f.Location()),
		)
	}

	if m, ok := f.(ast.Method); ok {
		if m.HasModifier(ast.MOD_STATIC) {
			d.StaticMethods_[m.Name().Value] = m
		} else {
			d.Methods_[m.Name().Value] = m
		}
	} else if m, ok := f.(ast.Member); ok {
		if m.HasModifier(ast.MOD_STATIC) {
			d.StaticMembers_[m.Name().Value] = m
		} else {
			d.Members_[m.Name().Value] = m
		}
	}

	return nil
}

func (d *BaseDecl) Location() token.Location {
	return token.LocationBetween(&d.StartToken, &d.EndToken)
}

func (d *BaseDecl) Semantic(p ast.SemanticParser) io.Error {
	for _, decl := range d.Declarations_ {
		if err := decl.Semantic(p); err != nil {
			return err
		}
	}

	d.addStaticFields(p.Scope())
	for _, member := range d.StaticMembers_ {
		if err := member.Semantic(p); err != nil {
			return err
		}
	}

	for _, method := range d.StaticMethods_ {
		if err := method.Semantic(p); err != nil {
			return err
		}
	}

	d.addFields(p.Scope())
	for _, member := range d.Members_ {
		if err := member.Semantic(p); err != nil {
			return err
		}
	}
	for _, method := range d.Methods_ {
		if err := method.Semantic(p); err != nil {
			return err
		}
	}

	return nil
}

func parentsString(superClass ast.DeclType, implements []ast.DeclType) string {
	var parents []ast.DeclType
	if superClass != nil {
		parents = append(parents, superClass)
	}

	return strs.Strings(append(parents, implements...))
}

func Syntax(p ast.SyntaxParser) (ast.Declaration, io.Error) {
	var declaration ast.Declaration
	switch t := p.Peek(); t.Type {
	case token.TOK_INTERFACE:
		declaration = newInterface()
	case token.TOK_ABSTRACT:
		declaration = newAbstract()
	case token.TOK_CLASS:
		declaration = newClass()
	case token.TOK_STRUCT:
		declaration = newStruct()
	case token.TOK_ENUM:
		declaration = newEnum()
	default:
		return nil, io.NewError("expected declaration", zap.Any("token", t))
	}

	if err := declaration.Syntax(p); err != nil {
		return nil, err
	}

	return declaration, nil
}

func (d *BaseDecl) addStaticFields(scope *ast.Scope) io.Error {
	for _, fields := range []map[string]ast.Member{d.Members_, d.StaticMembers_} {
		for _, member := range fields {
			if err := scope.AddVariable(member); err != nil {
				return err
			}
		}
	}

	return nil
}

func (d *BaseDecl) addFields(scope *ast.Scope) io.Error {
	for _, fields := range []map[string]ast.Member{d.Members_, d.StaticMembers_} {
		for _, member := range fields {
			if err := scope.AddVariable(member); err != nil {
				return err
			}
		}
	}

	return nil
}

func (d *BaseDecl) LinkParents(p ast.SemanticParser, visitedDecls *data.AsyncSet[ast.Declaration], cycleMap map[string]struct{}) io.Error {
	for _, decl := range d.Declarations_ {
		if err := decl.LinkParents(p, visitedDecls); err != nil {
			return err
		}
	}
	return nil
}

func (child *BaseDecl) Extends(p ast.SemanticParser, parent ast.Declaration, visitedDecls *data.AsyncSet[ast.Declaration]) io.Error {
	for name, parentMethod := range parent.Methods() {
		if childMethod, exists := child.Methods_[name]; exists {
			if _, err := p.TypeContext().MustExtend(childMethod.Type(), parentMethod.Type()); err != nil {
				return err
			}
		} else if parentMethod.HasModifier(ast.MOD_VIRTUAL) && child.IsClass {
			return io.NewError("child missing method",
				zap.Any("method", name),
				zap.Any("parent", parent.Name()),
				zap.Any("location", child.Location()),
			)
		} else {
			child.Methods_[name] = parentMethod
		}
	}

	return nil
}
