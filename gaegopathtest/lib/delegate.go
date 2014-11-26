// +build !appengine
package lib

import "syscall"
import "fmt"

func init() {
	delegate = func() string {
		return fmt.Sprintf("Hello! pid:", syscall.Getpid())
	}
}
