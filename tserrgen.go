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
	if e := genApi(&m); e != nil {
		return e
	}
	if e := genApiTest(&m); e != nil {
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
		if i == 13 {
			continue
		}
		c := NewCode().Expr("errmsg" + v.Name).Literal("errmsg").Ident(fmt.Sprint(i)).List().Ident(v.Code).List().Ident("\"" + v.Msg + "\"").Close()
		if e := cf.WriteCode(c.String()); e != nil {
			return tserr.Op(&tserr.OpArgs{Op: "WriteCode", Fn: string(cf.Filepath()), Err: e})
		}
	}
	if e := cf.FinishFile(); e != nil {
		return tserr.Op(&tserr.OpArgs{Op: "FinishFile", Fn: string(cf.Filepath()), Err: e})
	}
	return nil
}

func genApiTest(m *tserrconfig) error {
	cf, err := NewCodefile(tsfio.Directory(m.Root.Path), tserr_api_test_go)
	if err != nil {
		return tserr.Op(&tserr.OpArgs{Op: "NewCodefile", Fn: m.Root.Path + string(tserr_api_test_go), Err: err})
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
		t, e = genApiTestFunc(&v)
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

func genApi(m *tserrconfig) error {
	cf, err := NewCodefile(tsfio.Directory(m.Root.Path), tserr_api_go)
	if err != nil {
		return tserr.Op(&tserr.OpArgs{Op: "NewCodefile", Fn: m.Root.Path + string(tserr_api_go), Err: err})
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
		t, e = genApiFunc(&v)
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
	c.If(&IfArgs{Operand1: "err", Operand2: "nil", Operator: "=="})
	c.SelFunc(&SelArgs{Var: "t", Sel: "Fatal"}).Ident("errNil").ParamCloseln().Close()
	c.Call("testValidJson").Ident("t").List().Ident("err").ParamCloseln()
	c.ShortVarDecl(&ShortVarDeclArgs{Ident: "emsg", Expr: "errmsg{"})
	c.SelVar(&SelArgs{Var: "em", Sel: "Id"}).Listln()
	c.SelVar(&SelArgs{Var: "em", Sel: "C"}).Listln()
	c.SelFunc(&SelArgs{Var: "fmt", Sel: "Sprintf"}).Ident("\"%v\"").List()
	c.SelFunc(&SelArgs{Var: "fmt", Sel: "Errorf"}).SelVar(&SelArgs{Var: "em", Sel: "M"}).List()
	c.Ident("a").ParamClose().ParamClose().Listln().Close()
	c.Call("testEqualJson").Ident("t").List().Ident("err").List().Ident("&emsg").ParamCloseln()
	c.FuncClose()
	return c.String(), nil
}

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
	c := NewCode().Comment(m.Comment).Comment(m.Param[0].Comment).Func1(&Func1Args{Name: m.Name, Var: m.Param[0].Name, Type: m.Param[0].Type, Return: "error"})
	c.Return().Call("errorf").Addr().Ident("errmsg" + m.Name).List().Ident(m.Param[0].Name).ParamCloseln().FuncClose()
	return c.String(), nil
}

func genApiTestFuncM(m *errmsg) (string, error) {
	return "", nil
}

func genApiFuncM(m *errmsg) (string, error) {
	l := len(m.Param)
	if l <= 1 {
		return "", tserr.Higher(&tserr.HigherArgs{Var: "number of parameters", Actual: int64(l), LowerBound: 2})
	}
	c := NewCode().Comment(m.Name + "Args holds the required arguments for the error function " + m.Name).TypeStruct(m.Name + "Args")
	for _, v := range m.Param {
		c.Comment(v.Comment).Var(v.Name, v.Type)
	}
	c.FuncClose().Comment(m.Comment).Func1(&Func1Args{Name: m.Name, Var: "a", Type: " *" + m.Name + "Args", Return: "error"})
	c.If(&IfArgs{Operand1: "a", Operand2: "nil", Operator: "=="}).Return().Call("NilPtr").ParamCloseln().Close()
	c.Return().Call("errorf").Addr().Ident("errmsg" + m.Name)
	for _, v := range m.Param {
		c.List().SelVar(&SelArgs{Var: "a", Sel: v.Name})
	}
	c.ParamCloseln().FuncClose()
	return c.String(), nil
}
