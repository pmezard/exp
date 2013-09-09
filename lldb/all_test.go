// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lldb

import (
	"encoding/hex"
	"os"
	"runtime"
	"time"
)

const (
	testDbName = "_test.db"
	walName    = "_wal"
)

func now() time.Time { return time.Now() }

func hdump(b []byte) string {
	return hex.Dump(b)
}

func die() {
	os.Exit(1)
}

func stack() string {
	buf := make([]byte, 1<<16)
	return string(buf[:runtime.Stack(buf, false)])
}
