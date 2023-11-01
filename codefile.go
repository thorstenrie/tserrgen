// Copyright (c) 2023 thorstenrie.
// All Rights Reserved. Use is governed with GNU Affero General Public License v3.0
// that can be found in the LICENSE file.
package tserrgen

import (
	"go/format"

	"github.com/thorstenrie/tserr"
	"github.com/thorstenrie/tsfio"
)

type Codefile struct {
	fn tsfio.Filename
	fp tsfio.Filename
}

const (
	headerSuffix = ".header"
	footerSuffix = ".footer"
)

func NewCodefile(dn tsfio.Directory, fn tsfio.Filename) (*Codefile, error) {
	f, err := tsfio.Path(dn, fn)
	if err != nil {
		return nil, tserr.Op(&tserr.OpArgs{Op: "Path", Fn: string(dn) + string(fn), Err: err})
	}
	cf := &Codefile{fp: f, fn: fn}
	return cf, nil
}

func (cf *Codefile) Filepath() tsfio.Filename {
	return cf.fp
}

func (cf *Codefile) StartFile() error {
	if e := tsfio.ResetFile(cf.fp); e != nil {
		return tserr.Op(&tserr.OpArgs{Op: "ResetFile", Fn: string(cf.fp), Err: e})
	}
	fh := cf.fn + headerSuffix
	if e := tsfio.AppendFile(&tsfio.Append{FileA: cf.fp, FileI: fh}); e != nil {
		return tserr.Op(&tserr.OpArgs{Op: "AppendFile", Fn: string(cf.fp), Err: e})
	}
	return nil
}

func (cf *Codefile) WriteCode(c string) error {
	return tsfio.WriteStr(cf.fp, c)
}

func (cf *Codefile) FinishFile() error {
	fe := cf.fn + footerSuffix
	if e := tsfio.AppendFile(&tsfio.Append{FileA: cf.fp, FileI: fe}); e != nil {
		return tserr.Op(&tserr.OpArgs{Op: "AppendFile", Fn: string(cf.fp), Err: e})
	}
	if e := cf.Format(); e != nil {
		return tserr.Op(&tserr.OpArgs{Op: "format", Fn: string(cf.fp), Err: e})
	}
	return nil
}

func (cf *Codefile) Format() error {
	i, e := tsfio.ReadFile(cf.fp)
	if e != nil {
		return tserr.Op(&tserr.OpArgs{Op: "ReadFile", Fn: string(cf.fp), Err: e})
	}
	o, e := format.Source(i)
	if e != nil {
		return tserr.Op(&tserr.OpArgs{Op: "Source", Fn: string(cf.fp), Err: e})
	}
	if e := tsfio.WriteSingleStr(cf.fp, string(o)); e != nil {
		return tserr.Op(&tserr.OpArgs{Op: "WriteSingleStr", Fn: string(cf.fp), Err: e})
	}
	return nil
}
