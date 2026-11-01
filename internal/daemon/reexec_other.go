//go:build !unix

package daemon

import "errors"

func execSelf(string, []string, []string) error { return errors.New("re-exec is not supported here") }
