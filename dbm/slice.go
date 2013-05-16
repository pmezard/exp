// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Functions analogous to package "os".

package dbm

import (
	"github.com/cznic/exp/lldb"
)

// Slice represents a slice of an Array.
type Slice struct {
	a        *Array
	prefix   []interface{}
	from, to []byte
}

// Do calls f for every subscripts-value pair in s in ascending collation order
// of the subsripts.  Do returns non nil error for general errors (eg. file
// read error).  If f returns false or a non nil error then Do terminates and
// returns the value of error from f.
//
// Note: f can get called with a subscripts-value pair which actually may no
// longer exist - if some other goroutine introduces such data race.
// Coordination required to avoid this situation, if applicable/desirable, must
// be provided by the client of dbm.
func (s *Slice) Do(f func(subscripts, value []interface{}) (bool, error)) (err error) {
	var (
		enum      *lldb.BTreeEnumerator
		bk, bv    []byte
		k, v      []interface{}
		hit, more bool
		skip      = len(s.prefix)
		db        = s.a.db
		noVal     bool
		from      = append(bpack(s.a.prefix), s.from...)
		to        = append(bpack(s.a.prefix), s.to...)
		bprefix   = s.a.prefix
	)

	ok, err := s.a.validate(false)
	if !ok && err != nil {
		return err
	}

	if s.a.tree == nil {
		return
	}

	if t := s.a.tree; !t.IsMem() && t.Handle() == 1 {
		noVal = true
	}

	doLeave := false
	defer func() {
		if e := recover(); e != nil {
			rdbg("PANIC %T %#v", e, e)
			err = fmt.Errorf("%v", e)
		}
		if doLeave {
			db.leave(&err)
		}
	}()

	switch {
	case s.from == nil && s.to == nil:
		if err = db.enter(); err != nil {
			return
		}

		doLeave = true
		if enum, _, err = s.a.tree.Seek(bprefix); err != nil {
			return noEof(err)
		}

		for {
			if bk, bv, err = enum.Current(); err != nil {
				return noEof(err)
			}

			if len(bprefix) != 0 && collate(bk[:len(bprefix)], bprefix) > 0 {
				return nil
			}

			if k, err = lldb.DecodeScalars(bk); err != nil {
				return noEof(err)
			}

			if v, err = lldb.DecodeScalars(bv); err != nil {
				return noEof(err)
			}

			doLeave = false
			if db.leave(&err) != nil {
				return
			}

			if noVal && v != nil {
				v = []interface{}{0}
			}
			if more, err = f(k[skip:], v); !more || err != nil {
				return noEof(err)
			}

			if err = db.enter(); err != nil {
				return
			}

			doLeave = true
			if enum, hit, err = s.a.tree.Seek(bk); err != nil {
				return noEof(err)
			}

			if hit {
				if err = enum.Next(); err != nil {
					return noEof(err)
				}
			}
		}
	case s.from == nil && s.to != nil:
		if err = db.enter(); err != nil {
			return
		}

		doLeave = true
		if enum, _, err = s.a.tree.Seek(from); err != nil {
			return noEof(err)
		}

		for {
			if bk, bv, err = enum.Current(); err != nil {
				return noEof(err)
			}

			if collate(bk, to) > 0 {
				return
			}

			if k, err = lldb.DecodeScalars(bk); err != nil {
				return noEof(err)
			}

			if v, err = lldb.DecodeScalars(bv); err != nil {
				return noEof(err)
			}

			doLeave = false
			if db.leave(&err) != nil {
				return
			}

			if noVal && v != nil {
				v = []interface{}{0}
			}
			if more, err = f(k[skip:], v); !more || err != nil {
				return noEof(err)
			}

			if err = db.enter(); err != nil {
				return
			}

			doLeave = true
			if enum, hit, err = s.a.tree.Seek(bk); err != nil {
				return noEof(err)
			}

			if hit {
				if err = enum.Next(); err != nil {
					return noEof(err)
				}
			}
		}
	case s.from != nil && s.to == nil:
		if err = db.enter(); err != nil {
			return
		}

		doLeave = true
		if enum, _, err = s.a.tree.Seek(from); err != nil {
			return noEof(err)
		}

		for {
			if bk, bv, err = enum.Current(); err != nil {
				return noEof(err)
			}

			if len(bprefix) != 0 && collate(bk[:len(bprefix)], bprefix) > 0 {
				return nil
			}

			if k, err = lldb.DecodeScalars(bk); err != nil {
				return noEof(err)
			}

			if v, err = lldb.DecodeScalars(bv); err != nil {
				return noEof(err)
			}

			doLeave = false
			if db.leave(&err) != nil {
				return
			}

			if noVal && v != nil {
				v = []interface{}{0}
			}
			if more, err = f(k[skip:], v); !more || err != nil {
				return noEof(err)
			}

			if err = db.enter(); err != nil {
				return
			}

			doLeave = true
			if enum, hit, err = s.a.tree.Seek(bk); err != nil {
				return noEof(err)
			}

			if hit {
				if err = enum.Next(); err != nil {
					return noEof(err)
				}
			}
		}
	case s.from != nil && s.to != nil:
		if err = db.enter(); err != nil {
			return
		}

		doLeave = true
		if enum, _, err = s.a.tree.Seek(from); err != nil {
			return noEof(err)
		}

		for {
			if bk, bv, err = enum.Current(); err != nil {
				return noEof(err)
			}

			if collate(bk, to) > 0 {
				return nil
			}

			if k, err = lldb.DecodeScalars(bk); err != nil {
				return noEof(err)
			}

			if v, err = lldb.DecodeScalars(bv); err != nil {
				return noEof(err)
			}

			doLeave = false
			if db.leave(&err) != nil {
				return
			}

			if noVal && v != nil {
				v = []interface{}{0}
			}
			if more, err = f(k[skip:], v); !more || err != nil {
				return noEof(err)
			}

			if err = db.enter(); err != nil {
				return
			}

			doLeave = true
			if enum, hit, err = s.a.tree.Seek(bk); err != nil {
				return noEof(err)
			}

			if hit {
				if err = enum.Next(); err != nil {
					return noEof(err)
				}
			}
		}
	}
	return noEof(err)
}
