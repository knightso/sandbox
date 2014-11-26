// +build !appengine
package lib

import "syscall"

func init() {
	delegate = func() string {
		wd, err := syscall.Getwd()
		if err != nil {
			return err.Error()
		}
		return "Hello! wd:" + wd
	}
}
