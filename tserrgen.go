// Copyright (c) 2023 thorstenrie.
// All Rights Reserved. Use is governed with GNU Affero General Public License v3.0
// that can be found in the LICENSE file.
package tserrgen

import (
	"encoding/json"
	"fmt"

	"github.com/thorstenrie/tserr"
	"github.com/thorstenrie/tsfio"
)

var (
	tserr_messages_go = tsfio.Filename("tserr_messages.go")
	tserr_api_go      = tsfio.Filename("tserr_api.go")
	tserr_api_test_go = tsfio.Filename("tserr_api_test.go")
	tserr_testcases   = map[string]string{
		"string":  "strFoo",
		"error":   "errFoo",
		"int64":   "intFoo",
		"float64": "floatFoo",
	}
)

func Generate(fn tsfio.Filename) error {
	b, err := tsfio.ReadFile(fn)
	if err != nil {
		return tserr.Op(&tserr.OpArgs{Op: "ReadFile", Fn: string(fn), Err: err})
	}
	var m tserrconfig
	if e := json.Unmarshal(b, &m); e != nil {
		return tserr.Op(&tserr.OpArgs{Op: "Unmarshal", Fn: string(fn), Err: err})
	}
	if e := genMessages(&m); e != nil {
		return e
	}
	if e := genApi(&m, tserr_api_go, genApiFunc); e != nil {
		return e
	}
	if e := genApi(&m, tserr_api_test_go, genApiTestFunc); e != nil {
		return e
	}
	return nil
}

func genMessages(m *tserrconfig) error {
	cf, err := NewCodefile(tsfio.Directory(m.Root.Path), tserr_messages_go)
	if err != nil {
		return tserr.Op(&tserr.OpArgs{Op: "NewCodefile", Fn: m.Root.Path + string(tserr_messages_go), Err: err})
	}
	if e := cf.StartFile(); e != nil {
		return tserr.Op(&tserr.OpArgs{Op: "StartFile", Fn: string(cf.Filepath()), Err: e})
	}
	for i, v := range m.Root.Errors {
		if i++; i >= 13 {
			i++
		}
		c := NewCode().Assignment(&AssignmentArgs{ExprLeft: "errmsg" + v.Name}).CompositeLit("errmsg").Ident(fmt.Sprint(i)).List().Ident(v.Code).List().Ident("\"" + v.Msg + "\"").BlockEnd()
		if e := cf.WriteCode(c.String()); e != nil {
			return tserr.Op(&tserr.OpArgs{Op: "WriteCode", Fn: string(cf.Filepath()), Err: e})
		}
	}
	if e := cf.FinishFile(); e != nil {
		return tserr.Op(&tserr.OpArgs{Op: "FinishFile", Fn: string(cf.Filepath()), Err: e})
	}
	return nil
}

func genApi(m *tserrconfig, fn tsfio.Filename, genf func(*errmsg) (string, error)) error {
	cf, err := NewCodefile(tsfio.Directory(m.Root.Path), fn)
	if err != nil {
		return tserr.Op(&tserr.OpArgs{Op: "NewCodefile", Fn: m.Root.Path + string(fn), Err: err})
	}
	if e := cf.StartFile(); e != nil {
		return tserr.Op(&tserr.OpArgs{Op: "StartFile", Fn: string(cf.Filepath()), Err: e})
	}
	for _, v := range m.Root.Errors {
		var (
			l int    = len(v.Param)
			t string = ""
			e error
		)
		if l == 0 {
			return tserr.Empty("Param")
		}
		t, e = genf(&v)
		if e != nil {
			return e
		}
		if e := cf.WriteCode(t); e != nil {
			return tserr.Op(&tserr.OpArgs{Op: "WriteCode", Fn: string(cf.Filepath()), Err: e})
		}
	}
	if e := cf.FinishFile(); e != nil {
		return tserr.Op(&tserr.OpArgs{Op: "FinishFile", Fn: string(cf.Filepath()), Err: e})
	}
	return nil
}

