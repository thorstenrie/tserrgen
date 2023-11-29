// Copyright (c) 2023 thorstenrie.
// All Rights Reserved. Use is governed with GNU Affero General Public License v3.0
// that can be found in the LICENSE file.
package tserrgen

import "fmt"

type Code struct {
	c string
}

func NewCode() *Code {
	return &Code{}
}

func (code *Code) String() string {
	if code == nil {
		return ""
	}
	return code.c
}

// Line comment
func (code *Code) LineComment(c string) *Code {
	if code == nil {
		return nil
	}
	code.c += fmt.Sprintf("// %v\n", c)
	return code
}

// Block End with new line
func (code *Code) FuncEnd() *Code {
	if code == nil {
		return nil
	}
	code.c += "}\n\n"
	return code
}

// Block End
func (code *Code) BlockEnd() *Code {
	if code == nil {
		return nil
	}
	code.c += "}\n"
	return code
}

// Function Call
func (code *Code) Call(n string) *Code {
	if code == nil {
		return nil
	}
	code.c += fmt.Sprintf("%v(", n)
	return code
}

// Parameters End + new line
func (code *Code) ParamEndln() *Code {
	if code == nil {
		return nil
	}
	code.c += ")\n"
	return code
}

// Parameters End
func (code *Code) ParamEnd() *Code {
	if code == nil {
		return nil
	}
	code.c += ")"
	return code
}

type Func1Args struct {
	Name, Var, Type, Return string
}

func (code *Code) Func1(a *Func1Args) *Code {
	if code == nil {
		return nil
	}
	code.c += fmt.Sprintf("func %v(%v %v) %v {\n", a.Name, a.Var, a.Type, a.Return)
	return code
}

// Type declaration for struct type
func (code *Code) TypeStruct(n string) *Code {
	if code == nil {
		return nil
	}
	code.c += fmt.Sprintf("type %v struct {\n", n)
	return code
}

type TypeArgs struct {
	Name, Type string
}

// Type
func (code *Code) Type(a *TypeArgs) *Code {
	if code == nil {
		return nil
	}
	code.c += fmt.Sprintf("%v %v\n", a.Name, a.Type)
	return code
}

// IdentifierLIst
func (code *Code) List() *Code {
	if code == nil {
		return nil
	}
	code.c += ", "
	return code
}

// IdentifierList with new line
func (code *Code) Listln() *Code {
	if code == nil {
		return nil
	}
	code.c += ",\n"
	return code
}

type SelArgs struct {
	Val, Sel string
}

// Field Selector
func (code *Code) SelField(a *SelArgs) *Code {
	if code == nil {
		return nil
	}
	code.c += fmt.Sprintf("%v.%v", a.Val, a.Sel)
	return code
}

// Method Selector
func (code *Code) SelMethod(a *SelArgs) *Code {
	if code == nil {
		return nil
	}
	code.c += fmt.Sprintf("%v.%v(", a.Val, a.Sel)
	return code
}

// Expression
type IfArgs struct {
	ExprLeft, ExprRight, Operator string
}

// If statement
func (code *Code) If(a *IfArgs) *Code {
	if code == nil {
		return nil
	}
	code.c += fmt.Sprintf("if %v %v %v {\n", a.ExprLeft, a.Operator, a.ExprRight)
	return code
}

type IfErrArgs struct {
	Method, Operator string
}

// If statement for error handling using a simple statement
func (code *Code) IfErr(a *IfErrArgs) *Code {
	if code == nil {
		return nil
	}
	code.c += fmt.Sprintf("if err := %v; err %v nil {\n", a.Method, a.Operator)
	return code
}

// Return
func (code *Code) Return() *Code {
	if code == nil {
		return nil
	}
	code.c += "return "
	return code
}

// Address operator
func (code *Code) Addr() *Code {
	if code == nil {
		return nil
	}
	code.c += "&"
	return code
}

// Identifier
func (code *Code) Ident(n string) *Code {
	if code == nil {
		return nil
	}
	code.c += n
	return code
}

type AssignmentArgs struct {
	ExprLeft, ExprRight string
}

// Assignment
func (code *Code) Assignment(a *AssignmentArgs) *Code {
	if code == nil {
		return nil
	}
	code.c += fmt.Sprintf("%v = %v", a.ExprLeft, a.ExprRight)
	return code
}

// Composite Literal
func (code *Code) CompositeLit(LiteralType string) *Code {
	if code == nil {
		return nil
	}
	code.c += fmt.Sprintf("%v{", LiteralType)
	return code
}

type ShortVarDeclArgs struct {
	Ident, Expr string
}

// Short variable declaration
func (code *Code) ShortVarDecl(a *ShortVarDeclArgs) *Code {
	if code == nil {
		return nil
	}
	code.c += fmt.Sprintf("%v := %v\n", a.Ident, a.Expr)
	return code
}

type KeyedElementArgs struct {
	Key, Element string
}

// Keyed element of a Composite literal
func (code *Code) KeyedElement(a *KeyedElementArgs) *Code {
	if code == nil {
		return nil
	}
	code.c += fmt.Sprintf("%v: %v,\n", a.Key, a.Element)
	return code
}
