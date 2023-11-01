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

func (code *Code) Comment(c string) *Code {
	if code == nil {
		return nil
	}
	code.c += fmt.Sprintf("// %v\n", c)
	return code
}

func (code *Code) FuncClose() *Code {
	if code == nil {
		return nil
	}
	code.c += "}\n\n"
	return code
}

func (code *Code) Close() *Code {
	if code == nil {
		return nil
	}
	code.c += "}\n"
	return code
}

func (code *Code) Call(n string) *Code {
	if code == nil {
		return nil
	}
	code.c += fmt.Sprintf("%v(", n)
	return code
}

func (code *Code) ParamCloseln() *Code {
	if code == nil {
		return nil
	}
	code.c += ")\n"
	return code
}

func (code *Code) ParamClose() *Code {
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

func (code *Code) TypeStruct(n string) *Code {
	if code == nil {
		return nil
	}
	code.c += fmt.Sprintf("type %v struct {\n", n)
	return code
}

func (code *Code) Var(a string, t string) *Code {
	if code == nil {
		return nil
	}
	code.c += fmt.Sprintf("%v %v\n", a, t)
	return code
}

func (code *Code) List() *Code {
	if code == nil {
		return nil
	}
	code.c += ", "
	return code
}

func (code *Code) Listln() *Code {
	if code == nil {
		return nil
	}
	code.c += ",\n"
	return code
}

type SelArgs struct {
	Var, Sel string
}

func (code *Code) SelVar(a *SelArgs) *Code {
	if code == nil {
		return nil
	}
	code.c += fmt.Sprintf("%v.%v", a.Var, a.Sel)
	return code
}

func (code *Code) SelFunc(a *SelArgs) *Code {
	if code == nil {
		return nil
	}
	code.c += fmt.Sprintf("%v.%v(", a.Var, a.Sel)
	return code
}

type IfArgs struct {
	Operand1, Operand2, Operator string
}

func (code *Code) If(a *IfArgs) *Code {
	if code == nil {
		return nil
	}
	code.c += fmt.Sprintf("if %v %v %v {\n", a.Operand1, a.Operator, a.Operand2)
	return code
}

func (code *Code) Return() *Code {
	if code == nil {
		return nil
	}
	code.c += "return "
	return code
}

func (code *Code) Addr() *Code {
	if code == nil {
		return nil
	}
	code.c += "&"
	return code
}

func (code *Code) Ident(n string) *Code {
	if code == nil {
		return nil
	}
	code.c += n
	return code
}

func (code *Code) Expr(n string) *Code {
	if code == nil {
		return nil
	}
	code.c += fmt.Sprintf("%v = ", n)
	return code
}

// t literal type
func (code *Code) Literal(t string) *Code {
	if code == nil {
		return nil
	}
	code.c += fmt.Sprintf("%v{", t)
	return code
}

type ShortVarDeclArgs struct {
	Ident, Expr string
}

func (code *Code) ShortVarDecl(a *ShortVarDeclArgs) *Code {
	if code == nil {
		return nil
	}
	code.c += fmt.Sprintf("%v := %v\n", a.Ident, a.Expr)
	return code
}