func genApiTestFunc(m *errmsg) (string, error) {
	l := len(m.Param)
	if l == 1 {
		return genApiTestFunc1(m)
	} else if l > 1 {
		return genApiTestFuncM(m)
	} else {
		return "", tserr.Higher(&tserr.HigherArgs{Var: "number of parameters", Actual: int64(l), LowerBound: 1})
	}
}

func genApiFunc(m *errmsg) (string, error) {
	l := len(m.Param)
	if l == 1 {
		return genApiFunc1(m)
	} else if l > 1 {
		return genApiFuncM(m)
	} else {
		return "", tserr.Higher(&tserr.HigherArgs{Var: "number of parameters", Actual: int64(l), LowerBound: 1})
	}
}

func genApiTestFunc1(m *errmsg) (string, error) {
	if m == nil {
		return "", tserr.NilPtr()
	}
	if m.Param == nil {
		return "", tserr.NilPtr()
	}
	l := len(m.Param)
	if l != 1 {
		return "", tserr.Equal(&tserr.EqualArgs{Var: "number of parameters", Actual: int64(l), Want: 1})
	}
	c := NewCode().Func1(&Func1Args{Name: "Test" + m.Name, Var: "t", Type: "*testing.T", Return: ""})
	t, e := foo(m.Param[0].Type)
	if e != nil {
		return "", e
	}
	c.ShortVarDecl(&ShortVarDeclArgs{Ident: "a", Expr: t})
	c.ShortVarDecl(&ShortVarDeclArgs{Ident: "em", Expr: "&errmsg" + m.Name})
	c.ShortVarDecl(&ShortVarDeclArgs{Ident: "err", Expr: m.Name + "(a)"})
	c.If(&IfArgs{ExprLeft: "err", ExprRight: "nil", Operator: "=="})
	c.SelMethod(&SelArgs{Val: "t", Sel: "Fatal"}).Ident("errNil").ParamEndln().BlockEnd()
	c.Call("testValidJson").Ident("t").List().Ident("err").ParamEndln()
	c, e = genEmsgTest(c, m)
	if e != nil {
		return "", e
	}
	c.Call("testEqualJson").Ident("t").List().Ident("err").List().Ident("&emsg").ParamEndln()
	c.FuncEnd()
	return c.String(), nil
}

func genEmsgTest(c *Code, m *errmsg) (*Code, error) {
	if (c == nil) || (m == nil) {
		return nil, tserr.NilPtr()
	}
	if m.Param == nil {
		return nil, tserr.NilPtr()
	}
	l := len(m.Param)
	if l == 0 {
		return nil, tserr.Higher(&tserr.HigherArgs{Var: "number of parameters", Actual: int64(l), LowerBound: 1})
	}
	c.ShortVarDecl(&ShortVarDeclArgs{Ident: "emsg", Expr: "errmsg{"})
	c.SelField(&SelArgs{Val: "em", Sel: "Id"}).Listln()
	c.SelField(&SelArgs{Val: "em", Sel: "C"}).Listln()
	c.SelMethod(&SelArgs{Val: "fmt", Sel: "Sprintf"}).Ident("\"%v\"").List()
	c.SelMethod(&SelArgs{Val: "fmt", Sel: "Errorf"}).SelField(&SelArgs{Val: "em", Sel: "M"}).List()
	if l == 1 {
		c.Ident("a")
	} else if l > 1 {
		for _, v := range m.Param {
			c.SelField(&SelArgs{Val: "a", Sel: v.Name}).List()
		}
	}
	c.ParamEnd().ParamEnd().Listln().BlockEnd()
	return c, nil
}

// Todo: Automate test variables
func foo(t string) (string, error) {
	r, ok := tserr_testcases[t]
	if ok {
		return r, nil
	}
	types := "[ "
	for i := range tserr_testcases {
		types += i + " "
	}
	types += "]"
	return "", tserr.TypeNotMatching(&tserr.TypeNotMatchingArgs{Act: t, Want: types})
}

// Todo non-printable
func genApiFunc1(m *errmsg) (string, error) {
	l := len(m.Param)
	if l != 1 {
		return "", tserr.Equal(&tserr.EqualArgs{Var: "number of parameters", Actual: int64(l), Want: 1})
	}
	c := NewCode().LineComment(m.Comment).LineComment(m.Param[0].Comment).Func1(&Func1Args{Name: m.Name, Var: m.Param[0].Name, Type: m.Param[0].Type, Return: "error"})
	c.Return().Call("errorf").Addr().Ident("errmsg" + m.Name).List().Ident(m.Param[0].Name).ParamEndln().FuncEnd()
	return c.String(), nil
}

func genApiTestFuncM(m *errmsg) (string, error) {
	l := len(m.Param)
	if l <= 1 {
		return "", tserr.Higher(&tserr.HigherArgs{Var: "number of parameters", Actual: int64(l), LowerBound: 2})
	}
	c := NewCode().Func1(&Func1Args{Name: "Test" + m.Name + "Nil", Var: "t", Type: "*testing.T", Return: ""})
	c.IfErr(&IfErrArgs{Method: m.Name + "(nil)", Operator: "=="}).SelMethod(&SelArgs{Val: "t", Sel: "Errorf"})
	c.Ident("errNil").ParamEndln().BlockEnd().FuncEnd()
	c.Func1(&Func1Args{Name: "Test" + m.Name, Var: "t", Type: "*testing.T", Return: ""})
	c.ShortVarDecl(&ShortVarDeclArgs{Ident: "a", Expr: m.Name + "Args{"})
	for _, v := range m.Param {
		t, e := foo(v.Type)
		if e != nil {
			return "", e
		}
		c.KeyedElement(&KeyedElementArgs{Key: v.Name, Element: t})
	}
	c.BlockEnd()
	c.ShortVarDecl(&ShortVarDeclArgs{Ident: "em", Expr: "&errmsg" + m.Name})
	c.ShortVarDecl(&ShortVarDeclArgs{Ident: "err", Expr: m.Name + "(&a)"})
	c.If(&IfArgs{ExprLeft: "err", ExprRight: "nil", Operator: "=="})
	c.SelMethod(&SelArgs{Val: "t", Sel: "Fatal"}).Ident("errNil").ParamEndln().BlockEnd()
	c.Call("testValidJson").Ident("t").List().Ident("err").ParamEndln()
	c, e := genEmsgTest(c, m)
	if e != nil {
		return "", e
	}
	c.Call("testEqualJson").Ident("t").List().Ident("err").List().Addr().Ident("emsg").ParamEndln()
	c.FuncEnd()
	return c.String(), nil
}

func genApiFuncM(m *errmsg) (string, error) {
	l := len(m.Param)
	if l <= 1 {
		return "", tserr.Higher(&tserr.HigherArgs{Var: "number of parameters", Actual: int64(l), LowerBound: 2})
	}
	c := NewCode().LineComment(m.Name + "Args holds the required arguments for the error function " + m.Name).TypeStruct(m.Name + "Args")
	for _, v := range m.Param {
		c.LineComment(v.Comment).Type(&TypeArgs{Name: v.Name, Type: v.Type})
	}
	c.FuncEnd().LineComment(m.Comment).Func1(&Func1Args{Name: m.Name, Var: "a", Type: " *" + m.Name + "Args", Return: "error"})
	c.If(&IfArgs{ExprLeft: "a", ExprRight: "nil", Operator: "=="}).Return().Call("NilPtr").ParamEndln().BlockEnd()
	c.Return().Call("errorf").Addr().Ident("errmsg" + m.Name)
	for _, v := range m.Param {
		c.List().SelField(&SelArgs{Val: "a", Sel: v.Name})
	}
	c.ParamEndln().FuncEnd()
	return c.String(), nil
}
